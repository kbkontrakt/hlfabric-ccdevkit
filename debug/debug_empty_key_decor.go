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
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var debugLog = shim.NewLogger("debug_stubs")

type debugGetStateByEmptyKey struct {
	shim.ChaincodeStubInterface
}

func (s *debugGetStateByEmptyKey) GetState(key string) ([]byte, error) {
	if key == "" {
		debugLog.Debugf("[DEBUG] Catch! GetState with empty key in [%s]\n", GetStacktrace(false))
	}
	return s.ChaincodeStubInterface.GetState(key)
}

func (s *debugGetStateByEmptyKey) GetPrivateData(collectionName, key string) ([]byte, error) {
	if key == "" {
		debugLog.Debugf("[DEBUG] Catch! GetPrivateData with empty key in [%s]\n", GetStacktrace(false))
	}
	return s.ChaincodeStubInterface.GetPrivateData(collectionName, key)
}

// NewDebugGetStateByEmptyKey notify places of passed empty keys with stacktrace.
func NewDebugGetStateByEmptyKey(stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &debugGetStateByEmptyKey{
		ChaincodeStubInterface: stub,
	}
}
