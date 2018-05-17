package webapi

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type grpcJSONParamer struct {
}

func (j grpcJSONParamer) ParseRequest(req *http.Request, in interface{}) error {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(in)
}

func (j grpcJSONParamer) HandleResponse(rw http.ResponseWriter, out interface{}) {
	if Written(rw) {
		return
	}

	if out == nil {
		return
	}

	var res Result
	res.Res.Code = int32(codes.OK)

	switch o := out.(type) {
	case error:
		res.Res.Code = int32(grpc.Code(o))
		res.Res.Message = o.Error()

	default:
		res.Data = out
	}

	data, err := json.Marshal(res)
	if err != nil {
		var mres Result
		mres.Res.Code = int32(codes.Internal)
		mres.Res.Message = err.Error()

		data, _ = json.Marshal(mres)
	}

	rw.Write(data)
	rw.WriteHeader(http.StatusOK)
}
