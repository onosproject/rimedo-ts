// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0

package ts

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/onosproject/onos-api/go/onos/ransim/types"

	//"context"
	"testing"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/rimedo-ts/pkg/manager"
	"github.com/onosproject/rimedo-ts/pkg/northbound/a1"
	"github.com/onosproject/rimedo-ts/pkg/sdran"
	"github.com/onosproject/rimedo-ts/test/utils"
)

// TestTsSm is the function for Helmit-based integration test
func (s *TestSuite) TestTsSm(t *testing.T) {
	// update xApp's IP and port number
	hostname, err := os.Hostname()
	if err != nil {
		t.Error(err)
		return
	}
	addr, err := net.LookupIP(hostname)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Hostname: %v", hostname)
	t.Logf("IP: %v", addr)
	err = os.Setenv("POD_NAME", hostname)
	if err != nil {
		t.Error(err)
		return
	}

	err = os.Setenv("POD_IP", addr[0].String())
	if err != nil {
		t.Error(err)
		return
	}

	sdranConfig := sdran.Config{
		AppID:              "rimedo-ts",
		E2tAddress:         "onos-e2t",
		E2tPort:            5150,
		TopoAddress:        "onos-topo",
		TopoPort:           5150,
		SMName:             "oran-e2sm-mho",
		SMVersion:          "v2",
		TSPolicySchemePath: "/data/schemas/ORAN_TrafficSteeringPreference_v102.json",
	}

	a1Config := a1.Config{
		PolicyName:        "ORAN_TrafficSteeringPreference",
		PolicyVersion:     "2.0.0",
		PolicyID:          "ORAN_TrafficSteeringPreference_2.0.0",
		PolicyDescription: "O-RAN traffic steering",
		A1tPort:           5150,
	}

	_, err = certs.HandleCertPaths("", "", "", true)
	if err != nil {
		t.Error(err)
		return
	}

	mgr := manager.NewManager(sdranConfig, a1Config, false)
	mgr.Run()

	// get UE ID
	ues, err := utils.GetUEID()
	if err != nil {
		t.Error(err)
		return
	} else if len(ues) != 1 {
		t.Errorf("the number of UEs should be 1, currently it is %v", len(ues))
		return
	}

	ueID := fmt.Sprintf("%016d", ues[0].IMSI)
	t.Logf("ueID: %s", ueID)

	// get a1 policy string
	a1Policy := strings.Replace(utils.A1Policy, "<IMSI>", ueID, 1)
	t.Log(a1Policy)

	// 341114881 - 138426010504510
	// 341098497 - 13842601454c001

	// put a1 policy
	err = utils.PutA1Policy(a1Policy)
	if err != nil {
		t.Error(err)
		return
	}

	// guard interval
	time.Sleep(10 * time.Second)

	// verification step
	// check serving cell - periodically for 2 mins
	// get sCell before verification
	var pNCGI types.NCGI

	for i := 0; i < utils.VerificationTimer; i++ {
		time.Sleep(1 * time.Second)
		if i < 30 {
			ues, err = utils.GetUEID()
			if err != nil {
				t.Error(err)
				return
			} else if len(ues) != 1 {
				t.Errorf("the number of UEs should be 1, currently it is %v", len(ues))
				return
			}
			pNCGI = ues[0].ServingTower
			t.Logf("[In Guard interval] UEID: %016d / NCGI: %x", ues[0].IMSI, ues[0].ServingTower)
			continue
		}

		ues, err := utils.GetUEID()
		if err != nil {
			t.Error(err)
			return
		} else if len(ues) != 1 {
			t.Errorf("the number of UEs should be 1, currently it is %v", len(ues))
			return
		}

		t.Logf("[Checking] UEID: %016d / NCGI: %x", ues[0].IMSI, ues[0].ServingTower)
		// if it is different, HO happens - test failed because the installed policy prohibit the HO.
		if ues[0].ServingTower != pNCGI {
			t.Errorf("Wrong sCell NCGI: original %v / changed %x", pNCGI, ues[0].ServingTower)
			return
		}
	}
}
