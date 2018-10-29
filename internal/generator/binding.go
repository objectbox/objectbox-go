package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
	"strings"
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
	PropertyId *Property
}

type Property struct {
	Name   string
	ObType string
	GoType string
	FbType string
	FbSlot int
	Flags  []string
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
				return fmt.Errorf("struct %s has a f with an invalid number of names, one expected, got %v",
					entity.Name, len(f.Names))
			}

			property := &Property{
				Name: f.Names[0].Name,

				// TODO what about backward compatibility? We need to keep track of previous properties
				FbSlot: len(entity.Properties),
			}

			if err = property.loadType(f.Type); err != nil {
				return err
			}

			if strings.ToLower(property.Name) == "id" {
				if entity.PropertyId != nil {
					return fmt.Errorf("struct %s has multiple ID properties - %s and %s",
						entity.Name, entity.PropertyId.Name, property.Name)
				}
				entity.PropertyId = property
			}

			entity.Properties = append(entity.Properties, property)
		}
	}

	binding.Entities = append(binding.Entities, entity)
	return nil
}

func (property *Property) loadType(t ast.Expr) error {
	property.GoType = types.ExprString(t)

	// TODO check thoroughly

	ts := property.GoType
	if property.GoType == "string" {
		property.ObType = "String"
		property.FbType = "UOffsetT"
	} else if ts == "int" || ts == "int64" {
		property.ObType = "Long"
		property.FbType = "Int64"
	} else if ts == "uint" || ts == "uint64" {
		property.ObType = "Long"
		property.FbType = "Uint64"
	} else if ts == "int32" || ts == "rune" {
		property.ObType = "Int"
		property.FbType = "Int32"
	} else if ts == "uint32" {
		property.ObType = "Int"
		property.FbType = "Uint32"
	} else if ts == "int8" {
		property.ObType = "Short"
		property.FbType = "Int8"
	} else if ts == "uint8" {
		property.ObType = "Short"
		property.FbType = "Uint8"
	} else if ts == "int16" {
		property.ObType = "Short"
		property.FbType = "Int16"
	} else if ts == "uint16" {
		property.ObType = "Short"
		property.FbType = "Uint16"
	} else if ts == "float32" {
		property.ObType = "Float"
		property.FbType = "Float32"
	} else if ts == "float64" {
		property.ObType = "Double"
		property.FbType = "Float64"
	} else if ts == "byte" {
		property.ObType = "Byte"
		property.FbType = "Byte"
	} else if ts == "bool" {
		property.ObType = "Bool"
		property.FbType = "Bool"
	} else {
		return fmt.Errorf("unknown type %s", ts)
	}

	// TODO Date (through tags)
	// TODO relation
	// TODO []byte byte vector

	return nil
}

// calculates flatbuffers vTableOffset
// called from the template
func (property *Property) VTableOffset() string {
	// TODO verify this, derived from the FB generated code & https://google.github.io/flatbuffers/md__internals.html
	return strconv.FormatInt(int64(4+2*property.FbSlot), 10)
}
