package extstub

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"kb-kontrakt.ru/hlfabric/ccdevkit/utils"
)

type (
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

// NewJSONMarshalState .
func NewJSONMarshalState(stub shim.ChaincodeStubInterface) MarshalStub {
	return &MarshalStateImpl{
		ChaincodeStubInterface: stub,

		marshalFunc:   utils.MarshalFuncJSON,
		unmarshalFunc: utils.UnmarshalFuncJSON,
	}
}
