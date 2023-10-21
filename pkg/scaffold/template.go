package scaffold

// protobufTemplate is a go-template to build gRPC service proto definition
var protobufTemplate = `syntax = "proto3";

package crd.{{ lower .Kind }}.{{ .APIVersion }};

option go_package = "{{ .Repo }}/gen/go/crd/{{ lower .Kind }}/{{ .APIVersion }};{{ lower .Kind }}{{ .APIVersion }}";

service {{ .Kind }}Service {
  rpc Create{{ .Kind }}(Create{{ .Kind }}Request) returns (Create{{ .Kind }}Response) {}
}

{{/* TODO: improve readability of template, looks messy */ -}}
message Create{{ .Kind }}Request {
  {{ $index := 1 }}

  {{- range $key, $value := .Spec -}}

  // {{ $value.Description }}
  {{ $value.Type }} {{ $key }} = {{ $index }};

  {{- $index = inc $index }}
  {{- end }}
}

message Create{{ .Kind }}Response {}
`

// serverTemplate is a go-template to build gRPC server implementation (optional)
var serverTemplate = `// Code generated by crdflow, DO NOT EDIT.

// Package server contains various functions to control Kubernetes CRDs via gRPC interface
package server

import (
	"context"
	pb "{{ .Repo }}/{{ .Output }}/gen"
	api "{{ .Repo }}/api/{{ .APIVersion }}/{{ lower .Kind }}"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Server contains initialized k8s client that used to make API calls to Kubernetes itself.
type Server struct {
	pb.UnimplementedExampleServiceServer

	kube client.Client
}

// Create{{ .Kind }} is creating custom resource in kubernetes cluster
// using initialized client that was provided in New()
func (s *Server) Create{{ .Kind }}(ctx context.Context, request *pb.Create{{ .Kind }}Request) (*pb.Create{{ .Kind }}Response, error) {
	if err := s.kube.Create(ctx, &api.{{ .Kind }}{
		Spec: api.{{ .Kind }}Spec {},
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Create{{ .Kind }}Response{}, nil
}

// New returns pointer to the Server instance and adding CRD to runtime scheme
// to be able to control it via client-go
func New(kubeClient client.Client) *Server {
	scheme := runtime.NewScheme()
	_ = api.AddToScheme(scheme)

	return &Server{
		kube: kubeClient,
	}
}
`
