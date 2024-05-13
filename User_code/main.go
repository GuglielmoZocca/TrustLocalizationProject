/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//Code to utilize by the user

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

// TargetInfo detail of a certain target
type TargetInfo struct {
	ID   string  `json:"TargetID"`    //ID of the target
	X    float32 `json:"coordinateX"` //X coordinate of the target
	Y    float32 `json:"coordinateY"` //Y coordinate of the target
	Date string  `json:"Date"`        //Date of the update
	Upd  bool    `json:"Update"`      //Indicate if target is updated
}

// Parameter of the application
const (
	collectionT = "TargetOrg1PrivateCollection" //Collection where are memorized the targets
)

// Info of User in organization 1
const (
	mspIDU1        = "Org1MSP"
	cryptoPathU1   = "../Network/organizations/peerOrganizations/org1.example.com"
	certPathU1     = cryptoPathU1 + "/users/User1@org1.example.com/msp/signcerts"
	keyPathU1      = cryptoPathU1 + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPathU1  = cryptoPathU1 + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpointU1 = "localhost:7051"
	gatewayPeerU1  = "peer0.org1.example.com"
)

// Info of User in organization 2
const (
	mspIDU2        = "Org2MSP"
	cryptoPathU2   = "../Network/organizations/peerOrganizations/org2.example.com"
	certPathU2     = cryptoPathU2 + "/users/User1@org2.example.com/msp/signcerts"
	keyPathU2      = cryptoPathU2 + "/users/User1@org2.example.com/msp/keystore"
	tlsCertPathU2  = cryptoPathU2 + "/peers/peer0.org2.example.com/tls/ca.crt"
	peerEndpointU2 = "localhost:7052"
	gatewayPeerU2  = "peer0.org2.example.com"
)

func main() {

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPathU1, gatewayPeerU1, peerEndpointU1)
	defer clientConnection.Close()

	//Creation of identity and Tls sign for the client
	id := newIdentity(certPathU1, mspIDU1)
	sign := newSign(keyPathU1)

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

	var idTarget string //id target user want to read

	var risp string //answer from the user

	var Target TargetInfo

	for true {

		fmt.Println("Do you want read the target info: yes or no")

		fmt.Scanln(&risp)

		if risp == "yes" {

			fmt.Println("Give the id of target")

			fmt.Scanln(&idTarget)

			err, Target = seeDataTarget(contract, "7")
			if err != nil {
				fmt.Printf("Problem with reading the target: %q\n", idTarget)
				continue
			}
			fmt.Printf("ID:%s, updated:%v, X: %f, Y: %f, Date: %s\n", Target.ID, Target.Upd, Target.X, Target.Y, Target.Date)

		}

		if risp == "no" {

			return

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

	//If target is not updated gives an error
	if !targetDataRcv.Upd {

		return fmt.Errorf("target isn't already updated"), targetDataRcv
	}

	return nil, targetDataRcv
}
