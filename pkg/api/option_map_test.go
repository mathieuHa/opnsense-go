package api

import (
	"encoding/json"
	"testing"
)

func TestOptionMapUnmarshal(t *testing.T) {
	var o OptionMap
	if err := json.Unmarshal([]byte(`{"u1":{"value":"a.rules","selected":0},"u2":{"value":"b.rules","selected":1}}`), &o); err != nil {
		t.Fatal(err)
	}
	if k, ok := o.KeyOf("b.rules"); !ok || k != "u2" {
		t.Fatalf("KeyOf(b.rules) = %q, %v", k, ok)
	}
	if _, ok := o.KeyOf("c.rules"); ok {
		t.Fatal("KeyOf(c.rules) found a key")
	}
	if err := json.Unmarshal([]byte(`[]`), &o); err != nil || len(o) != 0 {
		t.Fatalf("empty list: %v %v", o, err)
	}
	if err := json.Unmarshal([]byte(`"x"`), &o); err == nil {
		t.Fatal("expected an error for a string")
	}
}
