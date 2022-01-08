package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeName       = flag.String("type", "", "name of the type for which we'll generate options")
	optionTypeName = flag.String("optionType", "Option", "name of the option type that will be generated")
	prefix         = flag.String("prefix", "With", "prefix in front of each generated option function")
	output         = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	tagName        = flag.String("tag", "env", "the name of the env tag that should be scanned for, defaults to 'env'")
)

func main() {
	flag.Parse()
	if *typeName == "" {
		flag.Usage()
		os.Exit(2)
	}

	dst := *output
	if dst == "" {
		baseName := fmt.Sprintf("%s_opts.go", *typeName)
		dst = filepath.Join(os.Getenv("GOPATH"), strings.ToLower(baseName))
	}

	log.SetPrefix("envopts: ")
	log.SetFlags(0)
	_, err := Generate(*typeName, *optionTypeName, *prefix, dst, ".")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
