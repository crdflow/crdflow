package crd

import apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

type CRD struct {
	APIVersion string
	Kind       string
	Repo       string
	Output     string

	Spec map[string]apiextensionsv1.JSONSchemaProps
}
