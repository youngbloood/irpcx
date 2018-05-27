package binding

type Binding interface {
	Name() string
	Bind([]byte, interface{}) error
}
