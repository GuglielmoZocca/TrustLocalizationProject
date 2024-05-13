/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//File that describes functions that read from ledger

package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ReadTarget reads the target information from collection
func (s *PositionContract) ReadTarget(ctx contractapi.TransactionContextInterface, collection string) (*TargetInfo, error) {

	//Check that is it the right collection
	if len(collection) < 7 {
		return nil, fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Target" {
		return nil, fmt.Errorf("this isn't the right collection:" + collection)
	}

	// Get target data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: %v", err)
	}

	// target properties are private, therefore they get passed in transient field
	transientTargetJSON, ok := transientMap["target_data"]
	if !ok {
		return nil, fmt.Errorf("target not found in the transient map")
	}

	type targetData struct {
		ID string `json:"TargetID"`
	}

	//Unmarshal transient data
	var targetDataInput targetData
	err = json.Unmarshal(transientTargetJSON, &targetDataInput)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check transient data
	if len(targetDataInput.ID) == 0 {
		return nil, fmt.Errorf("targetID field must be a non-empty string")
	}

	// Verify that the client is submitting request to peer in their organization
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadTarget cannot be performed: Error %v", err)
	}

	//Gat target data
	log.Printf("ReadTarget: collection %v, ID %v", collection, targetDataInput.ID)
	targetJSON, err := ctx.GetStub().GetPrivateData(collection, targetDataInput.ID) //get the target from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read target: %v", err)
	}

	// No Target found, return empty response
	if targetJSON == nil {
		log.Printf("%v does not exist in collection %v", targetDataInput, collection)
		return nil, nil
	}

	//Unmarshal target data
	var target *TargetInfo
	err = json.Unmarshal(targetJSON, &target)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return target, nil

}

// ReadAllIDDevices reads all devices id information from collection
func (s *PositionContract) ReadAllIDDevices(ctx contractapi.TransactionContextInterface, collection string) ([]string, error) {

	//Check that is it the right collection
	if len(collection) < 7 {
		return nil, fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]

	if check != "Device" {
		return nil, fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return nil, fmt.Errorf("client is not an admin")
	}

	// Verify that the client is submitting request to peer in their organization
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadDevice cannot be performed: Error %v", err)
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collection, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []string

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset *Device
		err = json.Unmarshal(response.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		results = append(results, asset.ID)
	}

	return results, nil

}

// ReadDevice reads the device information from collection
func (s *PositionContract) ReadDevice(ctx contractapi.TransactionContextInterface, collection string) (*Device, error) {

	//Check that is it the right collection
	if len(collection) < 7 {
		return nil, fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Device" {
		return nil, fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return nil, fmt.Errorf("client is not an admin")
	}

	// Get device data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: %v", err)
	}

	// device properties are private, therefore they get passed in transient field
	transientDeviceJSON, ok := transientMap["device_data"]
	if !ok {
		return nil, fmt.Errorf("device not found in the transient map")
	}

	type deviceData struct {
		ID string `json:"deviceID"`
	}

	//Unmarshal transient data
	var deviceDataInput deviceData
	err = json.Unmarshal(transientDeviceJSON, &deviceDataInput)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check transient data
	if len(deviceDataInput.ID) == 0 {
		return nil, fmt.Errorf("deviceID field must be a non-empty string")
	}

	// Verify that the client is submitting request to peer in their organization
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadDevice cannot be performed: Error %v", err)
	}

	//Gat device data
	log.Printf("ReadDevice: collection %v, ID %v", collection, deviceDataInput.ID)
	deviceJSON, err := ctx.GetStub().GetPrivateData(collection, deviceDataInput.ID) //get the device from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read device: %v", err)
	}

	// No Device found
	if deviceJSON == nil {
		return nil, fmt.Errorf("%v does not exist in collection %v", deviceDataInput.ID, collection)
	}

	//Unmarshal device data
	var device *Device
	err = json.Unmarshal(deviceJSON, &device)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return device, nil

}
