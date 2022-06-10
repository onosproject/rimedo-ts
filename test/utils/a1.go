// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	A1TAddress = "http://onos-a1t:9639"
)

func PutA1Policy(policy string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/policytypes/ORAN_TrafficSteeringPreference_2.0.0/policies/1", A1TAddress)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(policy)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("status code is not 200 or 201: %v", resp.StatusCode)
	}

	return nil
}
