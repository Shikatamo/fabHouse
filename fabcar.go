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
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the house structure, with 4 properties.  Structure tags are used by encoding/json library
type House struct {
	Year   string `json:"year"`
	SquareFeets  string `json:"squarefeets"`
	Location string `json:"location"`
	Owner  string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "fabhouse" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabhouse"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryHouse" {
		return s.queryHouse(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createHouse" {
		return s.createHouse(APIstub, args)
	} else if function == "queryAllHouses" {
		return s.queryAllHouses(APIstub)
	} else if function == "changeHouseOwner" {
		return s.changeHouseOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryHouse(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	houseAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(houseAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	houses := []House{
		House{Year: "2007", SquareFeets: "300", Location: "Bayonne", Owner: "Tomoko"},
		House{Year: "1987", SquareFeets: "178", Location: "Anglet", Owner: "Brad"},
		House{Year: "1865", SquareFeets: "37", Location: "Bayonne", Owner: "Jin Soo"},
		House{Year: "1999", SquareFeets: "467", Location: "Anglet", Owner: "Max"},
		House{Year: "2007", SquareFeets: "2534", Location: "Bayonne", Owner: "Adriana"},
		House{Year: "1999", SquareFeets: "205", Location: "purple", Owner: "Michel"},
		House{Year: "2002", SquareFeets: "300", Location: "Biarritz", Owner: "Aarav"},
		House{Year: "2007", SquareFeets: "300", Location: "Biarritz", Owner: "Pari"},
		House{Year: "1989", SquareFeets: "125", Location: "Bayonne", Owner: "Valeria"},
		House{Year: "2007", SquareFeets: "125", Location: "Arruntz", Owner: "Shotaro"},
	}

	i := 0
	for i < len(houses) {
		fmt.Println("i is ", i)
		houseAsBytes, _ := json.Marshal(houses[i])
		APIstub.PutState("HOUSE"+strconv.Itoa(i), houseAsBytes)
		fmt.Println("Added", houses[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createHouse(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var house = House{Year: args[1], SquareFeets: args[2], Location: args[3], Owner: args[4]}

	houseAsBytes, _ := json.Marshal(house)
	APIstub.PutState(args[0], houseAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllHouses(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "HOUSE0"
	endKey := "HOUSE999"

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

	fmt.Printf("- queryAllHouses:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeHouseOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	houseAsBytes, _ := APIstub.GetState(args[0])
	house := House{}

	json.Unmarshal(houseAsBytes, &house)
	house.Owner = args[1]

	houseAsBytes, _ = json.Marshal(house)
	APIstub.PutState(args[0], houseAsBytes)

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
