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
package extstub

//go:generate mockgen -source=marshal_stub.go -package=repository -destination=marshal_stub_mocks.go

import (
	"io"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/kbkontrakt/hlfabric-ccdevkit/utils"
)

type (
	// FactoryFunc .
	FactoryFunc func(key string) interface{}

	// VisitFunc .
	VisitFunc func(key string, data interface{}) error

	// MarshalQueryState .
	MarshalQueryState interface {
		GetAllStates(query string, factFunc FactoryFunc, visitFunc VisitFunc) error
	}

	// MarshalState .
	MarshalState interface {
		// WriteState .
		WriteState(key string, value interface{}) error
		// ReadState .
		ReadState(key string, value interface{}) error
	}

	// MarshalPrivState .
	MarshalPrivState interface {
		// WriteState .
		WritePrivState(collection, key string, value interface{}) error
		// ReadState .
		ReadPrivState(collection, key string, value interface{}) error
	}

	// MarshalStub .
	MarshalStub interface {
		shim.ChaincodeStubInterface

		MarshalQueryState
		MarshalState
		MarshalPrivState
	}

	// MarshalStateImpl .
	MarshalStateImpl struct {
		shim.ChaincodeStubInterface

		marshalFunc   utils.MarshalFunc
		unmarshalFunc utils.UnmarshalFunc
	}
)

// GetAllStates .
func (stub *MarshalStateImpl) GetAllStates(query string, factFunc FactoryFunc, visitFunc VisitFunc) error {
	iterator, err := stub.GetQueryResult(query)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return err
		}

		obj := factFunc(kv.GetKey())
		err = stub.unmarshalFunc(kv.GetValue(), obj)
		if err != nil {
			return err
		}

		err = visitFunc(kv.GetKey(), obj)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}

// WriteState .
func (stub *MarshalStateImpl) WriteState(key string, value interface{}) error {
	bytes, err := stub.marshalFunc(value)
	if err != nil {
		return err
	}
	return stub.PutState(key, bytes)
}

// ReadState .
func (stub *MarshalStateImpl) ReadState(key string, value interface{}) error {
	data, err := stub.GetState(key)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNotFound
	}

	err = stub.unmarshalFunc(data, value)
	if err != nil {
		return err
	}

	return nil
}

// WritePrivState .
func (stub *MarshalStateImpl) WritePrivState(collection, key string, value interface{}) error {
	bytes, err := stub.marshalFunc(value)
	if err != nil {
		return err
	}
	return stub.PutPrivateData(collection, key, bytes)
}

// ReadPrivState .
func (stub *MarshalStateImpl) ReadPrivState(collection, key string, value interface{}) error {
	data, err := stub.GetPrivateData(collection, key)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNotFound
	}

	err = stub.unmarshalFunc(data, value)
	if err != nil {
		return err
	}

	return nil
}

// NewMarshalStateImpl .
func NewMarshalStateImpl(stub shim.ChaincodeStubInterface,
	marshalFunc utils.MarshalFunc, unmarshalFunc utils.UnmarshalFunc) MarshalStub {
	return &MarshalStateImpl{
		ChaincodeStubInterface: stub,

		marshalFunc:   marshalFunc,
		unmarshalFunc: unmarshalFunc,
	}
}
