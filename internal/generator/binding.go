package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"
)

type uid = uint64
type id = uint32

type Binding struct {
	Package  string
	Entities []*Entity

	recentEntity *ast.TypeSpec
	err          error
}

type Entity struct {
	Name       string
	Id         id
	Uid        uid
	Properties []*Property
	PropertyId *Property
}

type Property struct {
	Name        string
	Id          id
	Uid         uid
	Annotations map[string]*Annotation
	ObType      string
	ObFlags     []string // in ascending order
	GoType      string
	FbType      string
	FbSlot      int
}

type Annotation struct {
	Value string
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
		// TODO
		Id:  0,
		Uid: 0,
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

			if err = property.loadAnnotations(f.Tag.Value); err != nil {
				return err
			}

			if err = property.loadType(f.Type); err != nil {
				return err
			}

			if err = property.loadObFlags(f); err != nil {
				return err
			}

			// if this is an ID, set it as entity.PropertyId
			if x := sort.SearchStrings(property.ObFlags, "ID"); property.ObFlags != nil && property.ObFlags[x] == "ID" {
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

func (property *Property) loadAnnotations(tags string) error {
	if len(tags) > 1 && tags[0] == tags[len(tags)-1] && (tags[0] == '`' || tags[0] == '"') {
		tags = tags[1 : len(tags)-1]
	}
	if tags == "" {
		return nil
	}

	property.Annotations = make(map[string]*Annotation)

	// tags are space separated
	for _, tag := range strings.Split(tags, " ") {
		if len(tag) > 0 {
			ss := strings.Split(tags, ":")
			if len(ss) == 1 {
				property.Annotations[tag] = &Annotation{}
			} else if len(ss) == 2 {
				property.Annotations[ss[0]] = &Annotation{Value: ss[1]}
			} else {
				return fmt.Errorf("unkown tag format, multiple colons in the value %s", tag)
			}
		}
	}

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

	// TODO Date (through annotations)
	// TODO relation
	// TODO []byte byte vector

	return nil
}

func (property *Property) loadObFlags(f *ast.Field) error {
	if property.Annotations["id"] != nil || strings.ToLower(property.Name) == "id" {
		property.ObFlags = append(property.ObFlags, "ID")
	}

	// we guarantee that flags are in ascending order so that they could be searchable
	if property.ObFlags != nil {
		sort.Strings(property.ObFlags)
	}
	return nil
}

func (property *Property) loadIds(f *ast.Field) error {
	// TODO how to generate IDs?
	// maybe save them to the original file as tag? Or read the already existing generated file as well?

	return nil
}

// calculates flatbuffers vTableOffset
// called from the template
func (property *Property) VTableOffset() int {
	// TODO verify this, derived from the FB generated code & https://google.github.io/flatbuffers/md__internals.html
	return 4 + 2*property.FbSlot
}
