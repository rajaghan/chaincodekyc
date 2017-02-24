/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var customIndexStr = "_marbleindex"				//name for the key/value that will store a list of all known customers

type Customer struct{
	Nino string `json:"nino"`					//the fieldtags are needed to keep case from bouncing around
	Title string `json:"title"`
	FirstName string `json:"firstname"`
	MiddleName string `json:"middlename"`
	LastName string `json:"lastname"`					//the fieldtags are needed to keep case from bouncing around
	DOB string `json:"dob"`
	Gender string `json:"gender"`
	RS string `json:"rs"`
	Address string `json:"address"`					//the fieldtags are needed to keep case from bouncing around
	Email string `json:"email"`
	Landline string `json:"landline"`
	Mobile string `json:"mobile"`
	PC string `json:"PC"`					//the fieldtags are needed to keep case from bouncing around
	Ppnum string `json:"ppnum"`
	Dlnum string `json:"dlnum"`
	Title_Conf string `json:"title_conf"`
	FN_Conf string `json:"fn_conf"`
	MN_Conf string `json:"mn_conf"`
	LN_Conf string `json:"ln_conf"`
	DOB_Conf string `json:"dob_conf"`
	Gender_Conf string `json:"gender_conf"`
	RS_Conf string `json:"RS_conf"`
	Address_Conf string `json:"address_conf"`
	Email_Conf string `json:"email_conf"`
	Landline_Conf string `json:"landline_conf"`
	Mobile_Conf string `json:"mobile_conf"`
	PC_Conf string `json:"PC_conf"`
	Ppnum_Conf string `json:"ppnum_conf"`
	Dlnum_Conf string `json:"dlnum_conf"`
}


// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(customIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "delete" {										//deletes an entity from its state
		return t.Delete(stub, args)
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "init_customer" {									//create a new customer
		return t.init_customer(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the marble index
	customersAsBytes, err := stub.GetState(customIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get customer index")
	}
	var customerIndex []string
	json.Unmarshal(customersAsBytes, &customerIndex)								//un stringify it aka JSON.parse()
	
	//remove customer from index
	for i,val := range customerIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															//find the correct customer
			fmt.Println("found customer")
			customerIndex = append(customerIndex[:i], customerIndex[i+1:]...)			//remove it
			for x:= range customerIndex{											//debug prints...
				fmt.Println(string(x) + " - " + customerIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(customerIndex)									//save new index
	err = stub.PutState(customIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) init_customer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	//if len(args) != 4 {
		//return nil, errors.New("Incorrect number of arguments. Expecting 4")
	//}

	fmt.Println("- start init customer")
		
	str := `{"nino": "` + args[0] + `", "title": "` + args[1] + `", "firstname": "` + args[2] + `", "middlename": "` + args[3] + `" ,"lastname": "` + args[4] + `","dob": "` + args[5] + `", "gender": "` + args[6] + `", "rs": "` + args[7] + `","address": "` + args[8] + `","email": "` + args[9] + `", "landline": "` + args[10] + `", "mobile": "` + args[11] + `", "PC": "` + args[12] + `", "ppnum": "` + args[13] + `", "dlnum": "` + args[14]+`","title_conf": "` + args[15]+`","fn_conf": "` + args[16]+`","mn_conf": "` + args[17]+`","ln_conf": "` + args[18]+`","dob_conf": "` + args[19]+`","gender_conf": "` + args[20]+`","RS_conf": "` + args[21]+`","address_conf": "` + args[22]+`","email_conf": "` + args[23]+`","landline_conf": "` + args[24]+`","mobile_conf": "` + args[25]+`","PC_conf": "` + args[26]+`","ppnum_conf": "` + args[27]+`","dlnum_conf": "` + args[28]+`" }`
	
	err = stub.PutState(args[0], []byte(str))								//store marble with id as key
	
	if err != nil {
		return nil, err
	}
		
	//get the customer index
	customersAsBytes, err := stub.GetState(customIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get customer index")
	}
	var customerIndex []string
	json.Unmarshal(customersAsBytes, &customerIndex)							//un stringify it aka JSON.parse()
	
	//append
	customerIndex = append(customerIndex, args[0])								//add customer  nino to index list
	fmt.Println("! Customer index: ", customerIndex)
	jsonAsBytes, _ := json.Marshal(customerIndex)
	err = stub.PutState(customIndexStr, jsonAsBytes)						//store nino of customer

	fmt.Println("- end init customer")
	return nil, nil
}

