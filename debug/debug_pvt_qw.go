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

type debugPrivateQueryPutData struct {
	shim.ChaincodeStubInterface
	previousStack  string
	previousAction int
}

const (
	debugPrivateNoAction = 0
	debugPrivateQuery    = 1
	debugPrivateWrite    = 2
)

func (s *debugPrivateQueryPutData) checkAction(actionType int, collectionName, arg string) {
	if s.previousAction != debugPrivateNoAction && s.previousAction != actionType {
		if actionType == debugPrivateQuery {
			debugLog.Debugf("[DEBUG] Catch! Make a query after write for [%s] with [%s], stacktrace [%s] previous [%s]\n",
				collectionName, arg, GetStacktrace(false), s.previousStack)
		} else if actionType == debugPrivateWrite {
			debugLog.Debugf("[DEBUG] Catch! Make a write after query for [%s] with [%s], stacktrace [%s] previous [%s]\n",
				collectionName, arg, GetStacktrace(false), s.previousStack)
		}
	}

	s.previousAction = actionType
	s.previousStack = GetStacktrace(false)
}

// @TODO: Private Partial, Composite Key

func (s *debugPrivateQueryPutData) GetPrivateDataQueryResult(collectionName, query string) (shim.StateQueryIteratorInterface, error) {
	s.checkAction(debugPrivateQuery, collectionName, query)
	return s.ChaincodeStubInterface.GetPrivateDataQueryResult(collectionName, query)
}

func (s *debugPrivateQueryPutData) PutState(key string, data []byte) error {
	s.checkAction(debugPrivateWrite, "", key)
	return s.ChaincodeStubInterface.PutState(key, data)
}

func (s *debugPrivateQueryPutData) DelState(key string) error {
	s.checkAction(debugPrivateWrite, "", key)
	return s.ChaincodeStubInterface.DelState(key)
}

func (s *debugPrivateQueryPutData) PutPrivateData(collectionName, key string, data []byte) error {
	s.checkAction(debugPrivateWrite, collectionName, key)
	return s.ChaincodeStubInterface.PutPrivateData(collectionName, key, data)
}

func (s *debugPrivateQueryPutData) DelPrivateData(collectionName, key string) error {
	s.checkAction(debugPrivateWrite, collectionName, key)
	return s.ChaincodeStubInterface.DelPrivateData(collectionName, key)
}

// NewDebugPrivateQueryPutData notify places of passed empty keys with stacktrace.
func NewDebugPrivateQueryPutData(stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &debugPrivateQueryPutData{
		ChaincodeStubInterface: stub,

		previousAction: debugPrivateNoAction,
	}
}
