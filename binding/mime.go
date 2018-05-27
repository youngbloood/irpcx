package binding

type MIME string

const (
	MIMEJSON              MIME = "application/json"
	MIMEXML               MIME = "application/xml"
	MIMEXML2              MIME = "text/xml"
	MIMEPOSTForm          MIME = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm MIME = "multipart/form-data"
	MIMEPROTOBUF          MIME = "application/x-protobuf"
)

func Default(mime MIME) Binding {
	switch mime {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	// case MIMEPOSTForm, MIMEMultipartPOSTForm:
	// 	return FORM
	default:
		return JSON
	}
}

var (
	JSON = &jsonBinding{}
	XML  = &xmlBinding{}
	//FORM = &formBinding{}
)
