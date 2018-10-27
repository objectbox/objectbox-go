package generator

import (
	"fmt"
	"go/ast"
	"go/types"
)

type Binding struct {
	Package  string
	Entities []*Entity

	recentEntity *ast.TypeSpec
	err          error
}

type Entity struct {
	Name       string
	Properties []*Property
}

type Property struct {
	Name string
	Type string
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

func (binding *Binding) loadAstStruct(node ast.Node) (err error) {
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

			property := &Property{
				Name: f.Names[0].Name,
			}

			if property.Type, err = binding.getPropertyType(f.Type); err != nil {
				return err
			}

			entity.Properties = append(entity.Properties, property)
		}
	}

	binding.Entities = append(binding.Entities, entity)
	return nil
}

func (binding *Binding) getPropertyType(t ast.Expr) (typ string, err error) {
	// NOTE potential optimization if we didn't need to convert to string first
	ts := types.ExprString(t)

	if ts == "string" {
		typ = "String"
	} else if ts == "int" || ts == "uint" || ts == "int64" || ts == "uint64" {
		typ = "Long"
	} else if ts == "int32" || ts == "uint32" || ts == "rune" {
		typ = "Int"
	} else if ts == "int8" || ts == "int16" || ts == "uint8" || ts == "uint16" {
		typ = "Short"
	} else if ts == "float32" {
		typ = "Float"
	} else if ts == "float64" {
		typ = "Double"
	} else if ts == "byte" {
		typ = "Byte"
	} else if ts == "bool" {
		typ = "Bool"
	}

	// TODO Date (through tags)
	// TODO relation
	// TODO []byte byte vector

	if len(typ) == 0 {
		err = fmt.Errorf("unknown type %s", t)
	}

	return
}
