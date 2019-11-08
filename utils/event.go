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
package utils

import (
	"errors"
)

type (
	// EventHandler .
	EventHandler = func(args ...interface{}) error

	// Event .
	Event interface {
		On(handler EventHandler)
		Emit(args ...interface{}) error
	}

	// eventImpl .
	eventImpl struct {
		handlers []EventHandler
	}
)

// Emit .
func (ev *eventImpl) Emit(args ...interface{}) error {
	var err error
	for _, handler := range ev.handlers {
		err = handler(args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// On .
func (ev *eventImpl) On(handler EventHandler) {
	ev.handlers = append(ev.handlers, handler)
}

// WrapOnSingleArg .
func WrapOnSingleArg(mFunc MarshalFunc, handler EventHandler) EventHandler {
	return func(args ...interface{}) error {
		if len(args) != 1 {
			return errors.New("failed listen event: not single argument was passed")
		}

		data, err := mFunc(args[0])
		if err != nil {
			panic(err)
		}

		return handler(data)
	}
}

// NewSyncEvent creates default implementation
func NewSyncEvent() Event {
	return &eventImpl{
		make([]EventHandler, 0, 1),
	}
}
