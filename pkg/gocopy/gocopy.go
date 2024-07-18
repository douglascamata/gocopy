package gocopy

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strings"
)

type ErrFunctionNotFound struct {
	functionName string
}

func (e ErrFunctionNotFound) Error() string {
	return fmt.Sprintf("function '%s' not found", e.functionName)
}

type ErrTypeNotFound struct {
	typeName string
}

func (e ErrTypeNotFound) Error() string {
	return fmt.Sprintf("type '%s' not found", e.typeName)
}

type ErrInvalidGoSource struct {
	parseError string
}

func (e ErrInvalidGoSource) Error() string {
	return fmt.Sprintf("invalid Go source: %s", e.parseError)
}

func CopyFunction(fetcher sourceFetcher, fnName string) (string, error) {
	fset, node, err := fetchParseGoSource(fetcher)
	if err != nil {
		return "", fmt.Errorf("error copying function: %w", err)
	}

	buffer := strings.Builder{}
	var found bool
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name != fnName {
				return true
			}
			_ = format.Node(&buffer, fset, x)
			found = true
			return false
		}
		return true
	})
	if !found {
		return "", ErrFunctionNotFound{functionName: fnName}
	}
	return buffer.String(), nil
}

func CopyType(fetcher sourceFetcher, typeName string, includeMethods bool) (string, error) {
	fset, node, err := fetchParseGoSource(fetcher)
	if err != nil {
		return "", fmt.Errorf("error copying type: %w", err)
	}

	buffer := strings.Builder{}
	var found bool
	ast.Inspect(node, func(n ast.Node) bool {
		switch parsedNode := n.(type) {
		case *ast.GenDecl:
			if parsedNode.Tok != token.TYPE {
				return true
			}
			isSingleTypeSpec := len(parsedNode.Specs) == 1
			for _, spec := range parsedNode.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if typeSpec.Name.Name != typeName {
					continue
				}
				found = true
				if isSingleTypeSpec && parsedNode.Doc != nil {
					for _, line := range parsedNode.Doc.List {
						_, _ = io.WriteString(&buffer, line.Text)
					}
					_, _ = io.WriteString(&buffer, "\n")
					_, _ = io.WriteString(&buffer, "type ")
				}
				_ = format.Node(&buffer, fset, typeSpec)
				_, _ = io.WriteString(&buffer, "\n")
				if includeMethods {
					printMethods(&buffer, fset, node, typeSpec.Name.Name)
				}

				return false
			}
		}
		return true
	})
	if !found {
		return "", ErrTypeNotFound{typeName: typeName}
	}
	return buffer.String(), nil
}

func printMethods(writer io.Writer, fset *token.FileSet, node *ast.File, typeName string) {
	for _, f := range node.Decls {
		if funcDecl, ok := f.(*ast.FuncDecl); ok {
			// Check if it's a method (has a receiver) and matches our type
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				// Extracting the type of the receiver
				var receiverType string
				switch typeExpr := funcDecl.Recv.List[0].Type.(type) {
				case *ast.StarExpr: // For pointer receivers e.g., *MyType
					receiverType = typeExpr.X.(*ast.Ident).Name
				case *ast.Ident: // For value receivers e.g., MyType
					receiverType = typeExpr.Name
				}

				if receiverType == typeName {
					// Printing the full method signature
					_ = format.Node(writer, fset, funcDecl)
					_, _ = writer.Write([]byte("\n"))
				}
			}
		}
	}
}

func fetchParseGoSource(fetcher sourceFetcher) (*token.FileSet, *ast.File, error) {
	source, err := fetcher.Fetch()
	if err != nil {
		return nil, nil, err
	}
	return parseGoSource(source)
}

func parseGoSource(source []byte) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", source, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, nil, ErrInvalidGoSource{parseError: err.Error()}
	}
	return fset, node, nil
}
