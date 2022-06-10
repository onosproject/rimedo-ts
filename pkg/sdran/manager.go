// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0

package sdran

import (
	"context"
	"strconv"
	"sync"

	policyAPI "github.com/onosproject/onos-a1-dm/go/policy_schemas/traffic_steering_preference/v2"
	e2tAPI "github.com/onosproject/onos-api/go/onos/e2t/e2"
	e2api "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/pdubuilder"
	e2sm_v2_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-v2-ies"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/logging/service"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	control "github.com/onosproject/onos-mho/pkg/mho"
	"github.com/onosproject/onos-mho/pkg/store"
	"github.com/onosproject/rimedo-ts/pkg/mho"
	"github.com/onosproject/rimedo-ts/pkg/policy"
	"github.com/onosproject/rimedo-ts/pkg/rnib"
	"github.com/onosproject/rimedo-ts/pkg/southbound/e2"
)

var log = logging.GetLogger("rimedo-ts", "sdran", "manager")

type Config struct {
	AppID              string
	E2tAddress         string
	E2tPort            int
	TopoAddress        string
	TopoPort           int
	SMName             string
	SMVersion          string
	TSPolicySchemePath string
}

func NewManager(config Config, flag bool) *Manager {

	ueStore := store.NewStore()
	cellStore := store.NewStore()
	onosPolicyStore := store.NewStore()

	policyMap := make(map[string]*mho.PolicyData)

	indCh := make(chan *mho.E2NodeIndication)
	ctrlReqChs := make(map[string]chan *e2api.ControlMessage)

	options := e2.Options{
		AppID:       config.AppID,
		E2tAddress:  config.E2tAddress,
		E2tPort:     config.E2tPort,
		TopoAddress: config.TopoAddress,
		TopoPort:    config.TopoPort,
		SMName:      config.SMName,
		SMVersion:   config.SMVersion,
	}

	e2Manager, err := e2.NewManager(options, indCh, ctrlReqChs)
	if err != nil {
		log.Warn(err)
	}

	manager := &Manager{
		e2Manager:       e2Manager,
		mhoCtrl:         mho.NewController(indCh, ueStore, cellStore, onosPolicyStore, policyMap, flag),
		policyManager:   policy.NewPolicyManager(&policyMap),
		ueStore:         ueStore,
		cellStore:       cellStore,
		onosPolicyStore: onosPolicyStore,
		ctrlReqChs:      ctrlReqChs,
		services:        []service.Service{},
		mutex:           sync.RWMutex{},
	}
	return manager
}

type Manager struct {
	e2Manager       e2.Manager
	mhoCtrl         *mho.Controller
	policyManager   *policy.PolicyManager
	ueStore         store.Store
	cellStore       store.Store
	onosPolicyStore store.Store
	ctrlReqChs      map[string]chan *e2api.ControlMessage
	services        []service.Service
	mutex           sync.RWMutex
}

func (m *Manager) Run(flag *bool) {
	if err := m.start(flag); err != nil {
		log.Fatal("Unable to run Manager", err)
	}
}

func (m *Manager) start(flag *bool) error {
	m.startNorthboundServer()
	err := m.e2Manager.Start()
	if err != nil {
		log.Warn(err)
		return err
	}

	go m.mhoCtrl.Run(context.Background(), flag)

	return nil
}

func (m *Manager) startNorthboundServer() error {

	s := northbound.NewServer(northbound.NewServerCfg(
		"",
		"",
		"",
		int16(5150),
		true,
		northbound.SecurityConfig{}))

	for i := range m.services {
		s.AddService(m.services[i])
	}

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}

func (m *Manager) AddService(service service.Service) {

	m.services = append(m.services, service)

}

func (m *Manager) GetUEs(ctx context.Context) map[string]mho.UeData {
	output := make(map[string]mho.UeData)
	chEntries := make(chan *store.Entry, 1024)
	err := m.ueStore.Entries(ctx, chEntries)
	if err != nil {
		log.Warn(err)
		return output
	}
	for entry := range chEntries {
		ueData := entry.Value.(mho.UeData)
		output[ueData.UeID] = ueData
	}
	return output
}

func (m *Manager) GetCells(ctx context.Context) map[string]mho.CellData {
	output := make(map[string]mho.CellData)
	chEntries := make(chan *store.Entry, 1024)
	err := m.cellStore.Entries(ctx, chEntries)
	if err != nil {
		log.Warn(err)
		return output
	}
	for entry := range chEntries {
		cellData := entry.Value.(mho.CellData)
		output[cellData.CGIString] = cellData
	}
	return output
}

func (m *Manager) GetPolicies(ctx context.Context) map[string]mho.PolicyData {
	output := make(map[string]mho.PolicyData)
	chEntries := make(chan *store.Entry, 1024)
	err := m.onosPolicyStore.Entries(ctx, chEntries)
	if err != nil {
		log.Warn(err)
		return output
	}
	for entry := range chEntries {
		policyData := entry.Value.(mho.PolicyData)
		output[policyData.Key] = policyData
	}
	return output
}

func (m *Manager) GetCellTypes(ctx context.Context) map[string]rnib.Cell {
	return m.e2Manager.GetCellTypes(ctx)
}

func (m *Manager) SetCellType(ctx context.Context, cellID string, cellType string) error {
	return m.e2Manager.SetCellType(ctx, cellID, cellType)
}

func (m *Manager) GetCell(ctx context.Context, CGI string) *mho.CellData {

	return m.mhoCtrl.GetCell(ctx, CGI)

}

func (m *Manager) SetCell(ctx context.Context, cell *mho.CellData) {

	m.mhoCtrl.SetCell(ctx, cell)

}

func (m *Manager) AttachUe(ctx context.Context, ue *mho.UeData, CGI string, cgiObject *e2sm_v2_ies.Cgi) {

	m.mhoCtrl.AttachUe(ctx, ue, CGI, cgiObject)

}

func (m *Manager) GetUe(ctx context.Context, ueID string) *mho.UeData {

	return m.mhoCtrl.GetUe(ctx, ueID)

}

func (m *Manager) SetUe(ctx context.Context, ueData *mho.UeData) {

	m.mhoCtrl.SetUe(ctx, ueData)

}

func (m *Manager) CreatePolicy(ctx context.Context, key string, policy *policyAPI.API) *mho.PolicyData {

	return m.mhoCtrl.CreatePolicy(ctx, key, policy)

}

func (m *Manager) GetPolicy(ctx context.Context, key string) *mho.PolicyData {

	return m.mhoCtrl.GetPolicy(ctx, key)

}

func (m *Manager) SetPolicy(ctx context.Context, key string, policy *mho.PolicyData) {

	m.mhoCtrl.SetPolicy(ctx, key, policy)

}

func (m *Manager) DeletePolicy(ctx context.Context, key string) {

	m.mhoCtrl.DeletePolicy(ctx, key)

}

func (m *Manager) GetPolicyStore() *store.Store {
	return m.mhoCtrl.GetPolicyStore()
}

func (m *Manager) GetControlChannelsMap(ctx context.Context) map[string]chan *e2api.ControlMessage {
	return m.ctrlReqChs
}

func (m *Manager) GetPolicyManager() *policy.PolicyManager {
	return m.policyManager
}

func (m *Manager) SwitchUeBetweenCells(ctx context.Context, ueID string, targetCellCGI string) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	availableUes := m.GetUEs(ctx)
	chosenUe := availableUes[ueID]

	if shouldBeSwitched(chosenUe, targetCellCGI) {

		targetCell := m.GetCell(ctx, targetCellCGI)
		servingCell := m.GetCell(ctx, chosenUe.CGIString)

		targetCell.CumulativeHandoversOut++
		servingCell.CumulativeHandoversIn++

		chosenUe.Idle = false
		m.AttachUe(ctx, &chosenUe, targetCellCGI, targetCell.CGI)

		m.SetCell(ctx, targetCell)
		m.SetCell(ctx, servingCell)

		controlChannel := m.ctrlReqChs[chosenUe.E2NodeID]

		controlHandler := &control.E2SmMhoControlHandler{
			NodeID:            chosenUe.E2NodeID,
			ControlAckRequest: e2tAPI.ControlAckRequest_NO_ACK,
		}

		ueIDnum, err := strconv.Atoi(chosenUe.UeID)
		if err != nil {
			log.Errorf("SendHORequest() failed to convert string %v to decimal number - assumption is not satisfied (UEID is a decimal number): %v", chosenUe.UeID, err)
		}

		ueIdentity, err := pdubuilder.CreateUeIDGNb(int64(ueIDnum), []byte{0xAA, 0xBB, 0xCC}, []byte{0xDD}, []byte{0xCC, 0xC0}, []byte{0xFC})
		if err != nil {
			log.Errorf("SendHORequest() Failed to create UEID: %v", err)
		}

		servingPlmnIDBytes := servingCell.CGI.GetNRCgi().GetPLmnidentity().GetValue()
		servingNCI := servingCell.CGI.GetNRCgi().GetNRcellIdentity().GetValue().GetValue()
		servingNCILen := servingCell.CGI.GetNRCgi().GetNRcellIdentity().GetValue().GetLen()

		go func() {
			if controlHandler.ControlHeader, err = controlHandler.CreateMhoControlHeader(servingNCI, servingNCILen, 1, servingPlmnIDBytes); err == nil {

				if controlHandler.ControlMessage, err = controlHandler.CreateMhoControlMessage(servingCell.CGI, ueIdentity, targetCell.CGI); err == nil {

					if controlRequest, err := controlHandler.CreateMhoControlRequest(); err == nil {

						controlChannel <- controlRequest
						log.Infof("CONTROL MESSAGE: UE [ID:%v, 5QI:%v] switched between CELLs [CGI:%v -> CGI:%v]\n", chosenUe.UeID, chosenUe.FiveQi, servingCell.CGIString, targetCell.CGIString)

					} else {
						log.Warn("Control request problem!", err)
					}
				} else {
					log.Warn("Control message problem!", err)
				}
			} else {
				log.Warn("Control header problem!", err)
			}
		}()

	}

}

func shouldBeSwitched(ue mho.UeData, cgi string) bool {

	servingCgi := ue.CGIString
	if servingCgi == cgi {
		return false
	}
	return true

}
