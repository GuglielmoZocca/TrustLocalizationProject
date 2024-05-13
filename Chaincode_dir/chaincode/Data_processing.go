/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//File that describes functions that write on ledger

package chaincode

import (
	"cmp"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"math"
	"slices"
	"time"
)

// PositionContract of this fabric
type PositionContract struct {
	contractapi.Contract
}

// Device describes devices in this network
type Device struct {
	ID    string   `json:"deviceID"`     //ID of the device
	Key   string   `json:"Key""`         //Key of the decription of data from the device
	X     float32  `json:"coordinateX"`  //X coordinate of the device
	Y     float32  `json:"coordinateY"`  //Y coordinate of the device
	Obs   int      `json:"observation"`  //Obs distance observed from the device
	Conf  float32  `json:"Confidence"`   //Conf confidence of the device
	Ev    float32  `json:"Evidence"`     //Ev evidence of the device
	Rep   int      `json:"Reputation"`   //Rep reputation of the device
	Trust float32  `json:"Trust"`        //Trust of the device
	Neigh []string `json:"Neighborhood"` //Neighborhood id of the devices near the device
}

// TargetInfo detail of a certain target
type TargetInfo struct {
	ID   string  `json:"TargetID"`    //ID of the target
	X    float32 `json:"coordinateX"` //X coordinate of the target
	Y    float32 `json:"coordinateY"` //Y coordinate of the target
	Date string  `json:"Date"`        //Timestamp of the update
	Upd  bool    `json:"Update"`      //Indicate if target is updated
}

// CreateDevice creates a device object in a certain collection
func (s *PositionContract) CreateDevice(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Device" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get new device from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Device properties are private, therefore they get passed in transient field, instead of func args
	transientDeviceJSON, ok := transientMap["Device_properties"]
	if !ok {
		// log error to stdout
		return fmt.Errorf("Device not found in the transient map input")
	}

	type deviceTransientInput struct {
		ID    string   `json:"deviceID"`
		X     float32  `json:"coordinateX"`
		Y     float32  `json:"coordinateY"`
		Key   string   `json:"Key"`
		Neigh []string `json:"Neighborhood"`
		Rep   int      `json:"Reputation"`
	}

	//Unmarshal the transient data
	var deviceInput deviceTransientInput
	err = json.Unmarshal(transientDeviceJSON, &deviceInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check the transient data
	if len(deviceInput.ID) == 0 {
		return fmt.Errorf("deviceID field must be a non-empty string")
	}
	if len(deviceInput.Neigh) <= 1 {
		return fmt.Errorf("Device must have more than 1 neighbour")
	}
	if deviceInput.Rep <= 0 {
		return fmt.Errorf("Device reputation must be positive")
	}
	for _, n := range deviceInput.Neigh {
		if len(n) == 0 {
			return fmt.Errorf("Neighbour id field must be a non-empty string")
		}
	}

	// Check if Device already exists
	deviceAsBytes, err := ctx.GetStub().GetPrivateData(collection, deviceInput.ID)

	if err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	} else if deviceAsBytes != nil {
		fmt.Println("Device already exists: " + deviceInput.ID)
		return fmt.Errorf("this device already exists: " + deviceInput.ID)
	}

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("CreateDevice cannot be performed: Error %v", err)
	}

	// Make a new device
	device := Device{
		ID:    deviceInput.ID,
		Key:   deviceInput.Key,
		X:     deviceInput.X,
		Y:     deviceInput.Y,
		Obs:   0,
		Conf:  0,
		Ev:    0,
		Rep:   deviceInput.Rep,
		Trust: 0,
		Neigh: deviceInput.Neigh,
	}
	deviceJSONasBytes, err := json.Marshal(device)
	if err != nil {
		return fmt.Errorf("failed to marshal device into JSON: %v", err)
	}

	//Put the data in the organization
	err = ctx.GetStub().PutPrivateData(collection, deviceInput.ID, deviceJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put device into private data collecton: %v", err)
	}

	return nil
}

// UpdateDevice update the id, key and neigh of device object in a certain collection
func (s *PositionContract) UpdateDevice(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Device" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get new device from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Device properties are private, therefore they get passed in transient field, instead of func args
	transientDeviceJSON, ok := transientMap["Device_properties"]
	if !ok {
		// log error to stdout
		return fmt.Errorf("Device not found in the transient map input")
	}

	type deviceTransientInput struct {
		ID    string   `json:"deviceID"`
		Key   string   `json:"Key"`
		Neigh []string `json:"Neighborhood"`
	}

	//Unmarshal the transient data
	var deviceInput deviceTransientInput
	err = json.Unmarshal(transientDeviceJSON, &deviceInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check the transient data
	if len(deviceInput.ID) == 0 {
		return fmt.Errorf("deviceID field must be a non-empty string")
	}
	if len(deviceInput.Neigh) <= 1 {
		return fmt.Errorf("Device must have more than 1 neighbour")
	}
	for _, n := range deviceInput.Neigh {
		if len(n) == 0 {
			return fmt.Errorf("Neighbour id field must be a non-empty string")
		}
	}

	// Check if Device already exists
	DeviceAsBytes, err := ctx.GetStub().GetPrivateData(collection, deviceInput.ID)
	if err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	} else if DeviceAsBytes == nil {
		fmt.Println("Device not exist: " + deviceInput.ID)
		return fmt.Errorf("this device not exist: " + deviceInput.ID)
	}

	//Unmarshal the device data
	var deviceData Device
	err = json.Unmarshal(DeviceAsBytes, &deviceData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("CreateDevice cannot be performed: Error %v", err)
	}

	// Make a new device
	device := Device{
		ID:    deviceInput.ID,
		Key:   deviceInput.Key,
		X:     deviceData.X,
		Y:     deviceData.Y,
		Obs:   deviceData.Obs,
		Conf:  deviceData.Conf,
		Ev:    deviceData.Ev,
		Rep:   deviceData.Rep,
		Trust: deviceData.Trust,
		Neigh: deviceInput.Neigh,
	}
	deviceJSONasBytes, err := json.Marshal(device)
	if err != nil {
		return fmt.Errorf("failed to marshal device into JSON: %v", err)
	}

	//Put the data in the organization
	err = ctx.GetStub().PutPrivateData(collection, deviceInput.ID, deviceJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put device into private data collecton: %v", err)
	}

	return nil
}

// CreateTarget creates a target object in acertain collection
func (s *PositionContract) CreateTarget(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Target" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get new target from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Target properties are private, therefore they get passed in transient field, instead of func args
	transientTargetJSON, ok := transientMap["Target_properties"]
	if !ok {
		// log error to stdout
		return fmt.Errorf("Target not found in the transient map input")
	}

	type targetTransientInput struct {
		ID string `json:"TargetID"`
	}

	//Unmarshal the transient data
	var targetInput targetTransientInput
	err = json.Unmarshal(transientTargetJSON, &targetInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check the transient data
	if len(targetInput.ID) == 0 {
		return fmt.Errorf("targetID must be a non-empty string")
	}

	// Check if Target already exists
	TargetAsBytes, err := ctx.GetStub().GetPrivateData(collection, targetInput.ID)
	if err != nil {
		return fmt.Errorf("failed to get target: %v", err)
	} else if TargetAsBytes != nil {
		fmt.Println("Target already exists: " + targetInput.ID)
		return fmt.Errorf("this target already exists: " + targetInput.ID)
	}

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("CreateTarget cannot be performed: Error %v", err)
	}

	// Make a target
	target := TargetInfo{
		ID:   targetInput.ID,
		X:    0,
		Y:    0,
		Date: "",
		Upd:  false,
	}
	targetJSONasBytes, err := json.Marshal(target)
	if err != nil {
		return fmt.Errorf("failed to marshal target into JSON: %v", err)
	}

	//Put the data in the collection
	err = ctx.GetStub().PutPrivateData(collection, targetInput.ID, targetJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put target into private data collecton: %v", err)
	}

	return nil
}

// UpdateDeviceObsConf updates a device confidence and observation in a certain collection
func (s *PositionContract) UpdateDeviceObsConf(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Device" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get device data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Device properties are private, therefore they get passed in transient field, instead of func args
	transientDeviceJSON, ok := transientMap["Device_properties"]
	if !ok {
		// log error to stdout
		return fmt.Errorf("Device not found in the transient map input")
	}

	type deviceTransientInput struct {
		ID      string  `json:"deviceID"`
		Obs     int     `json:"observation"`
		Conf    float32 `json:"Confidence"`
		MinConf float32 `json:"MinConf"` //Minimum confidence of the device
		MaxConf float32 `json:"MaxConf"` //Maximum confidence of the device
	}

	//unmarshal the transient data
	var deviceInput deviceTransientInput
	err = json.Unmarshal(transientDeviceJSON, &deviceInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check the transient data
	if len(deviceInput.ID) == 0 {
		return fmt.Errorf("deviceID field must be a non-empty string")
	}

	//var Con float32 //minimum confidence
	//var max float32 //maximum confidence
	//minConf = 0.4
	//maxConf = 1

	if deviceInput.Conf > deviceInput.MaxConf || deviceInput.Conf < deviceInput.MinConf {
		return fmt.Errorf("Confidence must be between certain values")
	}

	// Check if Device already exists
	DeviceAsBytes, err := ctx.GetStub().GetPrivateData(collection, deviceInput.ID)
	if err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	} else if DeviceAsBytes == nil {
		fmt.Println("Device not exist: " + deviceInput.ID)
		return fmt.Errorf("this device not exist: " + deviceInput.ID)
	}

	//Unmarshal the device data
	var deviceData Device
	err = json.Unmarshal(DeviceAsBytes, &deviceData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("UpdateDevice cannot be performed: Error %v", err)
	}

	// Make a deiice
	device := Device{
		ID:    deviceInput.ID,
		Key:   deviceData.Key,
		X:     deviceData.X,
		Y:     deviceData.Y,
		Obs:   deviceInput.Obs,
		Conf:  deviceInput.Conf,
		Ev:    deviceData.Ev,
		Rep:   deviceData.Rep,
		Trust: deviceData.Trust,
		Neigh: deviceData.Neigh,
	}
	deviceJSONasBytes, err := json.Marshal(device)
	if err != nil {
		return fmt.Errorf("failed to marshal device into JSON: %v", err)
	}

	//Put the device in the collection
	err = ctx.GetStub().PutPrivateData(collection, deviceInput.ID, deviceJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put device into private data collecton: %v", err)
	}

	return nil
}

// Point represents a point in space.
type Point struct {
	X float32
	Y float32
}

// New returns a Point based on X and Y positions on a graph.
func New(x float32, y float32) Point {
	return Point{x, y}
}

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func Distance(p1 Point, p2 Point) float64 {
	first := math.Pow(float64(p2.X-p1.X), 2)
	second := math.Pow(float64(p2.Y-p1.Y), 2)
	return math.Sqrt(first + second)
}

// UpdateDeviceObsConf updates evidence, reputation and trust in a certain collection
func (s *PositionContract) UpdateDeviceEvRep(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Device" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get device data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Device properties are private, therefore they get passed in transient field, instead of func args
	transientDeviceJSON, ok := transientMap["Device_properties"]
	if !ok {
		// log error to stdout
		return fmt.Errorf("Device not found in the transient map input")
	}

	type deviceTransientInput struct {
		ID          string  `json:"deviceID"`
		PRH         int     `json:"PRH"`         //Reward or penalty in case of high evidence and high confidence, and in case of low evidence and high confidence respectively
		PRL         int     `json:"PRL"`         //Reward or penalty in case of high evidence and low confidence, and in case of low evidence and low confidence respectively
		ThreashConf float32 `json:"threashConf"` //threshold that indicate when a confidence value is high or low
		ThreashEv   float32 `json:"threashEv"`   //threshold that indicate when an evidence value is high or low
		MaxRep      int     `json:"maxRep"`      //max reputation
	}

	//unmarshal the transient data
	var deviceInput deviceTransientInput
	err = json.Unmarshal(transientDeviceJSON, &deviceInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check the transient data
	if len(deviceInput.ID) == 0 {
		return fmt.Errorf("deviceID field must be more than 0")
	}
	if deviceInput.MaxRep <= 0 {
		return fmt.Errorf("deviceID field must be a value grater than 0")
	}

	// Check if Device already exists
	DeviceAsBytes, err := ctx.GetStub().GetPrivateData(collection, deviceInput.ID)
	if err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	} else if DeviceAsBytes == nil {
		fmt.Println("Device not exist: " + deviceInput.ID)
		return fmt.Errorf("this device not exist: " + deviceInput.ID)
	}

	//Unmarshal the device data
	var deviceData Device
	err = json.Unmarshal(DeviceAsBytes, &deviceData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("UpdateDevice cannot be performed: Error %v", err)
	}

	var Evi float32 //Evidence of device

	//var threasDiff int //threashold to indicate that a observation support another
	//threasDiff = 1

	var num int //number of neighbour
	var tmp float32
	num = 0
	tmp = 0
	for _, idN := range deviceData.Neigh {
		num = num + 1
		// Check if Device already exists
		DeviceAsBytes, err = ctx.GetStub().GetPrivateData(collection, idN)
		if err != nil {
			return fmt.Errorf("failed to get device: %v", err)
		} else if DeviceAsBytes == nil {
			fmt.Println("Device not exist: " + idN)
			return fmt.Errorf("this device not exist: " + idN)
		}

		//Unmarshal the device data
		var deviceNeigh Device
		err = json.Unmarshal(DeviceAsBytes, &deviceNeigh)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		//Check triangle inequality
		if (float64(deviceData.Obs+deviceNeigh.Obs) >= Distance(New(deviceData.X, deviceData.Y), (New(deviceNeigh.X, deviceNeigh.Y)))*1000) && ((Distance(New(deviceData.X, deviceData.Y), (New(deviceNeigh.X, deviceNeigh.Y)))*1000 + float64(deviceNeigh.Obs)) >= float64(deviceData.Obs)) {
			tmp = tmp + deviceNeigh.Conf
		} else {
			tmp = tmp - deviceNeigh.Conf
		}
	}
	//Calculate the evidence of the device
	Evi = tmp * (float32(1) / float32(num))

	//PRH = 2
	//PRL = 1

	//threashConf = 0.7
	//threashEv = 0

	var Repu int //Reputation of device
	Repu = deviceData.Rep

	//Update the reputation of the device
	if deviceData.Conf >= deviceInput.ThreashConf && Evi >= deviceInput.ThreashEv {
		Repu = Repu + deviceInput.PRH
		if Repu > deviceInput.MaxRep {
			Repu = deviceInput.MaxRep
		}
	}
	if deviceData.Conf < deviceInput.ThreashConf && Evi >= deviceInput.ThreashEv {
		Repu = Repu + deviceInput.PRL
		if Repu > deviceInput.MaxRep {
			Repu = deviceInput.MaxRep
		}
	}
	if deviceData.Conf >= deviceInput.ThreashConf && Evi < deviceInput.ThreashEv {
		Repu = Repu - (deviceInput.PRH + 1)
		if Repu < deviceInput.MaxRep {
			Repu = deviceInput.MaxRep
		}
	}
	if deviceData.Conf < deviceInput.ThreashConf && Evi < deviceInput.ThreashEv {
		Repu = Repu - (deviceInput.PRL + 1)
		if Repu < deviceInput.MaxRep {
			Repu = deviceInput.MaxRep
		}
	}

	var Trustv float32 //Trust value to the observation

	//Calcutate the trust of the device
	Trustv = deviceData.Conf * Evi * float32(Repu)

	// Make a device
	device := Device{
		ID:    deviceInput.ID,
		Key:   deviceData.Key,
		X:     deviceData.X,
		Y:     deviceData.Y,
		Obs:   deviceData.Obs,
		Conf:  deviceData.Conf,
		Ev:    Evi,
		Rep:   Repu,
		Trust: Trustv,
		Neigh: deviceData.Neigh,
	}
	deviceJSONasBytes, err := json.Marshal(device)
	if err != nil {
		return fmt.Errorf("failed to marshal device into JSON: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(collection, deviceInput.ID, deviceJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put device into private data collecton: %v", err)
	}

	return nil
}

// DeleteDevice can be used to delete a device from a nextwork
func (s *PositionContract) DeleteDevice(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Device" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get device data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: %v", err)
	}

	// device properties are private, therefore they get passed in transient field
	transientDeleteJSON, ok := transientMap["device_delete"]
	if !ok {
		return fmt.Errorf("device to delete not found in the transient map")
	}

	type deviceDelete struct {
		ID string `json:"deviceID"`
	}

	//Unmarshal transient data
	var deviceDeleteInput deviceDelete
	err = json.Unmarshal(transientDeleteJSON, &deviceDeleteInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check transient data
	if len(deviceDeleteInput.ID) == 0 {
		return fmt.Errorf("deviceID field must be a non-empty string")
	}

	// Verify that the client is submitting request to peer in their organization
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("DeleteDevice cannot be performed: Error %v", err)
	}

	//Check if the device exists
	log.Printf("Deleting device: %v", deviceDeleteInput.ID)
	valAsbytes, err := ctx.GetStub().GetPrivateData(collection, deviceDeleteInput.ID) //get the device from chaincode state
	if err != nil {
		return fmt.Errorf("failed to read device: %v", err)
	}
	if valAsbytes == nil {
		return fmt.Errorf("device not found: %v", deviceDeleteInput.ID)
	}

	// delete the device from the collection
	err = ctx.GetStub().DelPrivateData(collection, deviceDeleteInput.ID)
	if err != nil {
		return fmt.Errorf("failed to delete state: %v", err)
	}

	return nil

}

// DeleteTarget can be used to delete a target from a nextwork
func (s *PositionContract) DeleteTarget(ctx contractapi.TransactionContextInterface, collection string) error {

	//Check that is it the right collection
	if len(collection) < 7 {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}
	check := collection[:6]
	if check != "Target" {
		return fmt.Errorf("this isn't the right collection:" + collection)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get target data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: %v", err)
	}

	// target properties are private, therefore they get passed in transient field
	transientDeleteJSON, ok := transientMap["target_delete"]
	if !ok {
		return fmt.Errorf("target to delete not found in the transient map")
	}

	type targetDelete struct {
		ID string `json:"targetID"`
	}

	//Unmarshal transient data
	var targetDeleteInput targetDelete
	err = json.Unmarshal(transientDeleteJSON, &targetDeleteInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check transient data
	if len(targetDeleteInput.ID) == 0 {
		return fmt.Errorf("targetID field must be a non-empty string")
	}

	// Verify that the client is submitting request to peer in their organization
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("DeleteTarget cannot be performed: Error %v", err)
	}

	//Check if the target exists
	log.Printf("Deleting target: %v", targetDeleteInput.ID)
	valAsbytes, err := ctx.GetStub().GetPrivateData(collection, targetDeleteInput.ID) //get the target from chaincode state
	if err != nil {
		return fmt.Errorf("failed to read target: %v", err)
	}
	if valAsbytes == nil {
		return fmt.Errorf("target not found: %v", targetDeleteInput.ID)
	}

	// delete the target from the collection
	err = ctx.GetStub().DelPrivateData(collection, targetDeleteInput.ID)
	if err != nil {
		return fmt.Errorf("failed to delete state: %v", err)
	}

	return nil

}

// Calculation position of a target
func (s *PositionContract) PositionTarget(ctx contractapi.TransactionContextInterface, collectionD string, collectionT string, dateU string) error {

	//Check that is it the right collection
	if len(collectionD) < 7 {
		return fmt.Errorf("this isn't the right device collection:" + collectionD)
	}
	check1 := collectionD[:6]
	if check1 != "Device" {
		return fmt.Errorf("this isn't the right device collection:" + collectionD)
	}

	//Check that is it the right collection
	if len(collectionT) < 7 {
		return fmt.Errorf("this isn't the right target collection:" + collectionT)
	}
	check2 := collectionT[:6]
	if check2 != "Target" {
		return fmt.Errorf("this isn't the right target collection:" + collectionT)
	}

	//Check if the client is an admin
	ident, err := cid.HasOUValue(ctx.GetStub(), "admin")
	if !ident {
		return fmt.Errorf("client is not an admin")
	}

	// Get target data from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: %v", err)
	}

	// target properties are private, therefore they get passed in transient field
	transientTargetJSON, ok := transientMap["target_Posisiton"]
	if !ok {
		return fmt.Errorf("target not found in the transient map")
	}

	type targetPosisiton struct {
		ID      string   `json:"TargetID"`
		Thresh  float64  `json:"ThreshErr"`
		Devices []string `json:"DevicesUp"`
	}

	//Unmarshal transient data
	var targetPosisitonInput targetPosisiton
	err = json.Unmarshal(transientTargetJSON, &targetPosisitonInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//Check transient data
	if len(targetPosisitonInput.ID) == 0 {
		return fmt.Errorf("targetID field must be a non-empty string")
	}

	idTarget := targetPosisitonInput.ID
	devices := targetPosisitonInput.Devices
	tresh := targetPosisitonInput.Thresh

	// Check if target already exists
	TargetAsBytes, err := ctx.GetStub().GetPrivateData(collectionT, idTarget)
	if err != nil {
		return fmt.Errorf("failed to get target: %v", err)
	} else if TargetAsBytes == nil {
		fmt.Println("Target not exist: " + idTarget)
		return fmt.Errorf("Target not exist: " + idTarget)
	}

	//Struct to describe the distance observeted from a device between it and the target
	type devicesDist struct {
		Id       string
		Distance int
		Trust    float32
		X        float32
		Y        float32
	}

	var devicesDistcol []devicesDist

	for _, id := range devices {

		//Get device from the collection
		log.Printf("ReadAsset: collection %v, ID %v", collectionD, id)
		deviceJSON, err := ctx.GetStub().GetPrivateData(collectionD, id) //get the device from chaincode state
		if err != nil {
			fmt.Errorf("failed to read device: %v", err)
			return nil
		}

		// No Device with data id
		if deviceJSON == nil {
			return fmt.Errorf("%v does not exist in collection %v", id, collectionD)
		}

		//Unmarshal the device data
		var device *Device
		err = json.Unmarshal(deviceJSON, &device)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		devicesDistcol = append(devicesDistcol, devicesDist{device.ID, device.Obs, device.Trust, device.X, device.Y})

	}

	// Stable sort by trust and descending order
	slices.SortStableFunc(devicesDistcol, func(a devicesDist, b devicesDist) int {
		return cmp.Compare(b.Trust, a.Trust)
	})

	//take only 3 devices wiht higher trust

	var d1 devicesDist
	var d2 devicesDist
	var d3 devicesDist

	for i, d := range devicesDistcol {
		if i == 0 {
			d1 = d
		}
		if i == 1 {
			d2 = d
		}
		if i == 2 {
			d3 = d
		}
	}

	x1 := float64(d1.X)                        //Coordinate x
	y1 := float64(d1.Y)                        //Coordinate y
	x2 := float64(d2.X)                        //Coordinate x
	y2 := float64(d2.Y)                        //Coordinate y
	x3 := float64(d3.X)                        //Coordinate x
	y3 := float64(d3.Y)                        //Coordinate y
	r1 := float64(d1.Distance) / float64(1000) //Distance observed
	r2 := float64(d2.Distance) / float64(1000) //Distance observed
	r3 := float64(d3.Distance) / float64(1000) //Distance observed

	var xtarg float64 //Coordinate x of the target
	var ytarg float64 //Coordinate y of the target

	// Calculate the distance between the centers of the circles
	d := math.Sqrt(math.Pow(x1-x2, float64(2)) + math.Pow(y1-y2, float64(2)))

	// Check if circles do not intersect
	if d > r1+r2 || d < math.Abs(r1-r2) {
		return fmt.Errorf("no solution")
	}

	// Check if circles are tangent to each other
	if d == r1+r2 || d == math.Abs(r1-r2) {
		fmt.Println("The circles are tangent to each other.")
		// Calculate the coordinates of the tangent point using the midpoint formula
		xtarg = (x1 + x2) / float64(2)
		ytarg = (y1 + y2) / float64(2)

		// Make submitting client the owner
		target := TargetInfo{
			ID:   idTarget,
			X:    float32(xtarg),
			Y:    float32(ytarg),
			Date: time.Now().Format(time.DateTime),
			Upd:  true,
		}
		targetJSONasBytes, err := json.Marshal(target)
		if err != nil {
			return fmt.Errorf("failed to marshal target into JSON: %v", err)
		}

		//Put the target in the collection
		err = ctx.GetStub().PutPrivateData(collectionT, idTarget, targetJSONasBytes)
		if err != nil {
			return fmt.Errorf("failed to put target into private data collecton: %v", err)
		}
	} else {
		// Calculate the distance from center_p to the intersection line
		a := (r1*r1 - r2*r2 + d*d) / (float64(2) * d)
		// Calculate the distance from intersection point to the intersection line
		h := math.Sqrt(r1*r1 - a*a)
		// Calculate the intersection points
		xposs1 := x1 + (a/d)*(x2-x1) + (h/d)*(y2-y1)
		yposs1 := y1 + (a/d)*(y2-y1) - (h/d)*(x2-x1)
		xposs2 := x1 + (a/d)*(x2-x1) - (h/d)*(y2-y1)
		yposs2 := y1 + (a/d)*(y2-y1) + (h/d)*(x2-x1)
		//Select the best solution
		minv1 := math.Abs((math.Pow(x3-xposs1, float64(2)) + math.Pow(y3-yposs1, float64(2))) - r3*r3)
		minv2 := math.Abs((math.Pow(x3-xposs2, float64(2)) + math.Pow(y3-yposs2, float64(2))) - r3*r3)
		if minv1 <= minv2 {
			if minv1 >= tresh {
				return fmt.Errorf("no solution")
			}
			xtarg = xposs1
			ytarg = yposs1
		} else {
			if minv2 >= tresh {
				return fmt.Errorf("no solution")
			}
			xtarg = xposs2
			ytarg = yposs2
		}

		// Verify that the client is submitting request to peer in their organization
		// This is to ensure that a client from another org doesn't attempt to read or
		// write private data from this peer.
		err = verifyClientOrgMatchesPeerOrg(ctx)
		if err != nil {
			return fmt.Errorf("UpdateTarget cannot be performed: Error %v", err)
		}

		// Make submitting client the owner
		target := TargetInfo{
			ID:   idTarget,
			X:    float32(xtarg),
			Y:    float32(ytarg),
			Date: dateU,
			Upd:  true,
		}
		targetJSONasBytes, err := json.Marshal(target)
		if err != nil {
			return fmt.Errorf("failed to marshal target into JSON: %v", err)
		}

		//Put the target in the collection
		err = ctx.GetStub().PutPrivateData(collectionT, idTarget, targetJSONasBytes)
		if err != nil {
			return fmt.Errorf("failed to put target into private data collecton: %v", err)
		}
	}

	return nil

}

// getCollectionName is an internal helper function to get collection of submitting client identity.
func getCollectionName(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get the MSP ID of submitting client identity
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get verified MSPID: %v", err)
	}

	// Create the collection name
	orgCollection := clientMSPID + "PrivateCollection"

	return orgCollection, nil
}

// verifyClientOrgMatchesPeerOrg is an internal function used verify client org id and matches peer org id.
func verifyClientOrgMatchesPeerOrg(ctx contractapi.TransactionContextInterface) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the peer's MSPID: %v", err)
	}

	if clientMSPID != peerMSPID {
		return fmt.Errorf("client from org %v is not authorized to read or write private data from an org %v peer", clientMSPID, peerMSPID)
	}

	return nil
}

func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}
