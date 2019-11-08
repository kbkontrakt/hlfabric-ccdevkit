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
	"reflect"
	"unsafe"
)

// CastSamePrimitiveSlices .
func CastSamePrimitiveSlices(sliceDst, sliceSrc interface{}) {
	if sliceSrc == nil || sliceDst == nil {
		panic("cast same primitive slices failed: passed nil src or dst")
	}

	dst := (*reflect.SliceHeader)((unsafe.Pointer(reflect.ValueOf(sliceDst).Pointer())))
	src := (*reflect.SliceHeader)((unsafe.Pointer(reflect.ValueOf(sliceSrc).Pointer())))

	dst.Data = src.Data
	dst.Len = src.Len
	dst.Cap = src.Cap
}
