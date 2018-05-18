package types

import (
	"net/http"

	"github.com/gogo/protobuf/proto"
)

// Server for http method registration
type Server interface {
	ParseRequest(req *http.Request, in proto.Message) error
	HandleResponse(rw http.ResponseWriter, out proto.Message, err error)
	Register(pattern string, handler http.Handler)
}
