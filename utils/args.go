/*
 *  Copyright 2017 - 2019 KB Kontrakt LLC - All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetArgsFromTransientMap returns function name and args from transientMap
func GetArgsFromTransientMap(stub shim.ChaincodeStubInterface) (fn string, args []string, err error) {
	args = []string{}

	tmap, err := stub.GetTransient()
	if err != nil {
		return
	}

	data, exists := tmap["Args"]
	if !exists {
		return
	}

	rawArgsItems := []interface{}{}

	err = json.Unmarshal(data, &rawArgsItems)
	if err != nil {
		return
	}

	if len(rawArgsItems) == 0 {
		return
	}

	bytes := []byte{}
	for inx, arg := range rawArgsItems {
		argStr := ""
		if arg != nil {
			if val, ok := arg.(string); ok {
				argStr = val
			} else {
				bytes, err = json.Marshal(arg)
				if err != nil {
					return
				}
				argStr = string(bytes)
			}
		}

		if inx == 0 {
			fn = argStr
		} else {
			args = append(args, argStr)
		}
	}

	return
}

// GetFnArgsOrFromTransientMap returns function name and args from args else transientMap
func GetFnArgsOrFromTransientMap(stub shim.ChaincodeStubInterface) (fn string, args []string, err error) {
	fn, args = stub.GetFunctionAndParameters()

	if (len(fn) == 0 || fn == "*") && len(args) == 0 {
		fn, args, err = GetArgsFromTransientMap(stub)
	}

	return
}
