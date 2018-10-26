package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type file struct {
	f *ast.File
}

func parseFile(filePath string) (f *file, err error) {
	f = &file{}

	// Create the AST by parsing src.
	fileset := token.NewFileSet() // positions are relative to fileset

	if f.f, err = parser.ParseFile(fileset, filePath, nil, 0); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(fnAsVisitor(fn), f.f)
}

// walker adapts a function to satisfy the ast.Visitor interface.
// The function return whether the walk should proceed into the node's children.
type fnAsVisitor func(ast.Node) bool

func (fn fnAsVisitor) Visit(node ast.Node) ast.Visitor {
	if fn(node) {
		return fn
	}
	return nil
}
