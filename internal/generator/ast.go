/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

	if f.f, err = parser.ParseFile(fileset, filePath, nil, parser.ParseComments); err != nil {
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
