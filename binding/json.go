package binding

import (
	"encoding/json"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(data []byte, obj interface{}) error {
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}
	return defaultValidate.ValidateStruct(obj)
}
