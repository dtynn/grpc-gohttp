package webapi

import (
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/status"
)

type grpcJSONParamer struct {
}

func (j grpcJSONParamer) ParseRequest(req *http.Request, in proto.Message) error {
	defer req.Body.Close()

	return jsonpb.Unmarshal(req.Body, in)
}

func (j grpcJSONParamer) HandleResponse(rw http.ResponseWriter, out proto.Message, err error) {
	if Written(rw) {
		return
	}

	rw.WriteHeader(http.StatusOK)

	m := new(jsonpb.Marshaler)

	if err != nil {
		m.Marshal(rw, status.New(status.Code(err), err.Error()).Proto())
		return
	}

	m.Marshal(rw, out)
}
