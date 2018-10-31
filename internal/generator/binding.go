package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

type uid = uint64
type id = uint32

type Binding struct {
	Package  string
	Entities []*Entity

	currentEntityName string
	err               error
}

type Entity struct {
	Name       string
	Id         id
	Uid        uid
	Properties []*Property
	PropertyId *Property // TODO what if ID is not defined? what about GetId function?
}

type Property struct {
	Name        string
	ObName      string
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
		binding.currentEntityName = v.Name.Name
		return true
	case *ast.StructType:
		if binding.currentEntityName == "" {
			// NOTE this should probably not happen
			binding.err = fmt.Errorf("encountered a struct without a name")
			return false
		} else {
			binding.err = binding.loadAstStruct(node)
			// reset after it has been "consumed"
			binding.currentEntityName = ""
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
		Name: binding.currentEntityName,
		// TODO
		Id:  0,
		Uid: 0,
	}

	var fullError = func(err error, property *Property) error {
		return fmt.Errorf("%s on property %s, entity %s", err, property.Name, entity.Name)
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

			if f.Tag != nil {
				if err = property.loadAnnotations(f.Tag.Value); err != nil {
					return fullError(err, property)
				}
			}

			// transient properties are not stored, thus no need to use it in the binding
			if property.Annotations["transient"] != nil {
				continue
			}

			if err = property.loadType(f.Type); err != nil {
				return fullError(err, property)
			}

			if err = property.loadObFlags(f); err != nil {
				return fullError(err, property)
			}

			// if this is an ID, set it as entity.PropertyId
			if property.Annotations["id"] != nil {
				if entity.PropertyId != nil {
					return fmt.Errorf("struct %s has multiple ID properties - %s and %s",
						entity.Name, entity.PropertyId.Name, property.Name)
				}
				entity.PropertyId = property
			}

			if property.Annotations["nameindb"] != nil {
				if len(property.Annotations["nameindb"].Value) == 0 {
					return fullError(fmt.Errorf("nameInDb annotation value must not be empty"), property)
				} else {
					property.ObName = property.Annotations["nameindb"].Value
				}
			} else {
				property.ObName = property.Name
			}

			entity.Properties = append(entity.Properties, property)
		}
	}

	if entity.PropertyId == nil {
		return fmt.Errorf("field annotated `id` is missing on entity %s", entity.Name)
	}

	binding.Entities = append(binding.Entities, entity)
	return nil
}

// Supported annotations:
// id
// index (value|hash|hash64)
// unique
// nameInDb
// transient
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
			var name string
			var value = &Annotation{}

			// if it contains a colon, it's a key:"value" pair
			if i := strings.IndexRune(tags, ':'); i >= 0 {
				name = tags[0:i]
				tags = tags[i+1:]

				if len(tags) > 1 && tags[0] == tags[len(tags)-1] && tags[0] == '"' {
					value.Value = strings.TrimSpace(tags[1 : len(tags)-1])
				} else {
					return fmt.Errorf("invalid annotation value %s for %s, expecting `name:\"value\"` format", tags, name)
				}
			} else {
				// otherwise there's no value
				name = tags
			}

			name = strings.ToLower(name)

			if property.Annotations[name] != nil {
				return fmt.Errorf("duplicate annotation %s", name)
			} else {
				property.Annotations[name] = value
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

func (property *Property) addObFlag(flag string) {
	property.ObFlags = append(property.ObFlags, flag)
}

func (property *Property) loadObFlags(f *ast.Field) error {
	if property.Annotations["id"] != nil {
		property.addObFlag("ID")
	}

	if property.Annotations["index"] != nil {
		property.addObFlag("INDEXED")
		switch strings.ToLower(property.Annotations["index"].Value) {
		case "":
			// default
		case "value":
			// TODO this doesn't seem to be implemented by the c-api?
			property.addObFlag("INDEX_VALUE")
		case "hash":
			property.addObFlag("INDEX_HASH")
		case "hash64":
			property.addObFlag("INDEX_HASH64")
		default:
			return fmt.Errorf("unknown index type %s", property.Annotations["index"].Value)
		}
	}

	if property.Annotations["unique"] != nil {
		property.addObFlag("UNIQUE")
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
