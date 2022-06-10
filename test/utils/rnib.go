// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"fmt"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	toposdk "github.com/onosproject/onos-ric-sdk-go/pkg/topo"
)

func getXappFilter() *topoapi.Filters {
	controlRelationFilter := &topoapi.Filters{
		KindFilter: &topoapi.Filter{
			Filter: &topoapi.Filter_Equal_{
				Equal_: &topoapi.EqualFilter{
					Value: topoapi.XAPP,
				},
			},
		},
	}
	return controlRelationFilter
}

func UpdateXAppA1InterfaceIPAddr(ipAddress string) error {
	sdkClient, err := toposdk.NewClient()
	if err != nil {
		return err
	}
	objects, err := sdkClient.List(context.Background(), toposdk.WithListFilters(getXappFilter()))
	if err != nil {
		return err
	}

	if len(objects) != 1 {
		return fmt.Errorf("number of xApp should be 1, currently %v", len(objects))
	}

	xAppObject := &topoapi.XAppInfo{}
	err = objects[0].GetAspect(xAppObject)
	if err != nil {
		return err
	}

	for i := 0; i < len(xAppObject.Interfaces); i++ {
		if xAppObject.Interfaces[i].Type == topoapi.Interface_INTERFACE_A1_XAPP {
			xAppObject.Interfaces[i].IP = ipAddress
		}
	}

	err = objects[0].SetAspect(xAppObject)
	if err != nil {
		return err
	}

	err = sdkClient.Update(context.Background(), &objects[0])
	if err != nil {
		return err
	}

	return nil
}
