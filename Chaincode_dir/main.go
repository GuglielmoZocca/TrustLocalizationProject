/*
SPDX-License-Identifier: Apache-2.0
*/
//File that creates the chaincode
package main

import (
	"log"

	"TrustLocalizationProject/Chaincode_dir/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	positionChaincode, err := contractapi.NewChaincode(&chaincode.PositionContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-private-data chaincode: %v", err)
	}

	if err := positionChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-private-data chaincode: %v", err)
	}
}
