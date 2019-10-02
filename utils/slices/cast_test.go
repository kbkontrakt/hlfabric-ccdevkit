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
