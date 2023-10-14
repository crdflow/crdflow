package cmd

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

var (
	crd        string
	apiVersion string
	repoName   string
	output     string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if err := initialGen(ctx); err != nil {
			return err
		}

		return nil
	},
}

func initialGen(ctx context.Context) error {
	file, err := os.Open(crd)
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
		if version.Name == apiVersion {
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
		APIVersion: apiVersion,
		Kind:       crdSchema.Spec.Names.Kind,
		Spec:       selectedSchema.Properties["spec"].Properties,
		Repo:       repoName,
		Output:     output,
	}

	sc := scaffold.New(output)

	err = sc.BuildGrpcService(resource)
	if err != nil {
		return fmt.Errorf("build grpc service: %w", err)
	}

	log.Println("Create server [y/n]")
	reader := bufio.NewReader(os.Stdin)

	if util.YesNo(reader) {
		// TODO: should be refactored for better readability and simplicity
		err = codegen.GenerateServer(ctx, codegen.GenerateServerOptions{
			RepoName:   repoName,
			ProtoPath:  output + "/api/crd/" + strings.ToLower(crdSchema.Spec.Names.Kind) + "/" + apiVersion,
			OutputPath: output,
			ProtoFile:  strings.ToLower(crdSchema.Spec.Names.Kind) + ".proto",
		})

		if err != nil {
			return fmt.Errorf("gerenate server: %w", err)
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

const (
	crdFlag        = "crd"
	apiVersionFlag = "version"
	repoNameFlag   = "repo"
	outputDirFlag  = "out"
)

func init() {
	initCmd.Flags().StringVar(&crd, crdFlag, "", "path to CRD (required)")
	_ = initCmd.MarkFlagRequired(crdFlag)

	initCmd.Flags().StringVar(&apiVersion, apiVersionFlag, "v1", "api version")

	initCmd.Flags().StringVar(&repoName, repoNameFlag, "", "name of the repository (required)")
	_ = initCmd.MarkFlagRequired(repoNameFlag)

	// TODO: add trailing slash validation
	initCmd.Flags().StringVar(&output, outputDirFlag, ".", "location to save codegen")

	rootCmd.AddCommand(initCmd)
}
