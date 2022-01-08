package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestGenNormal(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	gen, err := Load("FooEnv", "Option", "With", filepath.Join(wd, "example"))
	require.NoError(t, err)
	require.NotNil(t, gen.pkg)

	t.Run("generate header", func(t *testing.T) {
		require.NoError(t, gen.GenerateHeader())

		out := gen.Output().String()
		require.Contains(t, out, "DO NOT EDIT")
		require.Contains(t, out, "env/v6")
		require.Contains(t, out, `import "time"`)
	})

	t.Run("generate body", func(t *testing.T) {
		require.NoError(t, gen.GenerateBody())

		out := gen.Output().String()
		require.Contains(t, out, `type Option func(*FooEnv)`)
		require.Contains(t, out, `FromFooEnv`)
		require.Contains(t, out, `ParseEnv`)
		require.Contains(t, out, `ApplyOptions`)
		require.Contains(t, out, `func WithFoo(v []env.Options) Option {return func(o *FooEnv){o.Foo=v}}`)
	})

	t.Run("format", func(t *testing.T) {
		src, err := gen.Format()
		require.NoError(t, err)
		require.Contains(t, string(src), `ApplyOptions`)
	})
}

func TestGenPrivSuffix(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	src, err := Generate("fooEnv", "FooOption", "With",
		filepath.Join(wd, "example", "private_foo_env.go"),
		filepath.Join(wd, "example"))
	require.NoError(t, err)

	require.Contains(t, string(src), `type FooOption func(*fooEnv)`)
	require.Contains(t, string(src), `fromFooEnv`)
	require.Contains(t, string(src), `FooOptionParseEnv`)
	require.Contains(t, string(src), `applyFooOptions`)

	checkCompileErrors(t, filepath.Join(wd, "example"))
}

func TestGenPrefix(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	src, err := Generate("RabEnv", "RabEnvOption", "On",
		filepath.Join(wd, "example", "rabenv_opts.go"),
		filepath.Join(wd, "example"))
	require.NoError(t, err)

	require.Contains(t, string(src), `OnHome`)
	require.Contains(t, string(src), `configures:`)
	checkCompileErrors(t, filepath.Join(wd, "example"))
}

func TestGenOnlyPrivate(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	src, err := Generate("BarEnv", "BarEnvOption", "With",
		filepath.Join(wd, "example", "barenv_opts.go"),
		filepath.Join(wd, "example"))
	require.NoError(t, err)
	require.Nil(t, src, "shouldn't generate anything")
}

// checkCompileErrors is a helper that checks for compile errors
func checkCompileErrors(tb testing.TB, pat ...string) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}, pat...)
	require.NoError(tb, err)

	var errs []error
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			errs = append(errs, err)
		}
	})

	require.Len(tb, errs, 0, "at least one error")
}
