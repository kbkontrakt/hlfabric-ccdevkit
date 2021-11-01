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
package debug

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/sirupsen/logrus"
)

// Invoke processes debug methods
func Invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	const lockKey = "debuglock_"
	const batchEnumSize = 2048

	// logger := logrus.NewLogger("debugtools")

	logger := logrus.New()
	lockData, err := stub.GetState(lockKey)
	if err != nil {
		return nil, err
	}
	if lockData != nil {
		return nil, errors.New("unsupported function")
	}

	tmap, err := stub.GetTransient()
	if tmap != nil && len(tmap["args"]) > 0 {
		argsInfs := []interface{}{}
		err = json.Unmarshal(tmap["args"], &argsInfs)
		if err != nil {
			return nil, err
		}
		args = []string{}
		for _, arg := range argsInfs {
			switch val := arg.(type) {
			case string:
				args = append(args, val)
			default:
				bytes, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				args = append(args, string(bytes))
			}
		}
	}

	var isForPrivateData = false
	var delKeysStartSliceIndex = 2

	switch args[0] {
	case "DelPrivState":
		isForPrivateData = true
		delKeysStartSliceIndex = 3
		fallthrough

	case "DelState":
		if len(args) < delKeysStartSliceIndex+1 {
			return nil, errors.New("not enough arguments")
		}
		var err error
		var keys []string

		if args[1] == "query" {
			var iterator shim.StateQueryIteratorInterface
			if isForPrivateData {
				// iterator, err = stub.GetPrivateDataQueryResult(args[2], args[3])
				panic("it is not allowed to query and update in the same transaction")
			} else {
				iterator, err = stub.GetQueryResult(args[2])
			}
			if err != nil {
				return nil, err
			}
			keys, err = enumKeysFromIterator(iterator, batchEnumSize)
		} else {
			keys, err = enumKeysFromJSONArgs(args[delKeysStartSliceIndex:], batchEnumSize)
		}
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			if isForPrivateData {
				err = stub.DelPrivateData(args[2], key)
			} else {
				err = stub.DelState(key)
			}

			if err != nil {
				logger.Errorf("failed to delete state of [%+v]: %+v", key, err)
			}
		}

	case "Lock":
		err = stub.PutState(lockKey, []byte("1"))
		break

	default:
		return nil, errors.New("unknown debug method")
	}

	return nil, nil
}
