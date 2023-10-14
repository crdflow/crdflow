package scaffold

import (
	"fmt"
	"github.com/crdflow/crdflow/pkg/crd"
	"github.com/crdflow/crdflow/pkg/template_funcs"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

type Scaffold struct {
	SaveLocation string
}

func New(saveLocation string) *Scaffold {
	return &Scaffold{SaveLocation: saveLocation}
}

//TODO: convert validation from spec to protobuf
// go-grpc-validator or maybe buf-validate?

func (s *Scaffold) BuildGrpcService(data crd.CRD) error {
	t := template.
		New(protobufTmpl).
		Funcs(template_funcs.Func)

	tmpl, err := t.Parse(protobufTmpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	apiPath := "api/crd/" + strings.ToLower(data.Kind) + "/" + data.APIVersion

	if err = os.MkdirAll(s.SaveLocation+"/"+apiPath, fs.ModePerm); err != nil {
		return fmt.Errorf("create `api` folder: %w", err)
	}

	f, err := os.Create(s.SaveLocation + "/" + apiPath + "/" + strings.ToLower(data.Kind) + ".proto")
	if err != nil {
		return fmt.Errorf("create `.proto` file: %w", err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err = tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func (s *Scaffold) BuildStubs(data crd.CRD) error {
	t := template.
		New(serverTmpl).
		Funcs(template_funcs.Func)

	tmpl, err := t.Parse(serverTmpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err = os.MkdirAll(s.SaveLocation+"/"+"server", fs.ModePerm); err != nil {
		return fmt.Errorf("create `server` folder: %w", err)
	}

	f, err := os.Create(s.SaveLocation + "/" + "server/grpc.go")
	if err != nil {
		return fmt.Errorf("create `.proto` file: %w", err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err = tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
