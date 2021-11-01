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

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type structBasedKey struct {
	IDsmall *string `json:"id"`
	IDbig   *string `json:"ID"`
}

func enumKeysFromJSONArgs(args []string, limit int) ([]string, error) {
	var err error
	keys := make([]string, 0, limit)

	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}
		if arg[0] == '"' {
			key := ""
			err = json.Unmarshal([]byte(arg), &key)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
			continue
		}

		structs := []structBasedKey{}
		if arg[0] == '[' {
			err = json.Unmarshal([]byte(arg), &structs)
		} else {
			keyStruct := structBasedKey{}
			err = json.Unmarshal([]byte(arg), &keyStruct)
			structs = append(structs, keyStruct)
		}

		for i, l := 0, len(structs); i < l && limit >= 0; i, limit = i+1, limit-1 {
			if structs[i].IDsmall != nil {
				keys = append(keys, *structs[i].IDsmall)
			} else if structs[i].IDbig != nil {
				keys = append(keys, *structs[i].IDbig)
			}
		}
	}

	return keys, nil
}

func enumKeysFromIterator(iter shim.StateQueryIteratorInterface, limit int) ([]string, error) {
	keys := make([]string, 0, limit)

	for ; iter.HasNext() && limit >= 0; limit-- {
		pair, err := iter.Next()

		if err != nil {
			return nil, err
		}
		keys = append(keys, pair.GetKey())
	}

	return keys, nil
}
