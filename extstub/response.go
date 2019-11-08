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
		data = "empty contents of error"
	}
	if err, ok := data.(error); ok {
		data = err.Error()
	}

	data = mr.formatError(data)

	bytes, err := mr.marshalFunc(data)
	if err != nil {
		data = mr.formatError(fmt.Sprintf("failed to marshal [%s] error response: %+v", data, err))
		bytes, err = mr.marshalFunc(data)
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to marshal [%s] error response after second call: %+v", data, err))
		}
	}

	return shim.Error(string(bytes))
}

func (mr *marshalResponse) Success(data interface{}) pb.Response {
	if bytes, ok := data.([]byte); ok {
		return shim.Success(bytes)
	}

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
