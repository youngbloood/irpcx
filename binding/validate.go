package binding

import (
	"reflect"
	"sync"

	validator "gopkg.in/go-playground/validator.v8"
)

type Validator interface {
	ValidateStruct(interface{}) error
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var defaultValidate Validator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		config := &validator.Config{
			TagName: "binding",
		}
		v.validate = validator.New(config)
		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
