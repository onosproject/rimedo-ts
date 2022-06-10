// Copy from onosproject/onos-mho/pkg/monitoring/monitor.go
// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0
// modified by RIMEDO-Labs team

package monitoring

import (
	"context"
	"github.com/onosproject/rimedo-ts/pkg/mho"

	e2api "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-mho/pkg/broker"
)

var log = logging.GetLogger("rimedo-ts", "monitoring")

func NewMonitor(streamReader broker.StreamReader, nodeID topoapi.ID, indChan chan *mho.E2NodeIndication, triggerType e2sm_mho.MhoTriggerType) *Monitor {
	return &Monitor{
		streamReader: streamReader,
		nodeID:       nodeID,
		indChan:      indChan,
		triggerType:  triggerType,
	}
}

type Monitor struct {
	streamReader broker.StreamReader
	nodeID       topoapi.ID
	indChan      chan *mho.E2NodeIndication
	triggerType  e2sm_mho.MhoTriggerType
}

func (m *Monitor) Start(ctx context.Context) error {
	errCh := make(chan error)
	go func() {
		for {
			indMsg, err := m.streamReader.Recv(ctx)
			if err != nil {
				log.Errorf("Error reading indication stream, chanID:%v, streamID:%v, err:%v", m.streamReader.ChannelID(), m.streamReader.StreamID(), err)
				errCh <- err
			}
			err = m.processIndication(ctx, indMsg, m.nodeID)
			if err != nil {
				log.Errorf("Error processing indication, err:%v", err)
				errCh <- err
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Monitor) processIndication(ctx context.Context, indication e2api.Indication, nodeID topoapi.ID) error {

	m.indChan <- &mho.E2NodeIndication{
		NodeID:      string(nodeID),
		TriggerType: m.triggerType,
		IndMsg: e2api.Indication{
			Payload: indication.Payload,
			Header:  indication.Header,
		},
	}

	return nil
}
