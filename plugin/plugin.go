package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

var _ generator.Plugin = (*Plugin)(nil)

// New return a new gohttp plugin
func New() *Plugin {
	return &Plugin{}
}

// Plugin protoc plugin for generating http1 handlers
type Plugin struct {
	*generator.Generator
	generator.PluginImports

	httpPkg     generator.Single
	grpcPkg     generator.Single
	codesPkg    generator.Single
	metadataPkg generator.Single
}

// Name returns plugin name
func (p *Plugin) Name() string {
	return "gohttp"
}

// Init initialize with generator
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
}

// Generate generate codes
func (p *Plugin) Generate(fd *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	p.httpPkg = p.NewImport("net/http")
	p.grpcPkg = p.NewImport("google.golang.org/grpc")
	p.codesPkg = p.NewImport("google.golang.org/grpc/codes")
	p.metadataPkg = p.NewImport("google.golang.org/grpc/metadata")

	p.generateRefs()
	p.generateInterfaces()

	for _, s := range fd.GetService() {
		p.generateService(s)
	}
}

func (p *Plugin) generateRefs() {

}

func (p *Plugin) generateInterfaces() {
	p.P("// Paramer handle with arguments and outputs")
	p.P("type Paramer interface {")
	p.P("	ParseRequest(req *", p.httpPkg.Use(), ".Request, in ", p.Pkg["proto"], ".Message) error")
	p.P("	HandleResponse(rw ", p.httpPkg.Use(), ".ResponseWriter, out ", p.Pkg["proto"], ".Message, err error)")
	p.P("}")
	p.P()

	p.P("// WebAPIService handle with api registration")
	p.P("type WebAPIService interface {")
	p.P("	Paramer")
	p.P("	Register(pattern string, handler ", p.httpPkg.Use(), ".Handler)")
	p.P("}")
	p.P()
}

func (p *Plugin) generateService(s *descriptor.ServiceDescriptorProto) {
	p.P("// RegisterWebAPI", s.GetName(), "Server register web api methods for ", s.GetName())
	p.P("func RegisterWebAPI", s.GetName(), "Server(s WebAPIService, srv ", s.GetName(), "Server) {")
	for _, m := range s.GetMethod() {
		p.P(p.methodRegisterName(s, m), "(s, srv)")
	}
	p.P("}")
	p.P()

	for _, m := range s.GetMethod() {
		p.generateMethodRegister(s, m)
	}
}

func (p *Plugin) methodRegisterName(s *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) string {
	return "_Register_" + s.GetName() + "_" + m.GetName() + "_Handler"
}

func (p *Plugin) generateMethodRegister(s *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) {
	p.P("func ", p.methodRegisterName(s, m), "(s WebAPIService, srv ", s.GetName(), "Server) {")

	if !m.GetServerStreaming() && !m.GetClientStreaming() {
		inType := p.TypeName(p.ObjectNamed(m.GetInputType()))
		p.RecordTypeUse(m.GetInputType())

		p.P("	s.Register(\"/proto.", s.GetName(), "/", m.GetName(), "\", ", p.httpPkg.Use(), ".HandlerFunc(func(rw ", p.httpPkg.Use(), ".ResponseWriter, req *", p.httpPkg.Use(), ".Request) {")
		p.P("		in := new(", inType, ")")
		p.P("		if err := s.ParseRequest(req, in); err != nil {")
		p.P("			s.HandleResponse(rw, nil, ", p.grpcPkg.Use(), ".Errorf(", p.codesPkg.Use(), ".InvalidArgument, \"parse request into %T: %s\", in, err))")
		p.P("			return")
		p.P("		}")
		p.P()
		p.P("		ctx := ", p.metadataPkg.Use(), ".NewIncomingContext(req.Context(), ", p.metadataPkg.Use(), ".MD(req.Header).Copy())")
		p.P("		out, err := srv.", m.GetName(), "(ctx, in)")
		p.P("		s.HandleResponse(rw, out, err)")
		p.P("	}))")
	}

	p.P("}")
	p.P()
}
