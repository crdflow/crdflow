package util

import (
	"context"
	"os/exec"
)

func Exec(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	_, err := cmd.Output()
	return err
}
