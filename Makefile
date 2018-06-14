build:
	go install ./...

clean:
	rm -f ./example/proto/*.pb.go

generate: clean
	protoc -I ./example/proto -I $(GOPATH)/src --go_out=plugins=grpc:./example/proto/ --gohttp_out=./example/proto/ --goenum_out=./example/proto/ ./example/proto/*.proto
