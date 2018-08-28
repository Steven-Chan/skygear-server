package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	input := os.Args[1]
	if input == "--" {
		input = os.Args[2]
	}
	inputPath := path.Join(dir, input)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Could not parse file: %s", err)
	}

	if f.Name == nil {
		panic("Could not determine package name of " + inputPath)
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Doc == nil {
			continue
		}

		for _, doc := range genDecl.Doc.List {
			if strings.Contains(doc.Text, "@getter") {
				types := strings.Split(strings.Split(doc.Text, "@getter")[1], ",")
				fmt.Printf("get: %+v\n", types)
				typeSepc := genDecl.Specs[0].(*ast.TypeSpec)
				fmt.Printf("type: %+v\n", typeSepc.Name)
			}

			if strings.Contains(doc.Text, "@setter") {
				types := strings.Split(strings.Split(doc.Text, "@setter")[1], ",")
				fmt.Printf("set: %+v\n", types)
			}
		}
	}
}
