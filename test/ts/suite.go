// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ts

import (
	"github.com/onosproject/rimedo-ts/test/utils"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/input"
	"github.com/onosproject/helmit/pkg/test"
	testutils "github.com/onosproject/onos-ric-sdk-go/pkg/utils"
)

// TestSuite has sdran release and test suite
type TestSuite struct {
	sdran *helm.HelmRelease
	test.Suite
	c *input.Context
}

// SetupTestSuite prepares test suite setup
func (s *TestSuite) SetupTestSuite(c *input.Context) error {
	s.c = c
	// write files
	err := utils.WriteFile("/tmp/tls.cacrt", utils.TLSCacrt)
	if err != nil {
		return err
	}
	err = utils.WriteFile("/tmp/tls.crt", utils.TLSCrt)
	if err != nil {
		return err
	}
	err = utils.WriteFile("/tmp/tls.key", utils.TLSKey)
	if err != nil {
		return err
	}
	err = utils.WriteFile("/tmp/config.json", utils.ConfigJSON)
	if err != nil {
		return err
	}

	sdran, err := utils.CreateSdranRelease(c)
	if err != nil {
		return err
	}
	s.sdran = sdran
	r := sdran.Install(true)
	testutils.StartTestProxy()
	return r
}

// TearDownTestSuite uninstalls helm chart released
func (s *TestSuite) TearDownTestSuite() error {
	testutils.StopTestProxy()
	return nil
}
