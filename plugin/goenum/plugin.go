package goenum

import (
	"strconv"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

var _ generator.Plugin = (*Plugin)(nil)

// New return a new goenum plugin
func New() *Plugin {
	return &Plugin{}
}

// Plugin protoc plugin for generating http1 handlers
type Plugin struct {
	*generator.Generator
	generator.PluginImports
}

// Name returns plugin name
func (p *Plugin) Name() string {
	return "goenum"
}

// Init initialize with generator
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
}

// Generate generate codes
func (p *Plugin) Generate(fd *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	for _, enum := range fd.Enums() {
		p.generateEnum(fd, enum)
	}
}

func (p *Plugin) generateEnum(fd *generator.FileDescriptor, enum *generator.EnumDescriptor) {
	name := enum.GetName()
	p.P("var Enum", name, " = struct {")

	for _, value := range enum.GetValue() {
		p.P("	", strings.TrimPrefix(value.GetName(), name), " ", name)
	}

	p.P("}{")

	for _, value := range enum.GetValue() {
		p.P("	", strings.TrimPrefix(value.GetName(), name), ": ", strconv.Itoa(int(value.GetNumber())), ",")
	}

	p.P("}")
	p.P()
}
