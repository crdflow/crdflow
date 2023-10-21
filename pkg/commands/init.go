package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/crdflow/crdflow/pkg/codegen"
	schema "github.com/crdflow/crdflow/pkg/crd"
	"github.com/crdflow/crdflow/pkg/scaffold"
	"github.com/crdflow/crdflow/pkg/util"
)

const (
	crdFlag        = "crd"
	apiVersionFlag = "version"
	repoNameFlag   = "repo"
	outputDirFlag  = "out"
)

// InitCommand ...
func InitCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "init",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			crd, err := cmd.Flags().GetString(crdFlag)
			if err != nil {
				return err
			}

			apiVersion, err := cmd.Flags().GetString(apiVersionFlag)
			if err != nil {
				return err
			}

			repoName, err := cmd.Flags().GetString(repoNameFlag)
			if err != nil {
				return err
			}

			outputDir, err := cmd.Flags().GetString(outputDirFlag)
			if err != nil {
				return err
			}

			if err = initialGen(ctx, initialGenOptions{
				crd:        crd,
				apiVersion: apiVersion,
				repoName:   repoName,
				output:     outputDir,
			}); err != nil {
				return err
			}

			return nil
		},
	}

	AddStringFlag(command, crdFlag, "", "path to CRD", true)
	AddStringFlag(command, apiVersionFlag, "v1", "api version", false)
	AddStringFlag(command, repoNameFlag, "", "name of the repository", true)
	AddStringFlag(command, outputDirFlag, "", "location to save codegen", true)

	return command
}

type initialGenOptions struct {
	crd        string
	apiVersion string
	repoName   string
	output     string
}

func initialGen(ctx context.Context, opts initialGenOptions) error {
	file, err := os.Open(opts.crd)
	if err != nil {
		return fmt.Errorf("open crd: %w", err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	crdSchema := apiextensionsv1.CustomResourceDefinition{}

	if err = k8syaml.Unmarshal(bytes, &crdSchema); err != nil {
		return fmt.Errorf("unmarshal crd: %w", err)
	}

	var selectedSchema *apiextensionsv1.JSONSchemaProps

	for _, version := range crdSchema.Spec.Versions {
		// TODO: maybe fallback to version that exists in spec?
		if version.Name == opts.apiVersion {
			selectedSchema = version.Schema.OpenAPIV3Schema
		} else {
			//fmt.Printf("CRD version %s not found. Fallback to existing ones...", apiVersion)
			//
			//// fallback to version that exists
			//selectedSchema = version.Schema.OpenAPIV3Schema
			//apiVersion = version.Name
		}

		if selectedSchema == nil {
			return errors.New("there's no spec for provided apiVersion")
		}
	}

	resource := schema.CRD{
		APIVersion: opts.apiVersion,
		Kind:       crdSchema.Spec.Names.Kind,
		Spec:       selectedSchema.Properties["spec"].Properties,
		Repo:       opts.repoName,
		Output:     opts.output,
	}

	sc := scaffold.New(scaffold.WithOutputLocation(opts.output))

	err = sc.BuildGrpcService(resource)
	if err != nil {
		return fmt.Errorf("build grpc service: %w", err)
	}

	log.Println("Create server [y/n]")
	reader := bufio.NewReader(os.Stdin)

	if util.YesNo(reader) {
		// TODO: should be refactored for better readability and simplicity
		err = codegen.GenerateServer(ctx, codegen.GenerateServerOptions{
			RepoName:   opts.repoName,
			ProtoPath:  opts.output + "/api/crd/" + strings.ToLower(crdSchema.Spec.Names.Kind) + "/" + opts.apiVersion,
			OutputPath: opts.output,
			ProtoFile:  strings.ToLower(crdSchema.Spec.Names.Kind) + ".proto",
		})

		if err != nil {
			return fmt.Errorf("generate server: %w", err)
		}
	}

	log.Println("Create server stubs [y/n]")
	reader = bufio.NewReader(os.Stdin)

	if util.YesNo(reader) {
		if err = sc.BuildStubs(resource); err != nil {
			return fmt.Errorf("build server: %w", err)
		}
	}

	return nil
}
