package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	names := []string{
		filepath.Join(wd, "example", "fooenv_opts.go"),
		filepath.Join(wd, "example", "private_foo_env.go"),
		filepath.Join(wd, "example", "rabenv_opts.go")}
	for _, fname := range names {
		require.NoError(t, os.Remove(fname))
	}

	cmd := exec.Command("go", "generate", "./example")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	require.NoError(t, cmd.Run())
}
