package binding

import "encoding/xml"

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "json"
}

func (xmlBinding) Bind(data []byte, obj interface{}) error {
	if err := xml.Unmarshal(data, obj); err != nil {
		return err
	}
	return defaultValidate.ValidateStruct(obj)
}
