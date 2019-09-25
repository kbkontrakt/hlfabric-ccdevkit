package extstub

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"kb-kontrakt.ru/hlfabric/ccdevkit/utils"
)

type (
	// Response .
	Response interface {
		Send(data interface{}) pb.Response
		SendResponse(data interface{}, err error) pb.Response
		Error(data interface{}) pb.Response
		Success(data interface{}) pb.Response
	}

	// FormatResponseFunc .
	FormatResponseFunc func(data interface{}) interface{}

	// marshalResponse .
	marshalResponse struct {
		marshalFunc   utils.MarshalFunc
		formatError   FormatResponseFunc
		formatSuccess FormatResponseFunc
	}
)

// DefaultFormatError .
func DefaultFormatError(err interface{}) interface{} {
	if err == nil {
		return ""
		// map[string]interface{}{
		// 	"status": "error",
		// }
	}
	return err
	// map[string]interface{}{
	// 	"status":  "error",
	// 	"message": err,
	// }
}

// DefaultFormatSuccess .
func DefaultFormatSuccess(data interface{}) interface{} {
	// if data == nil {
	// return map[string]interface{}{
	// 	"status": "success",
	// }
	// }
	return data
	// return map[string]interface{}{
	// 	"status": "success",
	// 	"data":   data,
	// }
}

func (mr *marshalResponse) Error(data interface{}) pb.Response {
	if data == nil {
		return shim.Error("empty contents of error")
	}

	data = mr.formatError(data)

	bytes, err := mr.marshalFunc(data)
	if err != nil {
		data = mr.formatError(fmt.Sprintf("failed to marshal error response: %+v", err))
		bytes, err = mr.marshalFunc(data)
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to marshal error response after second call: %+v", err))
		}
	}

	return shim.Error(string(bytes))
}

func (mr *marshalResponse) Success(data interface{}) pb.Response {
	data = mr.formatSuccess(data)

	bytes, err := mr.marshalFunc(data)
	if err != nil {
		return mr.Error(err)
	}

	return shim.Success(bytes)
}

// Send .
func (mr *marshalResponse) SendResponse(resp interface{}, err error) pb.Response {
	if err != nil {
		return mr.Error(err)
	}

	return mr.Success(resp)
}

// Send .
func (mr *marshalResponse) Send(resp interface{}) pb.Response {
	if resp == nil {
		return mr.Success(nil)
	}

	if err, ok := resp.(error); ok {
		return mr.Error(err)
	}

	return mr.Success(resp)
}

// NewJSONMarshalResponse .
func NewJSONMarshalResponse() Response {
	return &marshalResponse{
		marshalFunc:   json.Marshal,
		formatError:   DefaultFormatError,
		formatSuccess: DefaultFormatSuccess,
	}
}
