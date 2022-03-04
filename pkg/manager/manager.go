// Created by RIMEDO-Labs team
// based on any onosproject manager
package manager

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/RIMEDO-Labs/rimedo-ts/pkg/mho"
	"github.com/RIMEDO-Labs/rimedo-ts/pkg/northbound/a1"
	"github.com/RIMEDO-Labs/rimedo-ts/pkg/sdran"
	policyAPI "github.com/onosproject/onos-a1-dm/go/policy_schemas/traffic_steering_preference/v2"
	topoAPI "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("rimedo-ts", "ts-manager")

type Config struct {
	AppID       string
	E2tAddress  string
	E2tPort     int
	TopoAddress string
	TopoPort    int
	SMName      string
	SMVersion   string
}

func NewManager(sdranConfig sdran.Config, a1Config a1.Config) *Manager {

	sdranManager := sdran.NewManager(sdranConfig)

	a1PolicyTypes := make([]*topoAPI.A1PolicyType, 0)
	a1Policy := &topoAPI.A1PolicyType{
		Name:        topoAPI.PolicyTypeName(a1Config.PolicyName),
		Version:     topoAPI.PolicyTypeVersion(a1Config.PolicyVersion),
		ID:          topoAPI.PolicyTypeID(a1Config.PolicyID),
		Description: topoAPI.PolicyTypeDescription(a1Config.PolicyDescription),
	}
	a1PolicyTypes = append(a1PolicyTypes, a1Policy)

	a1Manager, err := a1.NewManager("", "", "", a1Config.A1tPort, sdranConfig.AppID, a1PolicyTypes)
	if err != nil {
		log.Warn(err)
	}

	manager := &Manager{
		sdranManager: sdranManager,
		a1Manager:    *a1Manager,
	}
	return manager
}

type Manager struct {
	sdranManager *sdran.Manager
	a1Manager    a1.Manager
}

func (m *Manager) Run() {

	if err := m.start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}

}

func (m *Manager) Close() {
	m.a1Manager.Close(context.Background())
}

func (m *Manager) start() error {

	ctx := context.Background()

	policyMap := make(map[string][]byte)

	policyChange := make(chan bool)

	m.sdranManager.AddService(a1.NewA1EIService())
	m.sdranManager.AddService(a1.NewA1PService(&policyMap, policyChange))

	m.sdranManager.Run()

	m.a1Manager.Start()

	go func() {
		for range policyChange {
			log.Debug("------------------------------------------")
			log.Debug("--------- POLICY STORE CHANGED! ----------")
			log.Debug("------------------------------------------")
			if err := m.updatePolicies(ctx, policyMap); err != nil {
				log.Warn("Some problems occured when updating Policy store!")
			}
			log.Debug("------------------------------------------")
			log.Debug("------------------------------------------")
			log.Debug("")
			m.checkPolicies(ctx, true)
		}

	}()
	flag := true
	time.Sleep(15 * time.Second)
	log.Info("\n\n\n\n\n\n\n\n\n\n")
	go func() {
		for {
			time.Sleep(1 * time.Second)
			m.checkPolicies(ctx, flag)
			flag = false
		}
	}()

	go func() {
		for {
			log.Debug("")
			log.Debug("------------------------------------------")
			log.Debug("----------- AVALIABLE NODES --------------")
			log.Debug("------------------------------------------")
			time.Sleep(15 * time.Second)
			m.showAvailableNodes(ctx)
			log.Debug("------------------------------------------")
			log.Debug("------------------------------------------")
			log.Debug("")
		}
	}()

	return nil
}

func (m *Manager) updatePolicies(ctx context.Context, policyMap map[string][]byte) error {
	policies := m.sdranManager.GetPolicies(ctx)
	for k := range policies {
		if _, ok := policyMap[k]; !ok {
			m.sdranManager.DeletePolicy(ctx, k)
			log.Infof("\nPOLICY  MESSAGE: Policy [ID:%v] deleted\n", k)
		}
	}
	for i := range policyMap {
		r, err := policyAPI.UnmarshalAPI(policyMap[i])
		if err == nil {
			policyObject := m.sdranManager.CreatePolicy(ctx, i, &r)
			info := fmt.Sprintf("\nPOLICY  MESSAGE: Policy [ID:%v] applied -> ", policyObject.Key)
			previous := false
			if policyObject.API.Scope.SliceID != nil {
				info = info + fmt.Sprintf("Slice [SD:%v, SST:%v, PLMN:(MCC:%v, MNC:%v)]", *policyObject.API.Scope.SliceID.SD, policyObject.API.Scope.SliceID.Sst, policyObject.API.Scope.SliceID.PlmnID.Mcc, policyObject.API.Scope.SliceID.PlmnID.Mnc)
				previous = true
			}
			if policyObject.API.Scope.UeID != nil {
				if previous {
					info = info + ", "
				}
				info = info + fmt.Sprintf("UE [ID:%v]", *policyObject.API.Scope.UeID)
				previous = true
			}
			if policyObject.API.Scope.QosID != nil {
				if previous {
					info = info + ", "
				}
				if policyObject.API.Scope.QosID.QcI != nil {
					info = info + fmt.Sprintf("QoS [QCI:%v]", *policyObject.API.Scope.QosID.QcI)
				}
				if policyObject.API.Scope.QosID.The5QI != nil {
					info = info + fmt.Sprintf("QoS [5QI:%v]", *policyObject.API.Scope.QosID.The5QI)
				}
			}
			if policyObject.API.Scope.CellID != nil {
				if previous {
					info = info + ", "
				}
				info = info + "CELL ["
				if policyObject.API.Scope.CellID.CID.NcI != nil {

					info = info + fmt.Sprintf("NCI:%v, ", *policyObject.API.Scope.CellID.CID.NcI)
				}
				if policyObject.API.Scope.CellID.CID.EcI != nil {

					info = info + fmt.Sprintf("ECI:%v, ", *policyObject.API.Scope.CellID.CID.EcI)
				}
				info = info + fmt.Sprintf("PLMN:(MCC:%v, MNC:%v)]", policyObject.API.Scope.CellID.PlmnID.Mcc, policyObject.API.Scope.CellID.PlmnID.Mnc)
			}
			for i := range policyObject.API.TSPResources {
				info = info + fmt.Sprintf(" - (%v) -", policyObject.API.TSPResources[i].Preference)
				for j := range policyObject.API.TSPResources[i].CellIDList {
					nci := *policyObject.API.TSPResources[i].CellIDList[j].CID.NcI
					plmnId, _ := mho.GetPlmnIdFromMccMnc(policyObject.API.TSPResources[i].CellIDList[j].PlmnID.Mcc, policyObject.API.TSPResources[i].CellIDList[j].PlmnID.Mnc)
					cgi := PlmnIDNciToCGI(plmnId, uint64(nci))
					info = info + fmt.Sprintf(" CELL [CGI:%v],", cgi)
				}
				info = info[0:len(info)-1] + "\n"

			}
			log.Info(info)
		} else {
			log.Warn("Can't unmarshal the JSON file!")
			return err
		}
	}
	return nil
}

func (m *Manager) deployPolicies(ctx context.Context) {

	policyManager := m.sdranManager.GetPolicyManager()
	ues := m.sdranManager.GetUEs(ctx)
	keys := make([]string, 0, len(ues))
	for k := range ues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {

		var cellIDs []policyAPI.CellID
		var rsrps []int
		fiveQi := ues[keys[i]].FiveQi
		sd := "456DEF"
		scopeUe := policyAPI.Scope{

			SliceID: &policyAPI.SliceID{
				SD:  &sd,
				Sst: 1,
				PlmnID: policyAPI.PlmnID{
					Mcc: "138",
					Mnc: "426",
				},
			},
			UeID: &keys[i],
			QosID: &policyAPI.QosID{
				The5QI: &fiveQi,
			},
		}

		cgiKeys := make([]string, 0, len(ues))
		for cgi := range ues[keys[i]].CgiTable {
			cgiKeys = append(cgiKeys, cgi)
		}

		for j := range cgiKeys {

			cgi := ues[keys[i]].CgiTable[cgiKeys[j]]
			nci := int64(mho.GetNciFromCellGlobalID(cgi))
			plmnIdBytes := mho.GetPlmnIDBytesFromCellGlobalID(cgi)
			plmnId := mho.PlmnIDBytesToInt(plmnIdBytes)
			mcc, mnc := mho.GetMccMncFromPlmnID(plmnId)
			cellID := policyAPI.CellID{
				CID: policyAPI.CID{
					NcI: &nci,
				},
				PlmnID: policyAPI.PlmnID{
					Mcc: mcc,
					Mnc: mnc,
				},
			}

			cellIDs = append(cellIDs, cellID)
			rsrps = append(rsrps, int(ues[keys[i]].RsrpTable[cgiKeys[j]]))

		}

		tsResult := policyManager.GetTsResultForUEV2(scopeUe, rsrps, cellIDs)
		plmnId, err := mho.GetPlmnIdFromMccMnc(tsResult.PlmnID.Mcc, tsResult.PlmnID.Mnc)
		if err != nil {
			log.Warnf("Cannot get PLMN ID from these MCC and MNC parameters:%v,%v.", tsResult.PlmnID.Mcc, tsResult.PlmnID.Mnc)
		} else {
			targetCellCGI := PlmnIDNciToCGI(plmnId, uint64(*tsResult.CID.NcI))
			m.sdranManager.SwitchUeBetweenCells(ctx, keys[i], targetCellCGI)
		}

		cellIDs = nil
		rsrps = nil

	}

}

func (m *Manager) checkPolicies(ctx context.Context, flag bool) {
	log.Debug("")
	log.Debug("------------------------------------------")
	log.Debug("-------------- POLIECIES -----------------")
	log.Debug("------------------------------------------")
	policies := m.sdranManager.GetPolicies(ctx)
	keys := make([]string, 0, len(policies))
	for k := range policies {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if flag && (len(policies) == 0) {
		log.Infof("\nPOLICY  MESSAGE: Default policy applied\n")
	}

	for _, key := range keys {
		log.Debugf("Key:%v Policy:%v", key, policies[key].API)
	}
	log.Debug("------------------------------------------")
	log.Debug("------------------------------------------")
	log.Debug("")
	m.deployPolicies(ctx)
}

func (m *Manager) showAvailableNodes(ctx context.Context) {
	log.Debug("")
	log.Debug("------------------------------------------")
	log.Debug("---------------- CELLS -------------------")
	log.Debug("------------------------------------------")
	cells := m.sdranManager.GetCellTypes(ctx)
	keys := make([]string, 0, len(cells))
	for k := range cells {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		log.Debugf("ID:%v CGI:%v CellType:%v", key, cells[key].CGI, cells[key].CellType)
	}
	log.Debug("------------------------------------------")
	log.Debug("------------------------------------------")
	log.Debug("")

	ues := m.sdranManager.GetUEs(ctx)
	keys = make([]string, 0, len(ues))
	for k := range ues {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	log.Debug("")
	log.Debug("------------------------------------------")
	log.Debug("----------------- UES --------------------")
	log.Debug("------------------------------------------")
	for _, key := range keys {
		ueIdString, _ := strconv.Atoi(key)
		info := fmt.Sprintf("ID:%v 5QI: %v CGI:%v RSRP:%v [", ueIdString, ues[key].FiveQi, ues[key].CGIString, ues[key].RsrpServing)

		neigh_keys := make([]string, 0, len(ues[key].RsrpNeighbors))
		for k := range ues[key].RsrpNeighbors {
			neigh_keys = append(neigh_keys, k)
		}
		sort.Strings(neigh_keys)
		for _, neigh := range neigh_keys {
			info += fmt.Sprintf("%v ", ues[key].RsrpNeighbors[neigh])
		}
		info += "]"
		log.Debug(info)
	}
	log.Debug("------------------------------------------")
	log.Debug("------------------------------------------")
	log.Debug("")
	log.Debug("")
	log.Debug("------------------------------------------")
	log.Debug("------------- UES [TABLE] ----------------")
	log.Debug("------------------------------------------")
	for _, key := range keys {
		ueIdString, _ := strconv.Atoi(key)
		info := fmt.Sprintf("ID:%v 5QI: %v CGI:%v [", ueIdString, ues[key].FiveQi, ues[key].CGIString)

		cgi_keys := make([]string, 0, len(ues[key].RsrpTable))
		for k := range ues[key].RsrpTable {
			cgi_keys = append(cgi_keys, k)
		}
		sort.Strings(cgi_keys)
		for _, cgi := range cgi_keys {
			info += fmt.Sprintf("%v ", cgi)
		}
		info += "]"
		log.Debug(info)
	}
	log.Debug("------------------------------------------")
	log.Debug("------------------------------------------")
	log.Debug("")
}

func (m *Manager) changeCellsTypes(ctx context.Context) {
	cellTypes := make(map[int]string)
	cellTypes[0] = "Macro"
	cellTypes[1] = "SmallCell"
	for {
		time.Sleep(10 * time.Second)
		cells := m.sdranManager.GetCellTypes(ctx)
		type_id := rand.Intn(len(cellTypes))
		for key, val := range cells {
			_ = val
			err := m.sdranManager.SetCellType(ctx, key, cellTypes[type_id])
			if err != nil {
				log.Warn(err)
			}
			break
		}

	}
}

func PlmnIDNciToCGI(plmnID uint64, nci uint64) string {
	cgi := strconv.FormatInt(int64(plmnID<<36|(nci&0xfffffffff)), 16)
	//cgi = cgi[0:2] + cgi[4:5] + cgi[3:4] + cgi[5:6] + cgi[2:3] + cgi[6:]
	return cgi[0:8] + cgi[13:14] + cgi[10:12] + cgi[8:10] + cgi[14:15] + cgi[12:13]
}
