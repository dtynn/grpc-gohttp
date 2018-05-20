package types

import (
	"net/http"

	"github.com/gogo/protobuf/proto"
)

// Paramer handling request and response
type Paramer interface {
	ParseRequest(req *http.Request, in proto.Message) error
	HandleResponse(rw http.ResponseWriter, out proto.Message, err error)
}

// Server for http method registration
type Server interface {
	Paramer
	Register(pattern string, handler http.Handler)
}
