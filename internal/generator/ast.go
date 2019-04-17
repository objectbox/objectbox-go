/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
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
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	f       *ast.File
	info    *types.Info
	fileset *token.FileSet
	files   []*ast.File
	dir     string
}

func parseFile(sourceFile string) (f *file, err error) {
	f = &file{
		dir:     filepath.Dir(sourceFile),
		fileset: token.NewFileSet(),
	}

	// get the main file's package name
	pkgName, err := getPackageName(sourceFile)
	if err != nil {
		return nil, err
	}

	// parse the whole directory to read & understand the used types
	var filter = func(file os.FileInfo) bool {
		// never skip the sourceFile
		if file.Name() == filepath.Base(sourceFile) {
			return true
		}
		return parserFilter(file)
	}
	var pkgs map[string]*ast.Package
	if pkgs, err = parser.ParseDir(f.fileset, f.dir, filter, parser.ParseComments); err != nil {
		return nil, err
	}

	if pkgs[pkgName] == nil {
		return nil, fmt.Errorf("couldn't find package %s in directory %s", pkgName, f.dir)
	}

	// create a list of types in the package the original file belongs to and
	for name, file := range pkgs[pkgName].Files {
		if name == sourceFile {
			f.f = file
		}
		f.files = append(f.files, file)
	}

	if f.f == nil {
		return nil, fmt.Errorf("the source file %s not found among the files processed in the directory", sourceFile)
	}

	return f, nil
}

func getPackageName(filePath string) (string, error) {
	if f, err := parser.ParseFile(&token.FileSet{}, filePath, nil, 0); err != nil {
		return "", err
	} else {
		return f.Name.Name, nil
	}
}

func parserFilter(file os.FileInfo) bool {
	// skip tests
	if strings.HasSuffix(file.Name(), "_test.go") {
		return false
	}

	// skip files starting with an underscore or a dot (ignored by go build)
	if strings.HasPrefix(file.Name(), "_") || strings.HasPrefix(file.Name(), ".") {
		return false
	}

	return true
}

func (f *file) getType(expr ast.Expr) (types.Type, error) {
	// load file info (resolved types) JiT if necessary
	if f.info == nil {
		// call types.Config.Check() to fill types.Info
		f.info = &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		}

		var conf = types.Config{
			IgnoreFuncBodies:         true,
			DisableUnusedImportCheck: true,
			Importer:                 importer.ForCompiler(f.fileset, "source", nil),
		}
		if _, err := conf.Check(f.dir, f.fileset, f.files, f.info); err != nil {
			return nil, fmt.Errorf("error running type-check: %s", err)
		}
	}

	if t := f.info.TypeOf(expr); t == nil {
		return nil, fmt.Errorf("type %s could not be resolved", expr)
	} else {
		return t, nil
	}
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
