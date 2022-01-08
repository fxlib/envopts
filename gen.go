package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// Generate will generate te source code for the provided typeName and optionTypeName
func Generate(typeName, optTypeName, prefix string, dst string, pat ...string) (src []byte, err error) {
	gen, err := Load(typeName, optTypeName, prefix, pat...)
	if err != nil {
		return nil, fmt.Errorf("failed to load source code: %w", err)
	}

	if err = gen.GenerateHeader(); err != nil {
		return nil, fmt.Errorf("failed to generate header: %w", err)
	}

	if err = gen.GenerateBody(); err != nil {
		return nil, fmt.Errorf("failed to generate body: %w", err)
	}

	if gen.resCount < 1 {
		return nil, nil //noting to generate for private field structs
	}

	src, err = gen.Format()
	if err != nil {
		return src, fmt.Errorf("failed to format generated code: %w", err)
	}

	if err = os.WriteFile(dst, src, 0644); err != nil {
		return src, fmt.Errorf("failed to write output file: %w", err)
	}

	gicmd := exec.Command("goimports", "-w", dst)
	gicmd.Stderr = os.Stderr
	gicmd.Stdout = os.Stdout
	if err = gicmd.Run(); err != nil {
		return src, fmt.Errorf("failed to run goimports: %w", err)
	}

	return src, nil
}

// Gen is the generator
type Gen struct {
	out  bytes.Buffer
	pkg  *packages.Package
	opts struct {
		typeName    string
		optTypeName string
		optSuffix   string
		prefix      string
	}

	// resCount counts how many options in the end were generated. If non got
	// generated we'll skip the writing of source code altogether.
	resCount int
}

// Load will setup the generator
func Load(typeName, optTypeName, prefix string, pat ...string) (g *Gen, err error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}, pat...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("require exactly one package, got '%d': %w", len(pkgs), err)
	}

	g = &Gen{pkg: pkgs[0]}
	g.opts.typeName = typeName
	g.opts.optTypeName = optTypeName
	g.opts.prefix = prefix

	// if the option type name is not default we need to suffix other identifiers as well
	if g.opts.optTypeName != "Option" {
		g.opts.optSuffix = g.opts.optTypeName
	}

	return
}

// GenerateHeader will generate the package and import statements
func (g *Gen) GenerateHeader() (err error) {
	fmt.Fprintf(&g.out, "// Code generated by \"envopts %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintf(&g.out, "package %s\n", g.pkg.Name)
	fmt.Fprintf(&g.out, `import "github.com/caarlos0/env/v6"`+"\n")

	for _, file := range g.pkg.Syntax {
		w := Walk{af: file, tf: g.pkg.Fset.File(file.Pos()), Gen: g}
		ast.Inspect(file, w.GenerateImports)
	}

	return
}

// GenerateBody will generate the rest of the file, past the import statements
func (g *Gen) GenerateBody() (err error) {

	// if the type is not exported we will not capitialize some generated functions either
	firstF, firstA := "F", "A"
	if !unicode.IsUpper(rune(g.opts.typeName[0])) {
		firstF, firstA = "f", "a"
	}

	// define the Option type
	fmt.Fprintf(&g.out, "type %s func(*%s)\n", g.opts.optTypeName, g.opts.typeName)

	// generate the function that generation an option func that completely sets the underlying option
	fmt.Fprintf(&g.out, `func %srom%s(v *%s) %s {
		return func(c *%s) { *c = *v }
	}`+"\n", firstF, strings.Title(g.opts.typeName), g.opts.typeName, g.opts.optTypeName, g.opts.typeName)

	// generate the function that parses the environment
	fmt.Fprintf(&g.out, `func %sParseEnv(eo env.Options) (opts []%s, err error) {
		var o %s
		opts = append(opts, %srom%s(&o))
		return opts, env.Parse(&o, eo)
	}`+"\n", g.opts.optSuffix, g.opts.optTypeName, g.opts.typeName, firstF, strings.Title(g.opts.typeName))

	// generate the function apply function
	fmt.Fprintf(&g.out, `func %spply%ss(opts ...%s) (res %s) {
		env.Parse(&res, env.Options{Environment: make(map[string]string)})
		for _, o := range opts {
			o(&res)
		}
		return
	}`+"\n", firstA, g.opts.optTypeName, g.opts.optTypeName, g.opts.typeName)

	// generate actual options for each field
	for _, file := range g.pkg.Syntax {
		w := Walk{af: file, tf: g.pkg.Fset.File(file.Pos()), Gen: g}
		ast.Inspect(file, w.GenerateOptionFuncs)
	}

	return
}

// Format the generated source code and return it
func (g *Gen) Format() (src []byte, err error) {
	src, err = format.Source(g.Output().Bytes())
	if err != nil {
		return g.out.Bytes(), fmt.Errorf("failed to format source: %w", err)
	}

	return
}

// Output returns the generated output
func (g *Gen) Output() *bytes.Buffer {
	return &g.out
}

// Walk represents a single walk through the ast
type Walk struct {
	err error
	af  *ast.File
	tf  *token.File
	*Gen
}

// GenerateOptionFuncs will generate the functional options for struct tags
func (w *Walk) GenerateOptionFuncs(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		return true
	}

	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.TypeSpec)
		if !ok || spec.Name.String() != w.opts.typeName {
			continue // not the right type
		}

		st, ok := spec.Type.(*ast.StructType)
		if !ok {
			continue // not a struct
		}

		if _, err := w.generateOptionFuncs(st); err != nil {
			w.err = fmt.Errorf("failed to generate option funcs: %w", err)
			return false
		}
	}

	return true
}

func (w Walk) generateOptionFuncs(st *ast.StructType) (ok bool, err error) {
	for _, field := range st.Fields.List {
		tag, _ := strconv.Unquote(field.Tag.Value)
		if _, ok := reflect.StructTag(tag).Lookup(*tagName); !ok {
			continue // skip fields without env tag
		}

		for _, fname := range field.Names {
			name := fname.Name
			if !unicode.IsUpper(rune(name[0])) {
				continue // skip unexported
			}

			d, err := os.ReadFile(w.tf.Name())
			if err != nil {
				return false, fmt.Errorf("failed to read source file for type expression: %w", err)
			}

			start, end := w.tf.Offset(field.Type.Pos()), w.tf.Offset(field.Type.End())
			typ := string(d[start:end])

			w.resCount++
			fmt.Fprintf(&w.out, `func %s%s%s(v %s) %s {return func(o *%s){o.%s=v}}`+"\n",
				w.opts.prefix, name, w.opts.optSuffix, typ, w.opts.optTypeName, w.opts.typeName, name)
		}
	}

	return true, nil
}

// GenerateImports will walk the node tree and copy over all imports from the
// source package over to the generated file.
func (w Walk) GenerateImports(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.IMPORT {
		return true
	}

	for _, spec := range decl.Specs {
		imprt, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		if imprt.Path == nil || imprt.Path.Value == `""` {
			continue
		}

		if imprt.Name != nil {
			fmt.Fprintf(&w.out, `import %s %s`+"\n", imprt.Name.String(), imprt.Path.Value)
		} else {
			fmt.Fprintf(&w.out, `import %s`+"\n", imprt.Path.Value)
		}
	}

	return true
}
