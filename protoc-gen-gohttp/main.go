package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"

	"github.com/dtynn/grpc-gohttp/plugin"
)

func main() {
	gen := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, gen.Request); err != nil {
		gen.Error(err, "parsing input proto")
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("no files to generate")
	}

	gen.CommandLineParameters(gen.Request.GetParameter())

	gen.WrapTypes()
	gen.SetPackageNames()
	gen.BuildTypeNameMap()
	gen.GeneratePlugin(plugin.New())

	for i := range gen.Response.File {
		gen.Response.File[i].Name = proto.String(strings.Replace(gen.Response.File[i].GetName(), ".pb.go", ".gohttp.pb.go", -1))
	}

	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "marshaling response")
	}

	if _, err := os.Stdout.Write(data); err != nil {
		gen.Error(err, "writing output proto")
	}
}
