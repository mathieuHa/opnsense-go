package api

import "encoding/json"

type OptionMap map[string]string

func (o *OptionMap) UnmarshalJSON(data []byte) error {
	var m map[string]struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		var empty []any
		if listErr := json.Unmarshal(data, &empty); listErr == nil {
			*o = OptionMap{}
			return nil
		}
		return err
	}
	out := OptionMap{}
	for k, v := range m {
		out[k] = v.Value
	}
	*o = out
	return nil
}

func (o OptionMap) KeyOf(value string) (string, bool) {
	for k, v := range o {
		if v == value {
			return k, true
		}
	}
	return "", false
}
