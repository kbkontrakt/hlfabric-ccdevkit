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
package slices

import (
	"testing"
)

func TestCastSamePrimitiveSlices(t *testing.T) {
	type strType string

	src := []strType{}
	dst := []string{}

	CastSamePrimitiveSlices(&dst, &src)

	if len(src) != 0 || len(dst) != 0 {
		t.Fatalf("cast result of src and dst length is not equal 0")
	}

	src = []strType{"str1", "str2"}

	CastSamePrimitiveSlices(&dst, &src)

	if len(dst) != len(src) || len(src) != 2 {
		t.Fatalf("cast result of src and dst length is not equal 2")
	}

	if src[0] != strType(dst[0]) || src[1] != strType(dst[1]) {
		t.Fatalf("cast result is not the same, src != dst")
	}
}
