// Package scaffold contains methods and helpers for code scaffolding
package scaffold

import (
	"fmt"
	"github.com/crdflow/crdflow/pkg/crd"
	"github.com/crdflow/crdflow/pkg/template_funcs"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Option sets options for Scaffold instance
type Option func(s *Scaffold)

// WithOutputLocation allows to specify location where generated files will be located.
// If empty string is provided - then files will be created in current directory.
func WithOutputLocation(location string) Option {
	return func(s *Scaffold) {
		s.location = location
	}
}

// Scaffold contains parameters that required for code scaffolding
type Scaffold struct {
	// location where files will be saved
	location string
}

// New returns instance of a Scaffold
func New(opts ...Option) *Scaffold {
	s := &Scaffold{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

//TODO: convert validation from spec to protobuf
// go-grpc-validator or maybe buf-validate?

// BuildGrpcService is building protobuf package with gRPC CRUD for the provided CustomResource
func (s *Scaffold) BuildGrpcService(customResource crd.CRD) error {
	t := template.
		New(protobufTemplate).
		Funcs(template_funcs.Func)

	tmpl, err := t.Parse(protobufTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := s.prepareFS4ProtobufDefinitions(customResource)
	if err != nil {
		return fmt.Errorf("prepare fs: %w", err)
	}
	defer f.Close()

	if err = tmpl.Execute(f, customResource); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

// prepareFS4ProtobufDefinitions is creating necessary folders in FS
// for correct protobuf modules work and other stuff
func (s *Scaffold) prepareFS4ProtobufDefinitions(customResource crd.CRD) (*os.File, error) {
	apiPath := filepath.Join("api", "crd", strings.ToLower(customResource.Kind), customResource.APIVersion)
	protoFilePath := filepath.Join(s.location, apiPath, strings.ToLower(customResource.Kind)+".proto")

	if err := os.MkdirAll(filepath.Join(s.location, apiPath), fs.ModePerm); err != nil {
		return nil, fmt.Errorf("create `api` folder: %w", err)
	}

	f, err := os.Create(protoFilePath)
	if err != nil {
		return nil, fmt.Errorf("create .proto file: %w", err)
	}

	return f, nil
}

const (
	// package + file name where scaffolded gRPC server will be located,
	// e.g. "location}/server.grpc.go"
	serverStubsPath = "server/grpc.go"
)

// BuildStubs is building protobuf stubs with gRPC CRUD from protobuf definitions of provided CustomResource
func (s *Scaffold) BuildStubs(customResource crd.CRD) error {
	t := template.
		New(serverTemplate).
		Funcs(template_funcs.Func)

	tmpl, err := t.Parse(serverTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := s.prepareFS4Stubs()
	if err != nil {
		return fmt.Errorf("prepare fs: %w", err)
	}
	defer f.Close()

	if err = tmpl.Execute(f, customResource); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

// prepareFS4Stubs is creating necessary folders and files where generated server stubs will be stored
func (s *Scaffold) prepareFS4Stubs() (*os.File, error) {
	if err := os.MkdirAll(filepath.Join(s.location, "server"), fs.ModePerm); err != nil {
		return nil, fmt.Errorf("create `server` folder: %w", err)
	}

	f, err := os.Create(filepath.Join(s.location, serverStubsPath))
	if err != nil {
		return nil, fmt.Errorf("create file with stubs: %w", err)
	}

	return f, nil
}
