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
package logs

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// dummyLogger .
type dummyLogger struct{}

// SetLevel .
func (c *dummyLogger) SetLevel(level shim.LoggingLevel) {}

// IsEnabledFor .
func (c *dummyLogger) IsEnabledFor(level shim.LoggingLevel) bool {
	return true
}

// Debug .
func (c *dummyLogger) Debug(args ...interface{}) {}

// Info .
func (c *dummyLogger) Info(args ...interface{}) {}

// Notice .
func (c *dummyLogger) Notice(args ...interface{}) {}

// Warning .
func (c *dummyLogger) Warning(args ...interface{}) {}

// Error .
func (c *dummyLogger) Error(args ...interface{}) {}

// Critical .
func (c *dummyLogger) Critical(args ...interface{}) {}

// Debugf .
func (c *dummyLogger) Debugf(format string, args ...interface{}) {}

// Infof .
func (c *dummyLogger) Infof(format string, args ...interface{}) {}

// Noticef .
func (c *dummyLogger) Noticef(format string, args ...interface{}) {}

// Warningf .
func (c *dummyLogger) Warningf(format string, args ...interface{}) {}

// Errorf .
func (c *dummyLogger) Errorf(format string, args ...interface{}) {}

// Criticalf .
func (c *dummyLogger) Criticalf(format string, args ...interface{}) {}

// DummyLogger .
func DummyLogger() Logger {
	return &dummyLogger{}
}
