package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	typeName       = flag.String("type", "", "name of the type for which we'll generate options")
	prefix         = flag.String("prefix", "With", "prefix in front of each generated option function")
	optionTypeName = flag.String("optionType", "Option", "name of the option type that will be generated")
	output         = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	tagName        = flag.String("tag", "env", "the name of the env tag that should be scanned for, defaults to 'env'")
)

func main() {
	flag.Parse()
	if *typeName == "" {
		flag.Usage()
		os.Exit(2)
	}

	log.SetPrefix("envopts: ")
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run() (err error) {
	var g Generator
	if err := g.Load(); err != nil {
		return fmt.Errorf("failed to load: %v", err)
	}

	if err := g.Generate(); err != nil {
		return fmt.Errorf("failed to generate: %v", err)
	}

	if g.count < 1 {
		return //don't write anything if count is 0
	}

	src, err := g.Format()
	if err != nil {
		return fmt.Errorf("failed to format: %v", err)
	}

	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_opts.go", *typeName)
		outputName = filepath.Join(os.Getenv("GOPATH"), strings.ToLower(baseName))
	}

	if err := ioutil.WriteFile(outputName, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}

	// remove any unused imports using goimports
	cmd := exec.Command(g.giexe, "-w", outputName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to ")
	}

	return
}

// Generator holds code generator state
type Generator struct {
	out      bytes.Buffer
	pkg      *packages.Package
	ts       *types.Struct
	giexe    string
	count    int
	exported bool
}

// Load the package
func (g *Generator) Load() (err error) {
	cfg := packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}

	pkgs, err := packages.Load(&cfg, ".")
	if err != nil {
		return fmt.Errorf("failed to load packages: %w", err)
	}

	if len(pkgs) != 1 {
		return fmt.Errorf("found not 1 packages but %d", len(pkgs))
	}

	g.pkg = pkgs[0]
	g.giexe, err = exec.LookPath("goimports")
	if err != nil {
		return fmt.Errorf("goimports not found in PATH: %w", err)
	}

	return g.lookup()
}

// lookup the type in the packages type definitions
func (g *Generator) lookup() (err error) {
	tobj, ok := g.pkg.Types.Scope().Lookup(*typeName), false
	if tobj == nil {
		return fmt.Errorf("type %s is cannot be found", *typeName)
	}

	g.exported = tobj.Exported()
	g.ts, ok = tobj.Type().Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("type %s is not a struct", *typeName)
	}

	return
}

// Generate the declerations
func (g *Generator) Generate() (err error) {
	fmt.Fprintf(&g.out, "// Code generated by \"stringer %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintf(&g.out, "package %s\n", g.pkg.Name)
	fmt.Fprintf(&g.out, `import "github.com/caarlos0/env/v6"`+"\n")
	fmt.Fprintf(&g.out, "type %s func(*%s)\n", *optionTypeName, *typeName)

	var optionSuffix string
	if *optionTypeName != "Option" {
		optionSuffix = *optionTypeName
	}

	for i := 0; i < g.ts.NumFields(); i++ {
		field := g.ts.Field(i)
		if !field.Exported() {
			continue // skip unexported
		}

		if _, ok := reflect.StructTag(g.ts.Tag(i)).Lookup(*tagName); !ok {
			continue // skip fields without env tag
		}

		g.count++
		fmt.Fprintf(&g.out, `func %s%s%s(v %s) %s {return func(o *%s){o.%s=v}}`+"\n",
			*prefix, field.Name(), optionSuffix, field.Type(), *optionTypeName, *typeName, field.Name())
	}

	firstF, firstA := "F", "A"
	if !g.exported {
		firstF, firstA = "f", "a"
	}

	// generate the function that generation an option func that completely sets the underlying option
	fmt.Fprintf(&g.out, `func %srom%s(v *%s) %s {
		return func(c *%s) { *c = *v }
	}`+"\n", firstF, *typeName, *typeName, *optionTypeName, *typeName)

	// generate the function that parses the environment
	fmt.Fprintf(&g.out, `func %sParseEnv(eo env.Options) (opts []%s, err error) {
		var o %s
		opts = append(opts, %srom%s(&o))
		return opts, env.Parse(&o, eo)
	}`+"\n", optionSuffix, *optionTypeName, *typeName, firstF, *typeName)

	fmt.Fprintf(&g.out, `func %spply%ss(opts ...%s) (res %s) {
		env.Parse(&res, env.Options{Environment: make(map[string]string)})
		for _, o := range opts {
			o(&res)
		}
		return
	}`+"\n", firstA, *optionTypeName, *optionTypeName, *typeName)

	return
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) Format() ([]byte, error) {
	src, err := format.Source(g.out.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return nil, fmt.Errorf("failed to format source: %w", err)
	}

	return src, nil
}
