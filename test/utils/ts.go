// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"io"

	"github.com/onosproject/onos-api/go/onos/ransim/trafficsim"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	RanSimulatorAddress = "ran-simulator:5150"
)

func GetUEID() ([]*types.Ue, error) {
	result := make([]*types.Ue, 0)
	tlsConfig, err := creds.GetClientCredentials()
	if err != nil {
		return []*types.Ue{}, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}
	conn, err := grpc.DialContext(ctx, RanSimulatorAddress, opts...)
	if err != nil {
		return []*types.Ue{}, err
	}
	trafficSimClient := trafficsim.NewTrafficClient(conn)
	ues, err := trafficSimClient.ListUes(ctx, &trafficsim.ListUesRequest{})
	if err != nil {
		return []*types.Ue{}, err
	}

	for {
		ue, err := ues.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return []*types.Ue{}, err
		}
		if ue != nil {
			result = append(result, ue.GetUe())
		}
	}

	return result, nil
}
