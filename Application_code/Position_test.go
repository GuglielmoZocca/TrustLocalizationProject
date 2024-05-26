// TEST APPLICATION
package main

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

//Test in the context to be in org 1

var contract *client.Contract
var gw *client.Gateway
var clientConnection *grpc.ClientConn

// initL initialize the connexion with the ledger
func initL(org int, clientA string) {
	var err error
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	if org == 1 && clientA == "A" {
		clientConnection = newGrpcConnection(tlsCertPathU1, gatewayPeerA1, peerEndpointA1)
		id := newIdentity(certPathA1, mspIDA1)
		sign := newSign(keyPathA1)
		gw, err = client.Connect(
			id,
			client.WithSign(sign),
			client.WithClientConnection(clientConnection),
			// Default timeouts for different gRPC calls
			client.WithEvaluateTimeout(5*time.Second),
			client.WithEndorseTimeout(15*time.Second),
			client.WithSubmitTimeout(5*time.Second),
			client.WithCommitStatusTimeout(1*time.Minute),
		)
		if err != nil {
			panic(err)
		}
		// Override default values for chaincode and channel name as they may differ in testing contexts.
		chaincodeName := "Position"
		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
			chaincodeName = ccname
		}

		channelName := "mychannel"
		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
			channelName = cname
		}
		//Get reference to the network and contract
		network := gw.GetNetwork(channelName)
		contract = network.GetContract(chaincodeName)

	}
	if org == 1 && clientA == "U" {
		clientConnection = newGrpcConnection(tlsCertPathU1, gatewayPeerA1, peerEndpointA1)
		id := newIdentity(certPathU1, mspIDU1)
		sign := newSign(keyPathU1)
		gw, err = client.Connect(
			id,
			client.WithSign(sign),
			client.WithClientConnection(clientConnection),
			// Default timeouts for different gRPC calls
			client.WithEvaluateTimeout(5*time.Second),
			client.WithEndorseTimeout(15*time.Second),
			client.WithSubmitTimeout(5*time.Second),
			client.WithCommitStatusTimeout(1*time.Minute),
		)
		if err != nil {
			panic(err)
		}
		// Override default values for chaincode and channel name as they may differ in testing contexts.
		chaincodeName := "Position"
		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
			chaincodeName = ccname
		}

		channelName := "mychannel"
		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
			channelName = cname
		}

		//Get reference to the network and contract
		network := gw.GetNetwork(channelName)
		contract = network.GetContract(chaincodeName)

	}
	if org == 2 && clientA == "A" {
		clientConnection = newGrpcConnection(tlsCertPathU2, gatewayPeerA2, peerEndpointA2)
		id := newIdentity(certPathA2, mspIDA2)
		sign := newSign(keyPathA2)
		gw, err = client.Connect(
			id,
			client.WithSign(sign),
			client.WithClientConnection(clientConnection),
			// Default timeouts for different gRPC calls
			client.WithEvaluateTimeout(5*time.Second),
			client.WithEndorseTimeout(15*time.Second),
			client.WithSubmitTimeout(5*time.Second),
			client.WithCommitStatusTimeout(1*time.Minute),
		)
		if err != nil {
			panic(err)
		}
		// Override default values for chaincode and channel name as they may differ in testing contexts.
		chaincodeName := "Position"
		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
			chaincodeName = ccname
		}

		channelName := "mychannel"
		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
			channelName = cname
		}

		//Get reference to the network and contract
		network := gw.GetNetwork(channelName)
		contract = network.GetContract(chaincodeName)

	}
	if org == 2 && clientA == "U" {
		clientConnection = newGrpcConnection(tlsCertPathU2, gatewayPeerA2, peerEndpointA2)
		id := newIdentity(certPathU2, mspIDU2)
		sign := newSign(keyPathU2)
		gw, err = client.Connect(
			id,
			client.WithSign(sign),
			client.WithClientConnection(clientConnection),
			// Default timeouts for different gRPC calls
			client.WithEvaluateTimeout(5*time.Second),
			client.WithEndorseTimeout(15*time.Second),
			client.WithSubmitTimeout(5*time.Second),
			client.WithCommitStatusTimeout(1*time.Minute),
		)
		if err != nil {
			panic(err)
		}
		// Override default values for chaincode and channel name as they may differ in testing contexts.
		chaincodeName := "Position"
		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
			chaincodeName = ccname
		}

		channelName := "mychannel"
		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
			channelName = cname
		}

		//Get reference to the network and contract
		network := gw.GetNetwork(channelName)
		contract = network.GetContract(chaincodeName)

	}

}

// TestInitAsAdminOrg1 test initialization ledger as Admin of organization 1
func TestInitAsAdminOrg1(t *testing.T) {

	initL(1, "A")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devices := []string{"3", "1", "2"}

	Deletedevices(contract, devices)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract) = %q, want  nil`, err)
	}

	err = Deletedevices(contract, devices)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}
}

// TestInitAsUserOrg1 test initialization ledger as User of organization 1
func TestInitAsUserOrg1(t *testing.T) {

	initL(1, "U")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devices := []string{"3", "1", "2"}

	Deletedevices(contract, devices)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err == nil {
		t.Fatalf(`You cannot create as user: %q, want not nil error`, err)
	}

	t.Logf("A user can’t access the devices collection:%q", err)

	err = Deletedevices(contract, devices)
	if err == nil {
		t.Fatalf(`You cannot create as user: %q, want not nil error`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err == nil {
		t.Fatalf(`You cannot create as user: %q, want not nil error`, err)
	}

}

// TestInitAsAdminOrg2 test initialization ledger as Admin of organization 2
func TestInitAsAdminOrg2(t *testing.T) {

	initL(2, "A")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devices := []string{"3", "1", "2"}

	Deletedevices(contract, devices)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err == nil {
		t.Fatalf(`You cannot create as user: %q, want not nil error`, err)
	}

	t.Logf("A client can’t access the devices collection of an organization that it doesn’t belong:%q", err)

	err = Deletedevices(contract, devices)
	if err == nil {
		t.Fatalf(`You cannot create as user: %q, want not nil error`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err == nil {
		t.Fatalf(`You cannot create as user: %q, want not nil error`, err)
	}

}

// TestAddObsAsAdminOrg1 test add observation of a device as Admin of organization 1
func TestAddObsAsAdminOrg1(t *testing.T) {

	initL(1, "A")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"3", "1", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract) = %q, want  nil`, err)
	}

	idDevice := "1"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice,idTarget) = %q, want nil`, err)
	}

	var deviceR Device

	err, deviceR = seeDataDevice(contract, idDevice)
	if err != nil {
		t.Fatalf(`seeDataDevice(contract, idDevice) = %q, want nil`, err)
	}

	if deviceR.Obs != 4242 {
		t.Fatalf(`it doesn't insert the correct obs, want %d: %q`, 4242, err)
	}

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}

// TestAddObsWithErrConfAsAdminOrg1 add observation of a device as Admin of organization 1, with a wrong confidence
func TestAddObsWithErrConfAsAdminOrg1(t *testing.T) {

	initL(1, "A")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"4", "1", "2", "3"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	device4 := deviceTransientInt{
		ID:    "4",
		X:     1,
		Y:     1,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3, device4)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract) = %q, want  nil`, err)
	}

	idDevice := "4"

	err = insertObsTest(contract, idDevice, idTarget)
	if err == nil {
		t.Fatalf(`Must not insert the observation, want not nil error`)
	}

	t.Logf("insert incorrect confidence: %q", err)

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}

// TestAddEvAsAdminOrg1 test an update of evidence, trust and reputation of a device as Admin of organization 1
func TestAddEvAsAdminOrg1(t *testing.T) {
	initL(1, "A")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"3", "1", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract) = %q, want  nil`, err)
	}

	idDevice := "1"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "2"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice,idTarget) = %q, want nil`, err)
	}

	idDevice = "3"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget)= %q, want nil`, err)
	}

	var devices []string

	devices = append(devices, "1", "2", "3")

	err = upDateDevicesTest(contract, devices)
	if err != nil {
		t.Fatalf(`upDateDevicesTest(contract, devices) = %q, want nil`, err)
	}

	var deviceR Device

	err, deviceR = seeDataDevice(contract, idDevice)
	if err != nil {
		t.Fatalf(`seeDataDevice(contract, idDevice) = %q, want nil`, err)
	}

	if deviceR.Ev != 1 || deviceR.Rep != 7 {
		t.Fatalf(`it calculate the incorrect evidence, want %d: %q`, 1, err)
	}

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}

// TestPosCalculateAsAdminOrg1 test the calculation of position of the target as Admin of organization 1
func TestPosCalculateAsAdminOrg1(t *testing.T) {
	initL(1, "A")

	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"3", "1", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract, devicesIn, idTarget) = %q, want  nil`, err)
	}

	idDevice := "1"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "2"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "3"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	var devices []string

	devices = append(devices, "1", "2", "3")

	err = upDateDevicesTest(contract, devices)
	if err != nil {
		t.Fatalf(`upDateDevicesTest(contract, devices) = %q, want nil`, err)
	}

	err = upDateTarget(contract, idTarget, devices)
	if err != nil {
		t.Fatalf(`upDateTarget(contract,idTarget,devices) = %q, want nil`, err)
	}

	var targetR TargetInfo

	err, targetR = seeDataTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`seeDataTarget(contract, idDevice) = %q, want nil`, err)
	}

	if targetR.Y == 0 || targetR.X == 0 {
		t.Fatalf(`it calculate the incorrect potiton, want not 0 position: %q`, err)
	}

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}

// TestSeeTargetAsUserOrg1 test access to the target as User of organization 1
func TestSeeTargetAsUserOrg1(t *testing.T) {
	initL(1, "A")

	idTarget := "7"
	devicesE := []string{"3", "1", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract, devicesIn, idTarget) = %q, want  nil`, err)
	}

	idDevice := "1"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "2"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "3"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	var devices []string

	devices = append(devices, "1", "2", "3")

	err = upDateDevicesTest(contract, devices)
	if err != nil {
		t.Fatalf(`upDateDevicesTest(contract, devices) = %q, want nil`, err)
	}

	err = upDateTarget(contract, idTarget, devices)
	if err != nil {
		t.Fatalf(`upDateTarget(contract,idTarget,devices) = %q, want nil`, err)
	}

	clientConnection.Close()
	gw.Close()

	initL(1, "U")

	var infT TargetInfo

	err, infT = seeDataTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`seeDataDevice(contract, idDevice) = %q, want nil`, err)
	}

	if infT.ID != idTarget {
		t.Fatalf("no correct id target inserted")
	}

	clientConnection.Close()
	gw.Close()

	initL(1, "A")

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

	clientConnection.Close()
	gw.Close()

	t.Logf("User can read a target of same organization")

}

// TestSeeTargetNOTUpdatedAsUserOrg1 test access to the target as User of organization 1 and not updated target
func TestSeeTargetNOTUpdatedAsUserOrg1(t *testing.T) {
	initL(1, "A")

	idTarget := "7"
	devicesE := []string{"3", "1", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract, devicesIn, idTarget) = %q, want  nil`, err)
	}

	clientConnection.Close()
	gw.Close()

	initL(1, "U")

	err, _ = seeDataTarget(contract, idTarget)
	if err == nil {
		t.Fatalf(`user can't see target not updated`)
	}

	t.Logf("user can't see target not updated")

	clientConnection.Close()
	gw.Close()

	initL(1, "A")

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

	clientConnection.Close()
	gw.Close()

}

// TestSeeTargetAsUserOrg2 test access to the target as Admin of organization 2
func TestSeeTargetAsAdminOrg2(t *testing.T) {
	initL(1, "A")

	idTarget := "7"
	devicesE := []string{"3", "1", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract, devicesIn, idTarget) = %q, want  nil`, err)
	}

	clientConnection.Close()
	gw.Close()

	initL(2, "A")

	err, _ = seeDataTarget(contract, idTarget)
	if err == nil {
		t.Fatalf(`Must not be able to read the target , want not nil error`)
	}

	t.Logf("A client can’t access the target collection of an organization that it doesn’t belong: %q", err)

	clientConnection.Close()
	gw.Close()

	initL(1, "A")

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

	clientConnection.Close()
	gw.Close()

}

// TestAddEvWithIncorrectObsAsAdminOrg1 test an update of evidence, trust and reputation of a device as Admin of organization 1, with impossible distance observation
func TestAddEvWithIncorrectObsAsAdminOrg1(t *testing.T) {
	initL(1, "A")
	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"3", "5", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "5"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "5",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"5", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract) = %q, want  nil`, err)
	}

	idDevice := "2"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "3"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "5"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	var devices []string

	devices = append(devices, "5", "2", "3")

	err = upDateDevicesTest(contract, devices)
	if err != nil {
		t.Fatalf(`upDateDevicesTest(contract, devices) = %q, want nil`, err)
	}

	var deviceR Device

	err, deviceR = seeDataDevice(contract, idDevice)
	if err != nil {
		t.Fatalf(`seeDataDevice(contract, idDevice) = %q, want nil`, err)
	}

	if deviceR.Ev != -1 && deviceR.Rep != 2 {
		t.Fatalf(`it calculate the incorrect evidence and reputation, want %d,%d: %q`, -1, 2, err)
	}

	t.Logf("Device %s has evidence -1 and reputation 2 becouse gives observation different from both other diveces", idDevice)

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}

// TestPosCalculateWithIncorrectObsAsAdminOrg1 test the calculation of position of the target as Admin of organization 1, with impossible distance observation
func TestPosCalculateWithIncorrectObsAsAdminOrg1(t *testing.T) {
	initL(1, "A")

	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"3", "5", "2"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "5"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "5",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"5", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract, devicesIn, idTarget) = %q, want  nil`, err)
	}

	idDevice := "5"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "2"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "3"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	var devices []string

	devices = append(devices, "5", "2", "3")

	err = upDateDevicesTest(contract, devices)
	if err != nil {
		t.Fatalf(`upDateDevicesTest(contract, devices) = %q, want nil`, err)
	}

	err = upDateTarget(contract, idTarget, devices)
	if err == nil {
		t.Fatalf(`You cannot obtain a potiton of the target with that observations received as input, want not nil error`)
	}

	t.Logf("You cannot obtain a potiton of the target with that observations received as input: %q", err)

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}

// TestPosCalculateWith4devicesAsAdminOrg1 test the calculation of position of the target as Admin of organization 1, with 4 devices where one gives impossible distance
func TestPosCalculateWith4devicesAsAdminOrg1(t *testing.T) {
	initL(1, "A")

	defer clientConnection.Close()
	defer gw.Close()

	idTarget := "7"
	devicesE := []string{"3", "1", "2", "5"}

	Deletedevices(contract, devicesE)
	DeleteTarget(contract, idTarget)

	var devicesIn []deviceTransientInt

	device1 := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	device2 := deviceTransientInt{
		ID:    "1",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	device3 := deviceTransientInt{
		ID:    "2",
		X:     10,
		Y:     4,
		Key:   "P",
		Neigh: []string{"1", "2"},
		Rep:   5,
	}

	device4 := deviceTransientInt{
		ID:    "5",
		X:     3,
		Y:     2,
		Key:   "P",
		Neigh: []string{"3", "2"},
		Rep:   5,
	}

	devicesIn = append(devicesIn, device1, device2, device3, device4)

	err := initLedger(contract, devicesIn, idTarget)
	if err != nil {
		t.Fatalf(`initLedger(contract, devicesIn, idTarget) = %q, want  nil`, err)
	}

	idDevice := "1"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "2"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "3"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	idDevice = "5"

	err = insertObsTest(contract, idDevice, idTarget)
	if err != nil {
		t.Fatalf(`insertObsTest(contract, idDevice, idTarget) = %q, want nil`, err)
	}

	var devices []string

	devices = append(devices, "5", "2", "3", "1")

	err = upDateDevicesTest(contract, devices)
	if err != nil {
		t.Fatalf(`upDateDevicesTest(contract, devices) = %q, want nil`, err)
	}

	err = upDateTarget(contract, idTarget, devices)
	if err != nil {
		t.Fatalf(`upDateTarget(contract,idTarget,devices) = %q, want nil`, err)
	}

	var targetR TargetInfo

	err, targetR = seeDataTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`seeDataDevice(contract, idDevice) = %q, want nil`, err)
	}

	if targetR.Y == 0 || targetR.X == 0 {
		t.Fatalf(`it calculate the incorrect potiton, want not 0 position: %q`, err)
	}

	t.Logf("Even if there is a device with incorrect observation, the network select devices with higher trust ")

	err = Deletedevices(contract, devicesE)
	if err != nil {
		t.Fatalf(`Deletedevices(contract,devices) = %q, want nil`, err)
	}

	err = DeleteTarget(contract, idTarget)
	if err != nil {
		t.Fatalf(`DeleteTarget(contract, idTarget) = %q, want nil`, err)
	}

}
