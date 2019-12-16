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
	"path"
	"path/filepath"
	"strings"
)

type file struct {
	ast            *ast.File
	info           *types.Info
	fileset        *token.FileSet
	files          []*ast.File
	dir            string
	pkgName        string
	typeCheckError error
}

func parseFile(sourceFile string) (f *file, err error) {
	f = &file{
		dir:     filepath.Dir(sourceFile),
		fileset: token.NewFileSet(),
	}

	{ // get the main file's package name
		parsed, err := parser.ParseFile(f.fileset, sourceFile, nil, 0)
		if err != nil {
			return nil, err
		}
		f.pkgName = parsed.Name.Name
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

	if pkgs[f.pkgName] == nil {
		return nil, fmt.Errorf("couldn't find package %s in directory %s", f.pkgName, f.dir)
	}

	// create a list of types in the package the original file belongs to and
	for name, file := range pkgs[f.pkgName].Files {
		if name == sourceFile {
			f.ast = file
		}
		f.files = append(f.files, file)
	}

	if f.ast == nil {
		return nil, fmt.Errorf("the source file %s not found among the files processed in the directory", sourceFile)
	}

	return f, nil
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

func (f *file) importedPackage(name string) (*types.Package, error) {
	for _, imp := range f.ast.Imports {
		if imp.Path == nil {
			return nil, fmt.Errorf("encountered an import without a path: %v", *imp)
		}

		var impPath = strings.Trim(imp.Path.Value, "\"'`")

		if imp.Name != nil && name == imp.Name.Name {
			return types.NewPackage(impPath, name), nil
		}
		if name == path.Base(impPath) {
			return types.NewPackage(impPath, name), nil
		}
	}
	return nil, fmt.Errorf("package %s not imported in the source file", name)
}

func (f *file) analyze() {
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
			// NOTE there is importer.ForCompiler() since 1.12 but it breaks our compatibility with 1.11.4
			// NOTE importer.Default() doesn't seem to work for local files - run the generator tests for more details
			Importer: importer.For("source", nil),
		}

		if _, err := conf.Check(f.dir, f.fileset, f.files, f.info); err != nil {
			// The type checker tries to go on even in case of an error to find out as much as it can.
			// Therefore, this may be an error on an unrelated field and we may still be able to get all the info we
			// need. If the type still can't be determined, we well fail bellow, printing this error as well.
			f.typeCheckError = err
		}

		// find all non-receiver functions (i.e. not related to any struct)
		// this can be used to verify converters exist and have correct signatures, however it only shows functions
		// imported in the package, e.g. it won't show `objectbox.StringIdConvertToEntityProperty`
		// TODO finish verification
		//for _, v := range f.info.Defs {
		//	if def, isFn := v.(*types.Func); isFn {
		//		if signature, isSig := def.Type().(*types.Signature); isSig {
		//			if signature.Recv() == nil {
		//				fmt.Println(def.Pkg().Name(), def.Name(), signature)
		//			}
		//		}
		//	}
		//}
	}
}

func (f *file) getType(expr ast.Expr) (types.Type, error) {
	f.analyze()

	t := f.info.TypeOf(expr)
	if t == nil {
		if f.typeCheckError != nil {
			// report the type checker error for more context
			return nil, fmt.Errorf("type-check error in %v, therefore type %s could not be resolved", f.typeCheckError, expr)
		}
		return nil, fmt.Errorf("type %s could not be resolved", expr)
	}
	return t, nil
}

/// funcSignature returns signature of a function. Can be used to verify converters - see unfinished code in analyze()
//func (f *file) funcSignature(name string) (*types.Signature, error) {
//	return nil, nil
//}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(fnAsVisitor(fn), f.ast)
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
