// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0
// Copy from https://github.com/woojoong88/onos-kpimon/tree/sample-a1t-xapp/pkg/northbound/a1

package a1

import (
	"context"

	a1tapi "github.com/onosproject/onos-api/go/onos/a1t/a1"
	"github.com/onosproject/onos-lib-go/pkg/logging/service"
	"google.golang.org/grpc"
)

func NewA1EIService() service.Service {
	log.Debugf("A1EI service created")
	return &A1EIService{}
}

type A1EIService struct {
}

func (a *A1EIService) Register(s *grpc.Server) {
	server := &A1EIServer{}
	a1tapi.RegisterEIServiceServer(s, server)
}

type A1EIServer struct {
}

func (a *A1EIServer) EIQuery(server a1tapi.EIService_EIQueryServer) error {
	log.Debug("EIQuery stream established")
	ch := make(chan bool)
	<-ch
	return nil
}

func (a *A1EIServer) EIJobSetup(server a1tapi.EIService_EIJobSetupServer) error {
	log.Debug("EIJobSetup stream established")
	ch := make(chan bool)
	<-ch
	return nil
}

func (a *A1EIServer) EIJobUpdate(server a1tapi.EIService_EIJobUpdateServer) error {
	log.Debug("EIJobUpdate stream established")
	ch := make(chan bool)
	<-ch
	return nil
}

func (a *A1EIServer) EIJobDelete(server a1tapi.EIService_EIJobDeleteServer) error {
	log.Debug("EIJobDelete stream established")
	ch := make(chan bool)
	<-ch
	return nil
}

func (a *A1EIServer) EIJobStatusQuery(server a1tapi.EIService_EIJobStatusQueryServer) error {
	log.Debug("EIJobStatusQuery stream established")
	ch := make(chan bool)
	<-ch
	return nil
}

func (a *A1EIServer) EIJobStatusNotify(ctx context.Context, message *a1tapi.EIStatusMessage) (*a1tapi.EIAckMessage, error) {
	log.Debug("EIJobStatusNotify called %v", message)
	return nil, nil
}

func (a *A1EIServer) EIJobResultDelivery(ctx context.Context, message *a1tapi.EIResultMessage) (*a1tapi.EIAckMessage, error) {
	log.Debug("EIJobResultDelivery called %v", message)
	return nil, nil
}
