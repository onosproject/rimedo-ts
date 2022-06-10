// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0
// Copy from https://github.com/woojoong88/onos-kpimon/tree/sample-a1t-xapp/pkg/northbound/a1
// modified by RIMEDO Labs team

package a1

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	a1tapi "github.com/onosproject/onos-api/go/onos/a1t/a1"
	"github.com/onosproject/onos-lib-go/pkg/logging/service"
	"google.golang.org/grpc"
)

var SampleJSON1 = `
{
  "scope": {
    "ueId": "0000000000000001"
  },
  "tspResources": [
    {
      "cellIdList": [
        {"plmnId": {"mcc": "248","mnc": "35"},
          "cId": {"ncI": 39}},
        {"plmnId": {"mcc": "248","mnc": "35"},
         "cId": {"ncI": 40}}
      ], 
      "preference": "PREFER"
    },
    {
      "cellIdList": [
        {"plmnId": {"mcc": "248","mnc": "35"},
          "cId": {"ncI": 81}},
        {"plmnId": {"mcc": "248","mnc": "35"},
          "cId": {"ncI": 82}},
        {"plmnId": {"mcc": "248","mnc": "35"},
         "cId": {"ncI": 83}}
      ],
      "preference": "FORBID"
    }
  ]
}
`

var SampleJSON2 = `
{
  "scope": {
    "ueId": "0000000000000002"
  },
  "tspResources": [
    {
      "cellIdList": [
        {"plmnId": {"mcc": "248","mnc": "35"},
          "cId": {"ncI": 39}},
        {"plmnId": {"mcc": "248","mnc": "35"},
         "cId": {"ncI": 40}}
      ], 
      "preference": "PREFER"
    },
    {
      "cellIdList": [
        {"plmnId": {"mcc": "248","mnc": "35"},
          "cId": {"ncI": 81}},
        {"plmnId": {"mcc": "248","mnc": "35"},
          "cId": {"ncI": 82}},
        {"plmnId": {"mcc": "248","mnc": "35"},
         "cId": {"ncI": 83}}
      ],
      "preference": "FORBID"
    }
  ]
}
`

var SampleEnforcedStatus = `
{
  "enforceStatus": "ENFORCED"
}
`

var SampleNotEnforcedStatus = `
{
  "enforceStatus": "NOT_ENFORCED",
  "enforceReason": "SCOPE_NOT_APPLICABLE"
}
`

var SampleNotEnforcedPolicyID = "2"

func NewA1PService(policyMap *map[string][]byte, notifier chan bool) service.Service {
	return &A1PService{
		TsPolicyTypeMap: policyMap,
		notifier:        notifier,
	}
}

type A1PService struct {
	TsPolicyTypeMap *map[string][]byte
	notifier        chan bool
}

func (a *A1PService) Register(s *grpc.Server) {
	server := &A1PServer{
		TsPolicyTypeMap: *a.TsPolicyTypeMap,
		StatusUpdateCh:  make(chan *a1tapi.PolicyStatusMessage),
		notifier:        a.notifier,
	}
	a1tapi.RegisterPolicyServiceServer(s, server)
}

type A1PServer struct {
	TsPolicyTypeMap map[string][]byte
	StatusUpdateCh  chan *a1tapi.PolicyStatusMessage
	notifier        chan bool
	mu              sync.RWMutex
}

func (a *A1PServer) PolicySetup(ctx context.Context, message *a1tapi.PolicyRequestMessage) (*a1tapi.PolicyResultMessage, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var result map[string]interface{}
	json.Unmarshal(message.Message.Payload, &result)

	if message.PolicyType.Id != "ORAN_TrafficSteeringPreference_2.0.0" {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy type does not support",
				},
			},
		}
		return res, nil
	}

	if _, ok := a.TsPolicyTypeMap[message.PolicyId]; ok {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy ID already exists",
				},
			},
		}
		return res, nil
	}

	a.TsPolicyTypeMap[message.PolicyId] = message.Message.Payload

	go func() {
		if message.NotificationDestination != "" {
			statusUpdateMsg := &a1tapi.PolicyStatusMessage{
				PolicyId:   message.PolicyId,
				PolicyType: message.PolicyType,
				Message: &a1tapi.StatusMessage{
					Header: &a1tapi.Header{
						RequestId:   uuid.New().String(),
						AppId:       message.Message.Header.AppId,
						Encoding:    message.Message.Header.Encoding,
						PayloadType: a1tapi.PayloadType_STATUS,
					},
				},
				NotificationDestination: message.NotificationDestination,
			}

			if message.PolicyId == SampleNotEnforcedPolicyID {
				statusUpdateMsg.Message.Payload = []byte(SampleNotEnforcedStatus)
			} else {
				statusUpdateMsg.Message.Payload = []byte(SampleEnforcedStatus)
			}

			a.StatusUpdateCh <- statusUpdateMsg
		}
	}()

	res := &a1tapi.PolicyResultMessage{
		PolicyId:   message.PolicyId,
		PolicyType: message.PolicyType,
		Message: &a1tapi.ResultMessage{
			Header: &a1tapi.Header{
				PayloadType: message.Message.Header.PayloadType,
				RequestId:   message.Message.Header.RequestId,
				Encoding:    message.Message.Header.Encoding,
				AppId:       message.Message.Header.AppId,
			},
			Payload: a.TsPolicyTypeMap[message.PolicyId],
			Result: &a1tapi.Result{
				Success: true,
			},
		},
		NotificationDestination: message.NotificationDestination,
	}
	a.notifier <- true
	return res, nil
}

func (a *A1PServer) PolicyUpdate(ctx context.Context, message *a1tapi.PolicyRequestMessage) (*a1tapi.PolicyResultMessage, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result map[string]interface{}
	json.Unmarshal(message.Message.Payload, &result)

	if message.PolicyType.Id != "ORAN_TrafficSteeringPreference_2.0.0" {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy type does not support",
				},
			},
		}
		return res, nil
	}

	if _, ok := a.TsPolicyTypeMap[message.PolicyId]; !ok {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy ID does not exists",
				},
			},
		}
		return res, nil
	}

	a.TsPolicyTypeMap[message.PolicyId] = message.Message.Payload

	go func() {
		if message.NotificationDestination != "" {
			statusUpdateMsg := &a1tapi.PolicyStatusMessage{
				PolicyId:   message.PolicyId,
				PolicyType: message.PolicyType,
				Message: &a1tapi.StatusMessage{
					Header: &a1tapi.Header{
						RequestId:   uuid.New().String(),
						AppId:       message.Message.Header.AppId,
						Encoding:    message.Message.Header.Encoding,
						PayloadType: a1tapi.PayloadType_STATUS,
					},
				},
				NotificationDestination: message.NotificationDestination,
			}

			if message.PolicyId == SampleNotEnforcedPolicyID {
				statusUpdateMsg.Message.Payload = []byte(SampleNotEnforcedStatus)
			} else {
				statusUpdateMsg.Message.Payload = []byte(SampleEnforcedStatus)
			}

			a.StatusUpdateCh <- statusUpdateMsg
		}
	}()

	res := &a1tapi.PolicyResultMessage{
		PolicyId:   message.PolicyId,
		PolicyType: message.PolicyType,
		Message: &a1tapi.ResultMessage{
			Header: &a1tapi.Header{
				PayloadType: message.Message.Header.PayloadType,
				RequestId:   message.Message.Header.RequestId,
				Encoding:    message.Message.Header.Encoding,
				AppId:       message.Message.Header.AppId,
			}, Payload: a.TsPolicyTypeMap[message.PolicyId],
			Result: &a1tapi.Result{
				Success: true,
			},
		},
		NotificationDestination: message.NotificationDestination,
	}
	a.notifier <- true
	return res, nil
}

func (a *A1PServer) PolicyDelete(ctx context.Context, message *a1tapi.PolicyRequestMessage) (*a1tapi.PolicyResultMessage, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result map[string]interface{}
	json.Unmarshal(message.Message.Payload, &result)

	if message.PolicyType.Id != "ORAN_TrafficSteeringPreference_2.0.0" {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy type does not support",
				},
			},
		}
		return res, nil
	}

	if _, ok := a.TsPolicyTypeMap[message.PolicyId]; !ok {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy ID does not exists",
				},
			},
		}
		return res, nil
	}

	delete(a.TsPolicyTypeMap, message.PolicyId)

	res := &a1tapi.PolicyResultMessage{
		PolicyId:   message.PolicyId,
		PolicyType: message.PolicyType,
		Message: &a1tapi.ResultMessage{
			Header: &a1tapi.Header{
				PayloadType: message.Message.Header.PayloadType,
				RequestId:   message.Message.Header.RequestId,
				Encoding:    message.Message.Header.Encoding,
				AppId:       message.Message.Header.AppId,
			}, Payload: a.TsPolicyTypeMap[message.PolicyId],
			Result: &a1tapi.Result{
				Success: true,
			},
		},
	}
	a.notifier <- true
	return res, nil
}

func (a *A1PServer) PolicyQuery(ctx context.Context, message *a1tapi.PolicyRequestMessage) (*a1tapi.PolicyResultMessage, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var result map[string]interface{}
	json.Unmarshal(message.Message.Payload, &result)

	if message.PolicyType.Id != "ORAN_TrafficSteeringPreference_2.0.0" {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy type does not support",
				},
			},
		}
		return res, nil
	}

	if message.PolicyId == "" {

		listPolicies := make([]string, 0)
		for k := range a.TsPolicyTypeMap {
			listPolicies = append(listPolicies, k)
		}

		listPoliciesJson, err := json.Marshal(listPolicies)
		if err != nil {
			log.Error(err)
		}

		res := &a1tapi.PolicyResultMessage{
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: listPoliciesJson,
				Result: &a1tapi.Result{
					Success: true,
				},
			},
		}
		return res, nil
	}

	if _, ok := a.TsPolicyTypeMap[message.PolicyId]; !ok {
		res := &a1tapi.PolicyResultMessage{
			PolicyId:   message.PolicyId,
			PolicyType: message.PolicyType,
			Message: &a1tapi.ResultMessage{
				Header: &a1tapi.Header{
					PayloadType: message.Message.Header.PayloadType,
					RequestId:   message.Message.Header.RequestId,
					Encoding:    message.Message.Header.Encoding,
					AppId:       message.Message.Header.AppId,
				}, Payload: message.Message.Payload,
				Result: &a1tapi.Result{
					Success: false,
					Reason:  "Policy ID does not exists",
				},
			},
		}
		return res, nil
	}

	resultMsg := &a1tapi.PolicyResultMessage{
		PolicyId:   message.PolicyId,
		PolicyType: message.PolicyType,
		Message: &a1tapi.ResultMessage{
			Header: &a1tapi.Header{
				PayloadType: message.Message.Header.PayloadType,
				RequestId:   message.Message.Header.RequestId,
				Encoding:    message.Message.Header.Encoding,
				AppId:       message.Message.Header.AppId,
			},
			Result: &a1tapi.Result{
				Success: true,
			},
		},
	}

	switch message.Message.Header.PayloadType {
	case a1tapi.PayloadType_POLICY:
		resultMsg.Message.Payload = a.TsPolicyTypeMap[message.PolicyId]
		resultMsg.Message.Header.PayloadType = a1tapi.PayloadType_POLICY
	case a1tapi.PayloadType_STATUS:
		resultMsg.Message.Header.PayloadType = a1tapi.PayloadType_STATUS
		if message.PolicyId == SampleNotEnforcedPolicyID {
			resultMsg.Message.Payload = []byte(SampleNotEnforcedStatus)

		} else {
			resultMsg.Message.Payload = []byte(SampleEnforcedStatus)
		}
	}
	return resultMsg, nil
}

func (a *A1PServer) PolicyStatus(server a1tapi.PolicyService_PolicyStatusServer) error {

	watchers := make(map[uuid.UUID]chan *a1tapi.PolicyAckMessage)
	mu := sync.RWMutex{}

	go func(m *sync.RWMutex) {
		for {
			ack, err := server.Recv()
			if err != nil {
				log.Error(err)
			}
			m.Lock()
			for _, v := range watchers {
				select {
				case v <- ack:
					log.Debugf("Sent msg %v on %v", ack, v)
				default:
					log.Debugf("Failed to send msg %v on %v", ack, v)
				}
			}
			m.Unlock()
		}
	}(&a.mu)

	for msg := range a.StatusUpdateCh {
		watcherID := uuid.New()
		ackCh := make(chan *a1tapi.PolicyAckMessage)
		timerCh := make(chan bool, 1)
		go func(ch chan bool) {
			time.Sleep(5 * time.Second)
			timerCh <- true
			close(timerCh)
		}(timerCh)

		go func(m *sync.RWMutex) {
			for {
				select {
				case ack := <-ackCh:
					if ack.Message.Header.RequestId == msg.Message.Header.RequestId {
						m.Lock()
						close(ackCh)
						delete(watchers, watcherID)
						m.Unlock()
						return
					}
				case <-timerCh:
					log.Error(fmt.Errorf("could not receive PolicyACKMessage in timer"))
					m.Lock()
					close(ackCh)
					delete(watchers, watcherID)
					m.Unlock()
					return
				}
			}
		}(&a.mu)

		mu.Lock()
		watchers[watcherID] = ackCh
		mu.Unlock()

		err := server.Send(msg)
		if err != nil {
			log.Error(err)
			mu.Lock()
			close(ackCh)
			delete(watchers, watcherID)
			mu.Unlock()
		}
	}

	return nil
}
