package main

import (
	"context"
	"io"

	"github.com/dtynn/grpc-gohttp/example/proto"
	"github.com/dtynn/grpc-gohttp/pkg/grpcapi"
	"github.com/dtynn/grpc-gohttp/pkg/mixapi"
	"github.com/dtynn/grpc-gohttp/pkg/webapi"
)

func main() {
	g := grpcapi.New()
	w := webapi.New()
	m := mixapi.Mix(g, w)

	srv := &echoServer{}
	proto.RegisterEchoServer(g.Server, srv)
	proto.RegisterWebAPIEchoServer(w, srv)

	m.Listen(":10181")
}

type echoServer struct {
}

func (e *echoServer) Ping(ctx context.Context, in *proto.In) (*proto.Out, error) {
	return &proto.Out{
		Typ: in.GetTyp(),
		Msg: "out " + in.GetMsg(),
		Num: in.GetNum() + 1,
	}, nil
}

func (e *echoServer) Stream(stream proto.Echo_StreamServer) error {
	var count int64
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&proto.Out{
			Typ: in.GetTyp(),
			Msg: "out " + in.GetMsg(),
			Num: count,
		}); err != nil {
			return err
		}

		count++
	}

	return nil
}
