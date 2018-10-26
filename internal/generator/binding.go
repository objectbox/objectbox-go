package generator

import (
	"fmt"
	"go/ast"
)

type Binding struct {
	Package  string
	Entities []*Entity

	recentEntity *ast.TypeSpec
	err          error
}

type Entity struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name string
}

func newBinding() (*Binding, error) {
	return &Binding{}, nil
}

func (binding *Binding) loadAstFile(f *file) (err error) {
	binding.Package = f.f.Name.Name // this is actually package name, not file name

	// process all structs
	f.walk(func(node ast.Node) bool {
		return binding.entityLoader(node)
	})

	if binding.err != nil {
		return binding.err
	}

	return nil
}

// this function only processes structs and cuts-off on types that can't contain a struct
func (binding *Binding) entityLoader(node ast.Node) bool {
	if binding.err != nil {
		return false
	}

	switch v := node.(type) {
	case *ast.TypeSpec:
		// this might be the name of the next struct
		binding.recentEntity = v
		return true
	case *ast.StructType:
		if binding.recentEntity == nil {
			// NOTE this should probably not happen
			binding.err = fmt.Errorf("encountered a struct without a name")
			return false
		} else {
			binding.err = binding.loadAstStruct(node)
			// reset after it has been "consumed"
			binding.recentEntity = nil
		}
		return true

	case *ast.GenDecl:
		return true
	case *ast.File:
		return true
	}

	return false
}

func (binding *Binding) loadAstStruct(node ast.Node) error {
	entity := &Entity{
		Name: binding.recentEntity.Name.Name,
	}

	switch t := node.(type) {
	case *ast.StructType:
		for _, f := range t.Fields.List {
			if len(f.Names) != 1 {
				return fmt.Errorf("Struct %s has a f with an invalid number of names, one expected, got %v",
					entity.Name, len(f.Names))
			}

			field := &Field{
				Name: f.Names[0].Name,
			}

			entity.Fields = append(entity.Fields, field)
		}
	}

	binding.Entities = append(binding.Entities, entity)
	return nil
}
