// TEST CONFIGURATION CODE
package main

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)

//Test in the context to be in org 1

var contract *client.Contract
var gw *client.Gateway
var clientConnection *grpc.ClientConn

// initL initialize the connexion with the ledger
func initL() {
	var err error
	// The gRPC client connection should be shared by all Gateway connections to this endpoint

	clientConnection = newGrpcConnection(tlsCertPathA1, gatewayPeerA1, peerEndpointA1)
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

// TestUpdateDeviceAsAdminOrg1 test update a device as Admin of organization 1
func TestUpdateDeviceAsAdminOrg1(t *testing.T) {

	initL()
	defer clientConnection.Close()
	defer gw.Close()

	device := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	DeleteDevice(contract, device.ID)

	err := insertDevice(contract, device)
	if err != nil {

		log.Fatalf("Error create device:%q", err)

	}

	deviceUpd := deviceTransientUp{
		ID:    device.ID,
		Key:   "V",
		Neigh: []string{"4", "2"},
	}

	err = updateDevice(contract, deviceUpd)
	if err != nil {

		log.Fatalf("error update device: %q", err)

	}

	err, deviceData := seeDataDevice(contract, deviceUpd.ID)
	if err != nil {

		log.Fatalf("seeDataDevice(contract, deviceUpd.ID) = %q, want nil", err)

	}

	if deviceData.Key != deviceUpd.Key || device.X != deviceData.X {

		log.Fatalf("non correct data in updated device: %q", err)

	}

	err = DeleteDevice(contract, deviceUpd.ID)
	if err != nil {
		log.Fatalf("error delete device: %q", err)

	}

}

// TestSeeTargetAsUserOrg1 test access to the device as Admin of organization 1
func TestSeeDevicetAsAdminOrg1(t *testing.T) {

	initL()
	defer clientConnection.Close()
	defer gw.Close()

	device := deviceTransientInt{
		ID:    "3",
		X:     5,
		Y:     8,
		Key:   "P",
		Neigh: []string{"2", "1"},
		Rep:   5,
	}

	DeleteDevice(contract, device.ID)

	err := insertDevice(contract, device)
	if err != nil {

		log.Fatalf("Error create device:%q", err)

	}

	var deviceInf Device

	err, deviceInf = seeDataDevice(contract, device.ID)
	if err != nil {
		t.Fatalf(`seeDataDevice(contract,device.ID) want not nil answer`)
	}

	if deviceInf.ID != device.ID {
		t.Fatalf(`read wrong devicer`)
	}

	err = DeleteDevice(contract, device.ID)
	if err != nil {
		log.Fatalf("error delete device: %q", err)

	}
}
