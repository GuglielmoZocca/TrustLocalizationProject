/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//Code to utilize in the gateway that communicate with the devices and do all operations of calculation of target position

package main

import (
	"bufio"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
	Minconf     = 0.4                             //minimum confidence fo a device
	Maxconf     = 1                               //max confidence of a device
	Prh         = 2                               //Reward or penalty in case of high evidence and high confidence, and in case of low evidence and high confidence respectively
	Prl         = 1                               //Reward or penalty in case of high evidence and low confidence, and in case of low evidence and low confidence respectively
	threashconf = 0.7                             //threshold that indicate when a confidence value is high or low
	threashev   = 0                               //threshold that indicate when an evidence value is high or low
	collectionD = "DeviceAdmin1PrivateCollection" //Collection where are memorized the devices
	collectionT = "TargetOrg1PrivateCollection"   //Collection where are memorized the targets
	treshErr    = 0.01                            //threshold that indicate when an evidence value is high or low
	mrep        = 20                              //max reputation
	numbRead    = 6                               //Dimension of a batch of observations to read before compute an average of the observation
)

// Info of User in organization 1
const (
	mspIDU1        = "Org1MSP"                                                     //ID of organization msp
	cryptoPathU1   = "../Network/organizations/peerOrganizations/org1.example.com" //Path where find all cryptographic material of organization
	certPathU1     = cryptoPathU1 + "/users/User1@org1.example.com/msp/signcerts"  //Path where to find client certificate
	keyPathU1      = cryptoPathU1 + "/users/User1@org1.example.com/msp/keystore"   //Path where to find client private key
	tlsCertPathU1  = cryptoPathU1 + "/peers/peer0.org1.example.com/tls/ca.crt"     //Path where to find peer tls certificate
	peerEndpointU1 = "localhost:7051"                                              //Address of endpoint peer
	gatewayPeerU1  = "peer0.org1.example.com"                                      //Name of endpoint peer
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

// Info of Admin in organization 2
const (
	mspIDA2        = "Org2MSP"
	cryptoPathA2   = "../Network/organizations/peerOrganizations/org2.example.com"
	certPathA2     = cryptoPathA2 + "/users/Admin@org2.example.com/msp/signcerts"
	keyPathA2      = cryptoPathA2 + "/users/Admin@org2.example.com/msp/keystore"
	tlsCertPathA2  = cryptoPathA2 + "/peers/peer0.org2.example.com/tls/ca.crt"
	peerEndpointA2 = "localhost:7052"
	gatewayPeerA2  = "peer0.org2.example.com"
)

// deviceTransientInt type for the initialization of the device
type deviceTransientInt struct {
	ID    string   `json:"deviceID"`
	X     float32  `json:"coordinateX"`
	Y     float32  `json:"coordinateY"`
	Key   string   `json:"Key"`
	Neigh []string `json:"Neighborhood"`
	Rep   int      `json:"Reputation"`
}

var numOb int //Number of observation inserted

var numUp int //Number of successful update of evidence, reputation and trust

func main() {

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPathU1, gatewayPeerA1, peerEndpointA1)
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

	//init the ledger

	var devicesIn []deviceTransientInt //Devices the initialize the network
	var devicesID []string             //Id of devices in the network currently

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

	idTarget := "7" //ID of the target

	//Initialize the network
	err = initLedger(contract, devicesIn, idTarget)
	if err != nil {
		fmt.Printf("%q", err)
		return
	}

	readMap := make(map[string]chan string) //Map of the channels associated to every device to receive data from them

	//open file of time for time operation of gateway
	f, err := os.OpenFile("TimeAppGateway.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		f.Close()
		return
	}

	defer f.Close()

	fmt.Println("start")

	//Execute loop to update the position of the target
	for true {
		t1 := time.Now()

		//Access to all id of devices in the network
		err, devicesID = seeAllIDDevice(contract)
		if err != nil {
			fmt.Printf("%q", err)
			continue
		}

		numOb = 0

		for _, device := range devicesID {
			//Create a new channel only if a device don't have already one
			if readMap[device] == nil {
				c := make(chan string, 10)
				readMap[device] = c
				//Start process to read from file data of the device
				go ReadStreamDevice(device, readMap[device])

			}

			//insert obs of data read
			go insertObs(contract, device, idTarget, readMap[device])

		}
		//wait until all obs are inserted correctly
		for numOb != len(devicesID) {
		}
		numOb = 0

		fmt.Println("insert obs")

		//Update evidence, reputation and trust of the device
		upDateDevices(contract, devicesID)

		fmt.Println("insert ev")

		//Update position of target
		err = upDateTarget(contract, idTarget, devicesID)
		if err != nil {
			fmt.Printf("error update position of the target: %q", err)
			continue
		}

		t2 := time.Now()

		diff := t2.Sub(t1)
		fmt.Println(diff)

		_, err = f.WriteString(diff.String() + "\n")
		if err != nil {
			return
		}

		fmt.Println("insert pos")

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

// initLedger Inizialize the ledger with initial devices and target
func initLedger(contract *client.Contract, devices []deviceTransientInt, idTarget string) error {

	//Initializtion devices
	var inv map[string][]byte
	for _, device := range devices {

		deviceInput, err := json.Marshal(device)

		inv = map[string][]byte{"Device_properties": deviceInput}

		_, err = contract.Submit("CreateDevice", client.WithTransient(inv), client.WithArguments(collectionD))
		if err != nil {
			return fmt.Errorf("failed to create device %s: %w", device.ID, err)
		}

	}

	//Inizialization target

	type targetTransientInput struct {
		ID string `json:"TargetID"`
	}

	targetData := targetTransientInput{
		ID: idTarget,
	}

	targetInput, err := json.Marshal(targetData)

	inv = map[string][]byte{"Target_properties": targetInput}

	_, err = contract.Submit("CreateTarget", client.WithTransient(inv), client.WithArguments(collectionT))
	if err != nil {
		return fmt.Errorf("failed to create target %s: %w", idTarget, err)
	}

	return nil

}

// EncryptDecrypt runs a XOR encryption on the input string, encrypting it if it hasn't already been,
// and decrypting it if it has, using the key provided. Decrypt data of the device
func EncryptDecrypt(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}

	return output
}

// ReadStreamDevice Read device file
func ReadStreamDevice(idDevice string, c chan string) error {
	noErr := true
	//Continue to try to read file until it doesn't read all file
	for noErr {

		//open file of encrypted data of the devices
		f, err := os.Open("Devices_cript/Device[" + idDevice + "]_ript.txt")

		if err != nil {
			f.Close()
			continue
		}

		// remember to close the file at the end of the program
		defer f.Close()

		//parse data and sent theme in the channel
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			//simulate time network and trasmission delay
			time.Sleep(100 * time.Millisecond)
			c <- scanner.Text()

		}

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

		close(c)

		noErr = false

	}

	return nil

}

// getHash calculate the hash of data device
func getHash(s string) uint64 {
	hash := uint64(0)
	for i := 0; i < len(s); i++ {
		hash = (hash * uint64(10)) + uint64(s[i]) - uint64('0')
	}

	return hash
}

// insertObs Insert the observation of a device from a file
func insertObs(contract *client.Contract, idDevice string, idTarget string, c chan string) error {
	var noErr bool //var to indicate if there is error
	noErr = true
	numRead := 0 //Number of read done in a batch
	//Continue to read from the channel until it has a batch of observation or it hasn't anything more to read
	for noErr {

		//Ask device data for the key of decryption
		type deviceData struct {
			ID string `json:"deviceID"`
		}

		deviceReq := deviceData{

			ID: idDevice,
		}

		deviceReqJason, err := json.Marshal(deviceReq)

		if err != nil {
			continue
		}

		inv := map[string][]byte{"device_data": deviceReqJason}

		evaluateResult, err := contract.Evaluate("ReadDevice", client.WithTransient(inv), client.WithArguments(collectionD))
		if err != nil {
			continue
		}

		var deviceDataRcv Device

		err = json.Unmarshal(evaluateResult, &deviceDataRcv)
		if err != nil {
			continue
		}

		key := deviceDataRcv.Key //Key for the decryption

		i := 0
		line1 := ""
		line2 := ""
		var split []string
		var hash string

		num := 0             //number of data
		Tot_distance := 0    //Total distance
		var Tot_conf float32 //total confidence
		Tot_conf = 0.0

		var tmpnum int
		var tmpconf float64

		//Scroll the channel
		for line := range c {

			if i%2 == 0 {
				line1 = line
				line1 = EncryptDecrypt(line1, key)
			} else {
				line2 = line
				hash = strconv.FormatUint(getHash(line1+"\n"), 10)
				//check if the data are not tampered
				if strings.Compare(line2, hash) == 0 {

					split = strings.Split(line1, ",")
					//Check if the data received are from the device requested
					if split[0] != idDevice {
						err = fmt.Errorf("rsv not data of correct device")
						break

					}
					//Check if the data received are toward the correct target
					if split[1] != idTarget {
						err = fmt.Errorf("rsv not data of correct target")
						break

					}
					num++
					//read distance
					tmpnum, err = strconv.Atoi(split[2])
					if err != nil {
						break
					}
					//read confidence
					tmpconf, err = strconv.ParseFloat(split[3], 32)
					if err != nil {
						break
					}

					Tot_distance = Tot_distance + tmpnum
					Tot_conf = Tot_conf + float32(tmpconf)

				} else {
					err = fmt.Errorf("no hash correspondance")
					break
				}

			}
			numRead++
			//finish to scroll the channel when arrive to the dimension of the batch
			if numRead == numbRead {
				numRead = 0
				break
			}
			i++
		}

		if err != nil {
			continue

		}

		if num == 0 {
			continue
		}

		/*if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}*/

		avg_distance := Tot_distance / num  //average distance
		avg_conf := Tot_conf / float32(num) //average confidence

		type deviceTransientInput struct {
			ID      string  `json:"deviceID"`
			Obs     int     `json:"observation"`
			Conf    float32 `json:"Confidence"`
			MinConf float32 `json:"MinConf"` //Minimum confidence of the device
			MaxConf float32 `json:"MaxConf"` //Maximum confidence of the device
		}

		newObConf := deviceTransientInput{
			ID:      idDevice,
			Obs:     avg_distance,
			Conf:    avg_conf,
			MinConf: float32(Minconf),
			MaxConf: float32(Maxconf),
		}

		deviceInput, err := json.Marshal(newObConf)

		inv = map[string][]byte{"Device_properties": deviceInput}

		if err != nil {
			continue
		}
		_, err = contract.Submit("UpdateDeviceObsConf", client.WithTransient(inv), client.WithArguments(collectionD))
		if err != nil {
			continue
		}

		noErr = false

	}

	numOb++

	return nil

}

// insertObs Insert the observation of a device from a file for testing
func insertObsTest(contract *client.Contract, idDevice string, idTarget string) error {

	//Read device for the key of decryption
	type deviceData struct {
		ID string `json:"deviceID"`
	}

	deviceReq := deviceData{

		ID: idDevice,
	}

	deviceReqJason, err := json.Marshal(deviceReq)

	if err != nil {
		return fmt.Errorf("failed marshalling device")
	}

	inv := map[string][]byte{"device_data": deviceReqJason}

	evaluateResult, err := contract.Evaluate("ReadDevice", client.WithTransient(inv), client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	var deviceDataRcv Device

	err = json.Unmarshal(evaluateResult, &deviceDataRcv)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	key := deviceDataRcv.Key //Key for the decryption

	// open file of encrypted data of the devices
	f, err := os.Open("Devices_cript/Device[" + idDevice + "]_ript.txt")

	if err != nil {

		return err
	}

	// remember to close the file at the end of the program
	defer f.Close()

	scanner := bufio.NewScanner(f)

	i := 0
	line1 := ""
	line2 := ""
	var split []string
	var hash string

	num := 0             //number of data
	Tot_distance := 0    //Total distance
	var Tot_conf float32 //total confidence
	Tot_conf = 0.0

	var tmpnum int
	var tmpconf float64

	//Scroll the file
	for scanner.Scan() {

		if i%2 == 0 {
			line1 = scanner.Text()
			line1 = EncryptDecrypt(line1, key)
		} else {
			line2 = scanner.Text()
			hash = strconv.FormatUint(getHash(line1+"\n"), 10)
			//check if the data are not tampered
			if strings.Compare(line2, hash) == 0 {

				//Check if the data received are from the device requested
				split = strings.Split(line1, ",")
				if split[0] != idDevice {
					return fmt.Errorf("rsv not data of correct device")

				}
				//Check if the data received are toward the correct target
				if split[1] != idTarget {
					return fmt.Errorf("rsv not data of correct target")

				}
				num++

				//Read the device distance
				tmpnum, err = strconv.Atoi(split[2])
				if err != nil {
					return fmt.Errorf("conversion incorrect of string")
				}
				//Read the device confidence
				tmpconf, err = strconv.ParseFloat(split[3], 32)
				if err != nil {
					return fmt.Errorf("conversion incorrect of string")
				}

				Tot_distance = Tot_distance + tmpnum
				Tot_conf = Tot_conf + float32(tmpconf)

			} else {
				return fmt.Errorf("no hash correspondance")
			}

		}
		i++
	}

	if num == 0 {
		return fmt.Errorf("no observation found")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	avg_distance := Tot_distance / num  //average distance
	avg_conf := Tot_conf / float32(num) //average confidence

	type deviceTransientInput struct {
		ID      string  `json:"deviceID"`
		Obs     int     `json:"observation"`
		Conf    float32 `json:"Confidence"`
		MinConf float32 `json:"MinConf"` //Minimum confidence of the device
		MaxConf float32 `json:"MaxConf"` //Maximum confidence of the device
	}

	newObConf := deviceTransientInput{
		ID:      idDevice,
		Obs:     avg_distance,
		Conf:    avg_conf,
		MinConf: float32(Minconf),
		MaxConf: float32(Maxconf),
	}

	deviceInput, err := json.Marshal(newObConf)

	inv = map[string][]byte{"Device_properties": deviceInput}

	if err != nil {
		return fmt.Errorf("failed marshalling asset")
	}
	_, err = contract.Submit("UpdateDeviceObsConf", client.WithTransient(inv), client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil

}

// Function that try to update the evidence, reputation and trust of a device
func updateEvReTr(contract *client.Contract, device string) {
	noErr := true
	//Try to update until it has success
	for noErr {

		type deviceTransientInput struct {
			ID          string  `json:"deviceID"`
			PRH         int     `json:"PRH"`         //Reward or penality in case of high evidence and high confidence, and in case of low evidence and high confidence respectively
			PRL         int     `json:"PRL"`         //Reward or penality in case of high evidence and low confidence, and in case of low evidence and low confidence respectively
			ThreashConf float32 `json:"threashConf"` //threashold that indicate when a confidence value is high or low
			ThreashEv   float32 `json:"threashEv"`   //threashold that indicate when a evidence value is high or low
			MaxRep      int     `json:"maxRep"`      //max reputation
		}

		var deviceInput deviceTransientInput

		var inv map[string][]byte

		deviceInput.ID = device
		deviceInput.PRH = Prh
		deviceInput.PRL = Prl
		deviceInput.ThreashEv = float32(threashev)
		deviceInput.ThreashConf = float32(threashconf)
		deviceInput.MaxRep = mrep

		deviceInputJason, err := json.Marshal(deviceInput)

		if err != nil {
			continue
		}

		inv = map[string][]byte{"Device_properties": deviceInputJason}

		_, err = contract.Submit("UpdateDeviceEvRep", client.WithTransient(inv), client.WithArguments(collectionD))
		if err != nil {
			continue
		}
		noErr = false
	}

	numUp++

}

// upDateDevices update the evidence, reputation and trust of the devices
func upDateDevices(contract *client.Contract, devices []string) error {

	numUp = 0 //Number of successful update of evidence, reputation and trust

	for _, elem := range devices {
		go updateEvReTr(contract, elem)
	}

	//Wait until all devices are updated
	for numUp != len(devices) {
	}

	numUp = 0

	return nil
}

// upDateDevices update the evidence, reputation and trust of the device for testing
func upDateDevicesTest(contract *client.Contract, devices []string) error {
	type deviceTransientInput struct {
		ID          string  `json:"deviceID"`
		PRH         int     `json:"PRH"`         //Reward or penality in case of high evidence and high confidence, and in case of low evidence and high confidence respectively
		PRL         int     `json:"PRL"`         //Reward or penality in case of high evidence and low confidence, and in case of low evidence and low confidence respectively
		ThreashConf float32 `json:"threashConf"` //threashold that indicate when a confidence value is high or low
		ThreashEv   float32 `json:"threashEv"`   //threashold that indicate when a evidence value is high or low
		MaxRep      int     `json:"maxRep"`      //max reputation
	}

	var deviceInput deviceTransientInput

	var inv map[string][]byte

	for _, elem := range devices {

		deviceInput.ID = elem
		deviceInput.PRH = Prh
		deviceInput.PRL = Prl
		deviceInput.ThreashEv = float32(threashev)
		deviceInput.ThreashConf = float32(threashconf)
		deviceInput.MaxRep = mrep

		deviceInputJason, err := json.Marshal(deviceInput)

		if err != nil {
			return fmt.Errorf("failed marshalling asset")
		}

		inv = map[string][]byte{"Device_properties": deviceInputJason}

		_, err = contract.Submit("UpdateDeviceEvRep", client.WithTransient(inv), client.WithArguments(collectionD))
		if err != nil {
			return fmt.Errorf("failed to update device evidence %s: %w", elem, err)
		}

	}

	return nil
}

// upDateTarget Calculate the position of the target idTarget
func upDateTarget(contract *client.Contract, idTarget string, devices []string) error {

	type targetPosisiton struct {
		ID      string   `json:"TargetID"`
		Thresh  float64  `json:"ThreshErr"` //Threshold of difference of distance considerable as error
		Devices []string `json:"DevicesUp"` //Devices to consider for the calculation
	}

	targetInput := targetPosisiton{
		ID:      idTarget,
		Thresh:  float64(treshErr),
		Devices: devices,
	}

	targetInputJason, err := json.Marshal(targetInput)

	if err != nil {
		return fmt.Errorf("failed marshalling asset")
	}

	inv := map[string][]byte{"target_Posisiton": targetInputJason}

	_, err = contract.Submit("PositionTarget", client.WithTransient(inv), client.WithArguments(collectionD, collectionT, time.Now().Format(time.DateTime)))
	if err != nil {
		return fmt.Errorf("failed to calculate position: %w", err)
	}

	return nil
}

// seeDataTarget Retrieve data target from id given
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

// seeDataDevice Retrieve data device from id given
func seeDataDevice(contract *client.Contract, idDevice string) (error, Device) {

	var deviceDataRcv Device

	//ask device for the key of decryption
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

// seeAllIDDevice Obtain all id of devices in the collection
func seeAllIDDevice(contract *client.Contract) (error, []string) {

	evaluateResult, err := contract.Evaluate("ReadAllIDDevices", client.WithArguments(collectionD))
	if err != nil {
		return fmt.Errorf("failed to read id devices: %q", err), nil
	}

	var devices []string
	err = json.Unmarshal(evaluateResult, &devices)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err), nil
	}

	return nil, devices

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

// DeleteTarget delete target idTarget
func DeleteTarget(contract *client.Contract, idTarget string) error {

	type targetData struct {
		ID string `json:"targetID"`
	}

	targetReq := targetData{

		ID: idTarget,
	}

	targetReqJason, err := json.Marshal(targetReq)

	if err != nil {
		return fmt.Errorf("failed marshalling target")
	}

	inv := map[string][]byte{"target_delete": targetReqJason}

	_, err = contract.Submit("DeleteTarget", client.WithTransient(inv), client.WithArguments(collectionT))
	if err != nil {
		return fmt.Errorf("failed to delete target %s: %w", idTarget, err)
	}

	return nil
}

// Deletedevices delete devices indicated
func Deletedevices(contract *client.Contract, devices []string) error {
	for _, d := range devices {

		err := DeleteDevice(contract, d)
		if err != nil {
			return fmt.Errorf("failed to delete device %s: %w", d, err)
		}
	}

	return nil

}
