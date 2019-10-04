package chrono

import (
	"encoding/json"
	"testing"
)

func TestDateTimeMarshalJSON(t *testing.T) {
	const strDate = "2006-01-02T15:04:05.999999999+07:00"
	dt := DateTimeFromStr(strDate)

	bytes, err := json.Marshal(dt)
	if err != nil {
		t.Fatal("expected no errors, got", err)
	}

	if string(bytes) != `"`+strDate+`"` {
		t.Fatal("expected formated datetime, got", string(bytes))
	}
}
