// Created by RIMEDO-Labs team
package main

import (
	"github.com/onosproject/onos-lib-go/pkg/env"
	"os"
	"os/signal"
	"syscall"

	"github.com/RIMEDO-Labs/rimedo-ts/pkg/manager"
	"github.com/RIMEDO-Labs/rimedo-ts/pkg/northbound/a1"
	"github.com/RIMEDO-Labs/rimedo-ts/pkg/sdran"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("rimedo-ts")

func main() {

	log.SetLevel(logging.InfoLevel)
	log.Info("Starting RIMEDO Labs Traffic Steering xAPP")

	sdranConfig := sdran.Config{
		AppID:              env.GetPodName(),
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

	_, err := certs.HandleCertPaths("", "", "", true)
	if err != nil {
		log.Fatal(err)
	}

	mgr := manager.NewManager(sdranConfig, a1Config)
	mgr.Run()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	log.Debug("app: received a shutdown signal:", <-killSignal)
	mgr.Close()
}
