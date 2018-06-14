package webapi

import (
	"net/http"

	"github.com/gogo/protobuf/proto"
)

// Codec is used to encode/decode data in handling http request
type Codec interface {
	In(req *http.Request, in proto.Message) error
	Out(rw http.ResponseWriter, out proto.Message, err error)
}

// Mux represents router multiplexier
type Mux interface {
	Post(pattern string, handler http.Handler)
}

// Interface interface for webapi registration
type Interface interface {
	Codec
	Mux
}
