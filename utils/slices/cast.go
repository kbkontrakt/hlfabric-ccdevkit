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
