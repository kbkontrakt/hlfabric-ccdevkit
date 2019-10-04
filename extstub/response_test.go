package extstub

import (
	"errors"
	"testing"
)

func TestResponseSendResponseShouldFormatError(t *testing.T) {
	resp := NewJSONMarshalResponse()

	result := resp.Error("error")
	if result.Status != 500 {
		t.Fatalf("expected error, got %d", result.GetStatus())
	}
	if result.GetMessage() != `"error"` {
		t.Fatalf("expected \"error\", got %s", result.GetMessage())
	}

	result = resp.SendResponse(nil, errors.New("error1"))

	if result.Status != 500 {
		t.Fatalf("expected error, got %d", result.GetStatus())
	}
	if result.GetMessage() != `"error1"` {
		t.Fatalf("expected \"error\", got %s", result.GetMessage())
	}
}
