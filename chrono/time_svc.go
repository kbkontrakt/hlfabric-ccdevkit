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
package chrono

//go:generate mockgen -source=time_svc.go -package=chrono -destination=time_svc_mocks.go

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	// TimeService .
	TimeService interface {
		NowDateTime() (DateTime, error)
	}

	timeServiceImpl struct {
		stub shim.ChaincodeStubInterface
	}
)

func (ts *timeServiceImpl) NowDateTime() (dt DateTime, err error) {
	var txTime *timestamp.Timestamp

	txTime, err = ts.stub.GetTxTimestamp()
	if err != nil {
		return
	}

	return DateTimeFromTimestamp(txTime), nil
}

// NewTimeServiceImpl creates default implementation
func NewTimeServiceImpl(stub shim.ChaincodeStubInterface) TimeService {
	return &timeServiceImpl{
		stub,
	}
}
