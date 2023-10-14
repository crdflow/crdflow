package codegen

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/crdflow/crdflow/pkg/util"
)

type GenerateServerOptions struct {
	RepoName   string
	ProtoPath  string
	OutputPath string
	ProtoFile  string
}

// GenerateServer gRPC server.
func GenerateServer(ctx context.Context, opts GenerateServerOptions) error {
	outputGen := filepath.Join(opts.OutputPath, "gen")
	if err := os.MkdirAll(outputGen, fs.ModePerm); err != nil {
		return err
	}

	if err := generateProto(ctx, opts.ProtoPath, outputGen, opts.ProtoFile); err != nil {
		return err
	}

	if err := goModInit(ctx, opts.OutputPath, opts.RepoName); err != nil {
		return err
	}

	return goModTidy(ctx, opts.OutputPath)
}

func generateProto(ctx context.Context, protoPath, outputGen, protoFile string) error {
	return util.Exec(
		ctx,
		"",
		"protoc",
		"--proto_path", protoPath,
		"--go_out", outputGen,
		"--go_opt=paths=source_relative",
		"--go-grpc_out", outputGen,
		"--go-grpc_opt=paths=source_relative",
		protoFile,
	)
}

func goModInit(ctx context.Context, outputPath, repoName string) error {
	return util.Exec(
		ctx,
		outputPath+"/",
		"go", "mod", "init", repoName,
	)
}

func goModTidy(ctx context.Context, outputPath string) error {
	return util.Exec(
		ctx,
		outputPath+"/",
		"go", "mod", "tidy",
	)
}
