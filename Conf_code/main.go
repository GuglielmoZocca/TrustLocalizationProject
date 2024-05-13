/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//Code to utilize by admin to operate on ledger

package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// deviceTransientInt type for the inizialitation of the device
type deviceTransientInt struct {
	ID    string   `json:"deviceID"`
	X     float32  `json:"coordinateX"`
	Y     float32  `json:"coordinateY"`
	Key   string   `json:"Key"`
	Neigh []string `json:"Neighborhood"`
	Rep   int      `json:"Reputation"`
}

// TargetInfo detail of a certain target
type TargetInfo struct {
	ID   string  `json:"TargetID"`    //ID of the target
	X    float32 `json:"coordinateX"` //X coordinate of the target
	Y    float32 `json:"coordinateY"` //Y coordinate of the target
	Date string  `json:"Date"`        //Date of the update
	Upd  bool    `json:"Update"`      //Indicate if target is updated
}

// deviceTransientUp type for the update of the device
type deviceTransientUp struct {
	ID    string   `json:"deviceID"`
	Key   string   `json:"Key"`          //New key
	Neigh []string `json:"Neighborhood"` //New neighborhood
}

// Device data inserted in the ledger
type Device struct {
	ID    string   `json:"deviceID"`     //ID of the device
	Key   string   `json:"Key""`         //Key of the decription of data of the device
	X     float32  `json:"coordinateX"`  //X coordinate of the device
	Y     float32  `json:"coordinateY"`  //Y coordinate of the device
	Obs   int      `json:"observation"`  //Obs distance observed from the device
	Conf  float32  `json:"Confidence"`   //Conf confidence of the device
	Ev    float32  `json:"Evidence"`     //Ev evidence of the device
	Rep   int      `json:"Reputation"`   //Rep reputation of the device
	Trust float32  `json:"Trust"`        //Trust of the device
	Neigh []string `json:"Neighborhood"` //Neighborhood id of the devices near of the device
}

// Parameter of the application
const (
	collectionD = "DeviceAdmin1PrivateCollection" //Collection where are memorized the devices
	collectionT = "TargetOrg1PrivateCollection"   //Collection where are memorized the targets
)

// Info of Admin in organization 1
const (
	mspIDA1        = "Org1MSP"
	cryptoPathA1   = "../Network/organizations/peerOrganizations/org1.example.com"
	certPathA1     = cryptoPathA1 + "/users/Admin@org1.example.com/msp/signcerts"
	keyPathA1      = cryptoPathA1 + "/users/Admin@org1.example.com/msp/keystore"
	tlsCertPathA1  = cryptoPathA1 + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpointA1 = "localhost:7051"
	gatewayPeerA1  = "peer0.org1.example.com"
)

func main() {

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPathA1, gatewayPeerA1, peerEndpointA1)
	defer clientConnection.Close()

	//Creation of identity and Tls sign for the client
	id := newIdentity(certPathA1, mspIDA1)
	sign := newSign(keyPathA1)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
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
	defer gw.Close()

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

	contract := network.GetContract(chaincodeName)

	var risp int //answer from admin

	var idD string //id of device to consider

	var idTarget string //id of the target to read

	var x float32 //coordinate x of device

	var y float32 //coordinate y of device

	var pubKe string //key of device

	var rep int //Reputation of the device

	var neigh string //neighbor of device

	var neighs []string //neighborhood of device

	var goon bool //indicate if continue the opearation or not

	var rispN string //answer from admin in the cycle

	var deviceinput deviceTransientInt //new device to insert

	var deviceUpd deviceTransientUp //new data to do the update of device

	var deviceData Device //data device read

	var targetData TargetInfo //target data read

	//Ask admin to what operation to do
	for true {

		fmt.Println("Do you want create(1), read(2), update(3) a device or read target(4):")

		fmt.Scanln(&risp)

		//Create a new device
		if risp == 1 {

			fmt.Println("What is the id of device")

			fmt.Scanln(&idD)

			fmt.Println("What is the x coordinate of the device:")

			_, err = fmt.Scanln(&x)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			fmt.Println("What is the y coordinate of the device:")

			_, err = fmt.Scanln(&y)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			fmt.Println("What is the key of the device:")

			_, err = fmt.Scanln(&pubKe)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			neighs = nil

			fmt.Println("Add a neigh:")

			_, err = fmt.Scanln(&neigh)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			neighs = append(neighs, neigh)

			fmt.Println("Add the second neigh:")

			_, err = fmt.Scanln(&neigh)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			neighs = append(neighs, neigh)

			goon = true
			for goon {

				fmt.Println("Want add a new neigh: yes or no")

				fmt.Scanln(&rispN)
				if rispN == "yes" {

					fmt.Println("Add a neigh:")

					_, err = fmt.Scanln(&neigh)
					if err != nil {
						fmt.Printf("error read input: %q\n", err)
						continue
					}

					neighs = append(neighs, neigh)

				}
				if rispN == "no" {

					goon = false

				}

			}

			fmt.Println("What is the reputation of the device:")

			_, err = fmt.Scanln(&rep)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			deviceinput = deviceTransientInt{
				ID:    idD,
				X:     x,
				Y:     y,
				Neigh: neighs,
				Rep:   rep,
			}

			err = insertDevice(contract, deviceinput)
			if err != nil {

				fmt.Printf("error insert new device: %q\n", err)

			}
		}

		//Read data of a device
		if risp == 2 {

			fmt.Println("What is the id of device")

			fmt.Scanln(&idD)

			err, deviceData = seeDataDevice(contract, idD)
			if err != nil {

				fmt.Printf("error in read device:%q\n", err)
				continue

			}

			fmt.Printf("ID: %s, Key: %s, X: %f, Y: %f, Reputation: %d, Trust: %f, Confidence: %f, Evidence: %f, Observation: %d, Neigh: %s\n", deviceData.ID, deviceData.Key, deviceData.X, deviceData.Y, deviceData.Rep, deviceData.Trust, deviceData.Conf, deviceData.Ev, deviceData.Obs, deviceData.Neigh)

		}

		//Update a device
		if risp == 3 {

			fmt.Println("What is the id of device")

			fmt.Scanln(&idD)

			fmt.Println("What is the key of the device:")

			_, err = fmt.Scanln(&pubKe)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			neighs = nil

			fmt.Println("Add a neigh:")

			_, err = fmt.Scanln(&neigh)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			neighs = append(neighs, neigh)

			fmt.Println("Add the second neigh:")

			_, err = fmt.Scanln(&neigh)
			if err != nil {
				fmt.Printf("error read input: %q\n", err)
				continue
			}

			neighs = append(neighs, neigh)

			goon = true
			for goon {

				fmt.Println("Want add a new neigh: yes or no")

				fmt.Scanln(&rispN)
				if rispN == "yes" {

					fmt.Println("Add a neigh:")

					_, err = fmt.Scanln(&neigh)
					if err != nil {
						fmt.Printf("error read input: %q\n", err)
						continue
					}

					neighs = append(neighs, neigh)

				}
				if rispN == "no" {

					goon = false

				}

			}

			deviceUpd = deviceTransientUp{
				ID:    idD,
				Key:   pubKe,
				Neigh: neighs,
			}

			err = updateDevice(contract, deviceUpd)
			if err != nil {

				fmt.Printf("error update device: %q\n", err)

			}

		}

		//Read data of a target
		if risp == 4 {
			fmt.Println("What is the id of target")

			fmt.Scanln(&idTarget)

			err, targetData = seeDataTarget(contract, idTarget)

			if err != nil {
				fmt.Printf("error read target: %q\n", err)

			}

			fmt.Printf("ID:%s, updated:%v, X: %f, Y: %f, Date: %s\n", targetData.ID, targetData.Upd, targetData.X, targetData.Y, targetData.Date)

		}

	}

}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection(tlsCertPath string, gatewayPeer string, peerEndpoint string) *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(certPath string, mspID string) *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign(keyPath string) identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// insertDevice create a device
func insertDevice(contract *client.Contract, device deviceTransientInt) error {

	deviceInput, err := json.Marshal(device)

	inv := map[string][]byte{"Device_properties": deviceInput}

	_, err = contract.Submit("CreateDevice", client.WithTransient(inv), client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to create device %s: %w", device.ID, err)
	}

	return nil
}

// DeleteDevice delete device idDevice
func DeleteDevice(contract *client.Contract, idDevice string) error {

	type deviceData struct {
		ID string `json:"deviceID"`
	}

	deviceReq := deviceData{

		ID: idDevice,
	}

	deviceReqJason, err := json.Marshal(deviceReq)

	if err != nil {
		return fmt.Errorf("failed marshalling asset")
	}

	inv := map[string][]byte{"device_delete": deviceReqJason}

	_, err = contract.Submit("DeleteDevice", client.WithTransient(inv), client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to delete device %s: %w", idDevice, err)
	}

	return nil
}

// seeDataDevice Retrive data device from id
func seeDataDevice(contract *client.Contract, idDevice string) (error, Device) {

	var deviceDataRcv Device

	type deviceData struct {
		ID string `json:"deviceID"`
	}

	deviceReq := deviceData{

		ID: idDevice,
	}

	deviceReqJason, err := json.Marshal(deviceReq)

	if err != nil {
		return fmt.Errorf("failed marshalling device"), deviceDataRcv
	}

	inv := map[string][]byte{"device_data": deviceReqJason}

	evaluateResult, err := contract.Evaluate("ReadDevice", client.WithTransient(inv), client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to read device %s: %w", idDevice, err), deviceDataRcv
	}

	err = json.Unmarshal(evaluateResult, &deviceDataRcv)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err), deviceDataRcv
	}

	return nil, deviceDataRcv
}

// updateDevice update device id or neighborhood or key
func updateDevice(contract *client.Contract, device deviceTransientUp) error {

	deviceInput, err := json.Marshal(device)

	inv := map[string][]byte{"Device_properties": deviceInput}

	_, err = contract.Submit("UpdateDevice", client.WithTransient(inv), client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to update device %s: %w", device.ID, err)
	}

	return nil
}

// seeDataTarget Retrive data target from id
func seeDataTarget(contract *client.Contract, idTarget string) (error, TargetInfo) {

	var targetDataRcv TargetInfo

	type targetData struct {
		ID string `json:"TargetID"`
	}

	targetReq := targetData{

		ID: idTarget,
	}

	targetReqJason, err := json.Marshal(targetReq)

	if err != nil {
		return fmt.Errorf("failed marshalling target"), targetDataRcv
	}

	inv := map[string][]byte{"target_data": targetReqJason}

	evaluateResult, err := contract.Evaluate("ReadTarget", client.WithTransient(inv), client.WithArguments(collectionT))
	if err != nil {
		return fmt.Errorf("failed to read target %s: %w", idTarget, err), targetDataRcv
	}

	err = json.Unmarshal(evaluateResult, &targetDataRcv)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err), targetDataRcv
	}

	return nil, targetDataRcv
}
