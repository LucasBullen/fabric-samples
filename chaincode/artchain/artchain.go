/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
//	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the art structure, with 6 properties.  Structure tags are used by encoding/json library
type Art struct {
	Id   string `json:"id"`
	Name  string `json:"name"`
	Author string `json:"author"`
	Description string `json:"description"`
	OwnerID  string `json:"owner_id"`
	BcnValue  string `json:"bcn_value"`
	IsListed  string `json:"is_listed"`
	AuctionStartTime  string `json:"start_time"`
	AuctionEndTime  string `json:"end_time"`
	HighestBidderId  string `json:"highest_bidder_id"`
	StartingBid  string `json:"starting_bid"`
}

/*
 * The Init method is called when the Smart Contract "artBlock" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "artBlock"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryArtById" {
		return s.queryArtById(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "list" {
		return s.list(APIstub, args)
	} else if function == "queryAllArt" {
		return s.queryAllArt(APIstub)
	} else if function == "setStatus" {
		return s.setStatus(APIstub, args)
	} else if function == "setPrice" {
		return s.setPrice(APIstub, args)
	} else if function == "setUpAuction" { 	//artid, auctionstartTime auctionEndTime, startingbid
		return s.setUpAuction(APIstub, args)
	} else if function == "bid" {			//artid, userid, bcn_ammount
		return s.bid(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	artIDs := []string{
		"6b86b273ff34fce19d6b804eff5a3f574",
		"d4735e3a265e16eee03f59718b9b5d030",
		"d16b72230967de01f640b7e4729b49fce",
		"0ce05c1decfe3ad16b72230967de01f64",
	}
	art := []Art{
		Art{ Name: "Wolf", 	Author: "Shimhaq", Description: "A colorful image of a wolf", 	OwnerID: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", BcnValue: "0.13", IsListed: "true", 	AuctionStartTime: "1505591349", AuctionEndTime: "1508184000", HighestBidderId: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", StartingBid: "0.02"},
		Art{ Name: "Owl", 	Author: "Shimhaq", Description: "A colorful image of an owl", 	OwnerID: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", BcnValue: "0.03", IsListed: "false", 	AuctionStartTime: "", AuctionEndTime: "", HighestBidderId: "", StartingBid: ""},
		Art{ Name: "Horse", Author: "Shimhaq", Description: "A colorful image of a horse", 	OwnerID: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", BcnValue: "0.1", IsListed: "true",	AuctionStartTime: "1505591349", AuctionEndTime: "1506184000", HighestBidderId: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", StartingBid: "0.1"},
		Art{ Name: "Lion", 	Author: "Shimhaq", Description: "A colorful image of a lion", 	OwnerID: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", BcnValue: "0.2", IsListed: "true",	AuctionStartTime: "1505591349", AuctionEndTime: "1505684000", HighestBidderId: "1MJY4td932qmFxH6FBQrpnySZBizV67rRC", StartingBid: "0.15"},
	}

	i := 0
	for i < len(art) {
		fmt.Println("i is ", i)
		artAsBytes, _ := json.Marshal(art[i])
		APIstub.PutState(artIDs[i], artAsBytes)
		fmt.Println("Added", art[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryArtById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	artAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(artAsBytes)
}

func (s *SmartContract) list(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	var art = Art{ Name: args[1], Author: args[2], Description: args[3], OwnerID: args[4], BcnValue: args[5], IsListed: args[6]}

	artAsBytes, _ := json.Marshal(art)
	APIstub.PutState(args[0], artAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllArt(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "000000000000000000000000000000000"
	endKey := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllArt:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) endAuction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	artAsBytes, _ := APIstub.GetState(args[0])
	art := Art{}

	json.Unmarshal(artAsBytes, &art)
	if(art.HighestBidderId != ""){
		art.OwnerID = art.HighestBidderId
	}
	art.IsListed = "false"
	art.AuctionStartTime = ""
	art.AuctionEndTime = ""
	art.HighestBidderId = ""
	art.StartingBid = ""

	artAsBytes, _ = json.Marshal(art)
	APIstub.PutState(args[0], artAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) setStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	artAsBytes, _ := APIstub.GetState(args[0])
	art := Art{}

	json.Unmarshal(artAsBytes, &art)
	art.IsListed = args[1]

	artAsBytes, _ = json.Marshal(art)
	APIstub.PutState(args[0], artAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) setPrice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	artAsBytes, _ := APIstub.GetState(args[0])
	art := Art{}

	json.Unmarshal(artAsBytes, &art)
	art.BcnValue = args[1]

	artAsBytes, _ = json.Marshal(art)
	APIstub.PutState(args[0], artAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) setUpAuction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	artAsBytes, _ := APIstub.GetState(args[0])
	art := Art{}

	json.Unmarshal(artAsBytes, &art)
	art.AuctionStartTime = args[1]
	art.AuctionEndTime = args[2]
	art.StartingBid = args[3]

	artAsBytes, _ = json.Marshal(art)
	APIstub.PutState(args[0], artAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) bid(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	artAsBytes, _ := APIstub.GetState(args[0])
	art := Art{}

	json.Unmarshal(artAsBytes, &art)
	art.HighestBidderId = args[1]
	art.BcnValue = args[2]

	artAsBytes, _ = json.Marshal(art)
	APIstub.PutState(args[0], artAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

