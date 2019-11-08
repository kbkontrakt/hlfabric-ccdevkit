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

import "runtime"

// GetStacktrace returns formatted string of stacktrace
func GetStacktrace(isForAll bool) string {
	outBuffer := make([]byte, 1024)

	for {
		n := runtime.Stack(outBuffer, isForAll)
		if n < len(outBuffer) {
			break
		}
		outBuffer = make([]byte, 2*len(outBuffer))
	}

	return string(outBuffer)
}
