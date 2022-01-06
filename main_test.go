package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	names := []string{"example/fooenv_opts.go", "example/private_foo_env.go", "example/rabenv_opts.go"}
	for _, fname := range names {
		require.NoError(t, os.Remove(fname))
	}

	cmd := exec.Command("go", "generate", "./example")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err)

	for _, fname := range names {
		_, err := os.ReadFile(fname)
		require.NoError(t, err)
	}
}
