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

import (
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

const dateTimeFormat = "2006-01-02T15:04:05.9999999Z"

// DateTime .
type DateTime struct {
	time.Time
}

// DateTimeFromTimestamp .
func DateTimeFromTimestamp(timestamp *timestamp.Timestamp) DateTime {
	seconds := timestamp.GetSeconds()
	nanos := timestamp.GetNanos()
	time := time.Unix(seconds, int64(nanos))
	return DateTimeFromStr(time.Format(dateTimeFormat))
}

// ParseDateTime .
func ParseDateTime(str string) (*DateTime, error) {
	time, err := time.Parse(dateTimeFormat, strings.Trim(str, `"`))
	if err != nil {
		return nil, err
	}

	dateTime := DateTime{time}

	return &dateTime, nil
}

// DateTimeFromStr .
func DateTimeFromStr(str string) DateTime {
	dateTime, err := ParseDateTime(str)
	if err != nil {
		panic(err)
	}

	return *dateTime
}

// UnmarshalJSON .
func (dateTime *DateTime) UnmarshalJSON(buf []byte) error {
	time, err := time.Parse(dateTimeFormat, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}

	dateTime.Time = time

	return nil
}

// MarshalJSON .
func (dateTime *DateTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + dateTime.Format(dateTimeFormat) + "\""), nil
}

// DateTimeFromStub .
func DateTimeFromStub(stub shim.ChaincodeStubInterface) (*DateTime, error) {
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, err
	}
	dateTime := DateTimeFromTimestamp(timestamp)
	return &dateTime, nil
}
