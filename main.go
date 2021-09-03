package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

var (
	src    = flag.String("src", "", "Input File")
	dest   = flag.String("dest", "test/", "Output Directory")
	header = `// DO NOT EDIT
// File contents generated automatically by go generate
//

`
)

//go:generate go run main.go -src ./pkg/dummy/main.go
func main() {
	flag.Parse()

	if *src == "" {
		log.Fatal("src is required")
	}

	if *dest == "" {
		log.Fatal("dest is required")
	}

	res, err := parseMethods(*src)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(*dest); os.IsNotExist(err) {
		err = os.Mkdir(*dest, 0755)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(filepath.Join(*dest, strings.Replace(filepath.Base(*src), ".go", "_test.go", -1)))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(header)
	f.Write(res)
}

func parseMethods(src string) ([]byte, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	data := &bytes.Buffer{}

	fmt.Fprintf(data, "package %s_test\n", f.Name.Name)

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			x.Name.Name = fmt.Sprintf("Test%s", strings.Title(x.Name.Name))
			if x.Recv != nil {
				x.Recv = nil
			}
			cg := &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Text:  fmt.Sprintf("// %s is a test", x.Name.Name),
						Slash: x.Pos() - 1,
					},
				},
			}
			x.Doc = cg
			printer.Fprint(data, fset, x)
			fmt.Fprintf(data, "\n\n")
		}
		return true
	})

	if len(data.Bytes()) == 0 {
		return nil, errors.New("function not found")
	}

	return imports.Process("", data.Bytes(), &imports.Options{
		Fragment: true,
		Comments: true,
	})
}
