generate:
	protoc -I ./example/proto -I $(GOPATH)/src --go_out=plugins=grpc:./example/proto/ --gohttp_out=./example/proto/ ./example/proto/*.proto
