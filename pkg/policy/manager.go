// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0
// Created by RIMEDO-Labs team
package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	policyAPI "github.com/onosproject/onos-a1-dm/go/policy_schemas/traffic_steering_preference/v2"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/rimedo-ts/pkg/mho"
	"github.com/xeipuuv/gojsonschema"
)

var log = logging.GetLogger("rimedo-ts", "policy")

func NewPolicySchemaValidatorV2(path string) *PolicySchemaValidatorV2 {

	return &PolicySchemaValidatorV2{
		schemePath: path,
	}

}

type PolicySchemaValidatorV2 struct {
	schemePath string
}

func NewPolicyManager(policyMap *map[string]*mho.PolicyData) *PolicyManager {

	var POLICY_WEIGHTS = map[string]int{
		"DEFAULT": 0.0,
		"PREFER":  16.0,
		"AVOID":   -16.0,
		"SHALL":   1000.0,
		"FORBID":  -1000.0,
	}

	return &PolicyManager{
		validator:     NewPolicySchemaValidatorV2("schemePath"),
		policyMap:     policyMap,
		preferenceMap: POLICY_WEIGHTS,
	}

}

type PolicyManager struct {
	validator     *PolicySchemaValidatorV2
	policyMap     *map[string]*mho.PolicyData
	preferenceMap map[string]int
}

func (m *PolicyManager) ReadPolicyObjectFromFileV2(jsonPath string, policyObject *mho.PolicyData) error {

	jsonFile, err := m.LoadPolicyJsonFromFileV2(jsonPath)
	if err != nil {
		log.Error("Couldn't read PolicyObject from file")
		return err
	}

	var ok bool
	ok, err = m.ValidatePolicyJsonSchemaV2(jsonPath)
	if err != nil {
		log.Error("Error validating json scheme")
		return err
	}
	if !ok {
		return errors.New("the json file is invalid")
	}
	if err = m.UnmarshalPolicyJsonV2(jsonFile, policyObject); err != nil {
		log.Error("Error unmarshaling json file")
		return err
	}

	return nil
}

func (m *PolicyManager) CheckPerUePolicyV2(ueScope policyAPI.Scope, policyObject *mho.PolicyData) bool {

	if policyObject.API.Scope.UeID == nil {
		return false
	}

	if *policyObject.API.Scope.UeID == "" {
		return false
	}

	if *policyObject.API.Scope.UeID != *ueScope.UeID {
		return false
	}

	if (policyObject.API.Scope.SliceID != nil) && (((policyObject.API.Scope.SliceID.SD == nil || (policyObject.API.Scope.SliceID.SD != nil && *policyObject.API.Scope.SliceID.SD == "")) ||
		policyObject.API.Scope.SliceID.Sst <= 0 || policyObject.API.Scope.SliceID.PlmnID.Mcc == "" || policyObject.API.Scope.SliceID.PlmnID.Mnc == "") ||
		(*policyObject.API.Scope.SliceID.SD != *ueScope.SliceID.SD || policyObject.API.Scope.SliceID.Sst != ueScope.SliceID.Sst ||
			policyObject.API.Scope.SliceID.PlmnID.Mcc != ueScope.SliceID.PlmnID.Mcc || policyObject.API.Scope.SliceID.PlmnID.Mnc != ueScope.SliceID.PlmnID.Mnc)) {
		return false
	}

	if (policyObject.API.Scope.QosID != nil) && ((policyObject.API.Scope.QosID.QcI == nil && policyObject.API.Scope.QosID.The5QI == nil) ||
		(policyObject.API.Scope.QosID.QcI != nil && *policyObject.API.Scope.QosID.QcI != *ueScope.QosID.QcI) ||
		(policyObject.API.Scope.QosID.The5QI != nil && *policyObject.API.Scope.QosID.The5QI != *ueScope.QosID.The5QI)) {
		return false
	}

	if (policyObject.API.Scope.CellID != nil) && (((policyObject.API.Scope.CellID.CID.NcI == nil && policyObject.API.Scope.CellID.CID.EcI == nil) ||
		(policyObject.API.Scope.CellID.CID.NcI != nil && *policyObject.API.Scope.CellID.CID.NcI != *ueScope.CellID.CID.NcI) ||
		(policyObject.API.Scope.CellID.CID.EcI != nil && *policyObject.API.Scope.CellID.CID.EcI != *ueScope.CellID.CID.EcI)) ||
		((policyObject.API.Scope.CellID.PlmnID.Mcc == "" || policyObject.API.Scope.CellID.PlmnID.Mnc == "") ||
			(policyObject.API.Scope.CellID.PlmnID.Mcc != ueScope.CellID.PlmnID.Mcc || policyObject.API.Scope.CellID.PlmnID.Mnc != ueScope.CellID.PlmnID.Mnc))) {
		return false
	}

	return true
}

func (m *PolicyManager) CheckPerSlicePolicyV2(ueScope policyAPI.Scope, policyObject *mho.PolicyData) bool {

	if policyObject.API.Scope.SliceID == nil {
		return false
	}

	if (policyObject.API.Scope.SliceID != nil && *policyObject.API.Scope.SliceID.SD == "") ||
		policyObject.API.Scope.SliceID.Sst <= 0 ||
		policyObject.API.Scope.SliceID.PlmnID.Mcc == "" ||
		policyObject.API.Scope.SliceID.PlmnID.Mnc == "" {
		return false
	}

	if (policyObject.API.Scope.UeID != nil) && !((*policyObject.API.Scope.UeID == "") || (*policyObject.API.Scope.UeID == *ueScope.UeID)) {
		return false
	}

	if (policyObject.API.Scope.QosID != nil) && ((policyObject.API.Scope.QosID.QcI == nil && policyObject.API.Scope.QosID.The5QI == nil) ||
		(policyObject.API.Scope.QosID.QcI != nil && *policyObject.API.Scope.QosID.QcI != *ueScope.QosID.QcI) ||
		(policyObject.API.Scope.QosID.The5QI != nil && *policyObject.API.Scope.QosID.The5QI != *ueScope.QosID.The5QI)) {
		return false
	}

	if (policyObject.API.Scope.CellID != nil) && (((policyObject.API.Scope.CellID.CID.NcI == nil && policyObject.API.Scope.CellID.CID.EcI == nil) ||
		(policyObject.API.Scope.CellID.CID.NcI != nil && *policyObject.API.Scope.CellID.CID.NcI != *ueScope.CellID.CID.NcI) ||
		(policyObject.API.Scope.CellID.CID.EcI != nil && *policyObject.API.Scope.CellID.CID.EcI != *ueScope.CellID.CID.EcI)) ||
		((policyObject.API.Scope.CellID.PlmnID.Mcc == "" || policyObject.API.Scope.CellID.PlmnID.Mnc == "") ||
			(policyObject.API.Scope.CellID.PlmnID.Mcc != ueScope.CellID.PlmnID.Mcc || policyObject.API.Scope.CellID.PlmnID.Mnc != ueScope.CellID.PlmnID.Mnc))) {
		return false
	}

	return true
}

func (m *PolicyManager) GetTsResultForUEV2(ueScope policyAPI.Scope, rsrps []int, cellIds []policyAPI.CellID) policyAPI.CellID {

	var bestCell policyAPI.CellID
	bestScore := -math.MaxFloat64
	for i := 0; i < len(rsrps); i++ {
		preferece := m.GetPreferenceV2(ueScope, cellIds[i])
		score := m.GetPreferenceScoresV2(preferece, rsrps[i])
		if score > bestScore {
			bestCell = cellIds[i]
			bestScore = score
		}
	}
	return bestCell
}

func (m *PolicyManager) GetPreferenceScoresV2(preference string, rsrp int) float64 {
	return float64(rsrp) + float64(m.preferenceMap[preference])
}

func (m *PolicyManager) GetPreferenceV2(ueScope policyAPI.Scope, queryCellId policyAPI.CellID) string {

	var preference string = "DEFAULT"
	for _, policy := range *m.policyMap {
		if policy.IsEnforced {
			if m.CheckPerSlicePolicyV2(ueScope, policy) || m.CheckPerUePolicyV2(ueScope, policy) {
				for _, tspResource := range policy.API.TSPResources {

					for _, cellId := range tspResource.CellIDList {
						if ((cellId.CID.NcI != nil && queryCellId.CID.NcI != nil && *cellId.CID.NcI == *queryCellId.CID.NcI) ||
							(cellId.CID.EcI != nil && queryCellId.CID.EcI != nil && *cellId.CID.EcI == *queryCellId.CID.EcI)) &&
							(cellId.PlmnID.Mcc == queryCellId.PlmnID.Mcc && cellId.PlmnID.Mnc == queryCellId.PlmnID.Mnc) {
							preference = string(tspResource.Preference)
						}
					}
				}
			}
		}
	}
	return preference
}

func (m *PolicyManager) AddPolicyV2(policyId string, policyDir string, policyObject *mho.PolicyData) (*mho.PolicyData, error) {

	policyPath := policyDir + policyId
	err := m.ReadPolicyObjectFromFileV2(policyPath, policyObject)
	if err != nil {
		log.Error(fmt.Sprintf("Couldn't read PolicyObject from file \n policyId: %s from: %s", policyId, policyPath))
		return nil, err
	}
	return policyObject, nil
}

func (m *PolicyManager) EnforcePolicyV2(policyId string) bool {

	if _, ok := (*m.policyMap)[policyId]; ok {
		(*m.policyMap)[policyId].IsEnforced = true
		return true
	}
	log.Error(fmt.Sprintf("Policy with policyId: %s, not enforced", policyId))
	return false
}

func (m *PolicyManager) DisablePolicyV2(policyId string) bool {

	if _, ok := (*m.policyMap)[policyId]; ok {
		(*m.policyMap)[policyId].IsEnforced = false
		return true
	}
	log.Error(fmt.Sprintf("Policy with policyId: %s, not enforced", policyId))
	return false
}

func (m *PolicyManager) GetPolicyV2(policyId string) (*mho.PolicyData, bool) {

	if val, ok := (*m.policyMap)[policyId]; ok {
		return val, ok
	}
	log.Error(fmt.Sprintf("Policy with policyId: %s, not enforced", policyId))
	return nil, false
}

func (m *PolicyManager) ValidatePolicyJsonSchemaV2(jsonPath string) (bool, error) {

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + m.validator.schemePath)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + jsonPath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (m *PolicyManager) UnmarshalPolicyJsonV2(jsonFile []byte, policyObject *mho.PolicyData) error {

	if err := json.Unmarshal(jsonFile, policyObject.API); err != nil {
		log.Error("Couldn't read PolicyObject from file")
		return err
	}
	policyObject.IsEnforced = false
	return nil

}

func (m *PolicyManager) LoadPolicyJsonFromFileV2(path string) ([]byte, error) {

	jsonFile, err := os.Open(path)
	if err != nil {
		log.Error("Failed to open policy JSON File")
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Error("Failed to read data from policy JSON File")
		return nil, err
	}
	return byteValue, nil

}

func (m *PolicyManager) isSimilarEnforced(policyData *mho.PolicyData) bool {
	for _, policy := range *m.policyMap {

		sameSlice := false
		sameUE := false
		sameQoS := false
		sameCellID := false

		if (policyData.API.Scope.SliceID == nil && policy.API.Scope.SliceID == nil) ||
			(policyData.API.Scope.SliceID != nil && policy.API.Scope.SliceID != nil &&
				policy.API.Scope.SliceID.Sst == policyData.API.Scope.SliceID.Sst &&
				policy.API.Scope.SliceID.SD != nil && policyData.API.Scope.SliceID.SD != nil &&
				*policy.API.Scope.SliceID.SD == *policyData.API.Scope.SliceID.SD &&
				policy.API.Scope.SliceID.PlmnID.Mcc == policyData.API.Scope.SliceID.PlmnID.Mcc &&
				policy.API.Scope.SliceID.PlmnID.Mnc == policyData.API.Scope.SliceID.PlmnID.Mnc) {
			sameSlice = true
		}

		if (policyData.API.Scope.UeID == nil && policy.API.Scope.UeID == nil) ||
			(policyData.API.Scope.UeID != nil && policy.API.Scope.UeID != nil && *policy.API.Scope.UeID == *policyData.API.Scope.UeID) {
			sameUE = true
		}

		if (policyData.API.Scope.QosID == nil && policy.API.Scope.QosID == nil) ||
			(policyData.API.Scope.QosID != nil && policy.API.Scope.QosID != nil &&
				((policy.API.Scope.QosID.QcI != nil && policyData.API.Scope.QosID.QcI != nil &&
					*policy.API.Scope.QosID.QcI == *policyData.API.Scope.QosID.QcI) ||
					(policy.API.Scope.QosID.The5QI != nil && policyData.API.Scope.QosID.The5QI != nil &&
						*policy.API.Scope.QosID.The5QI == *policyData.API.Scope.QosID.The5QI))) {
			sameQoS = true
		}

		if (policyData.API.Scope.CellID == nil && policy.API.Scope.CellID == nil) ||
			(policyData.API.Scope.CellID != nil && policy.API.Scope.CellID != nil &&
				((policy.API.Scope.CellID.CID.NcI != nil && policyData.API.Scope.CellID.CID.NcI != nil &&
					*policy.API.Scope.CellID.CID.NcI == *policyData.API.Scope.CellID.CID.NcI) ||
					(policy.API.Scope.CellID.CID.EcI != nil && policyData.API.Scope.CellID.CID.EcI != nil &&
						*policy.API.Scope.CellID.CID.EcI == *policyData.API.Scope.CellID.CID.EcI)) &&
				policy.API.Scope.CellID.PlmnID.Mcc == policyData.API.Scope.CellID.PlmnID.Mcc &&
				policy.API.Scope.CellID.PlmnID.Mnc == policyData.API.Scope.CellID.PlmnID.Mnc) {
			sameCellID = true
		}

		if sameSlice {
			if sameUE {
				if sameQoS {
					if !sameCellID {
						continue
					}
				} else {
					continue
				}
			} else {
				if sameQoS {
					if !sameCellID {
						continue
					}
				} else {
					continue
				}
			}
		} else {
			if sameUE {
				if sameQoS {
					if !sameCellID {
						continue
					}
				} else {
					continue
				}
			} else {
				continue
			}
		}

		if policy.IsEnforced {
			policy.IsEnforced = false
			return true
		}

	}
	return false
}
