package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	err := exec.Command("go", "generate", "./example").Run()
	require.NoError(t, err)

	f1, err := os.ReadFile("example/fooenv_opts.go")
	require.NoError(t, err)
	require.NotNil(t, f1)

	f2, err := os.ReadFile("example/private_foo_env.go")
	require.NoError(t, err)
	require.NotNil(t, f2)

	f3, err := os.ReadFile("example/rabenv_opts.go")
	require.NoError(t, err)
	require.NotNil(t, f3)
}
