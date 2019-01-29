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
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
)

type uid = uint64
type id = uint32

type Binding struct {
	Package  *types.Package
	Entities []*Entity
	Imports  map[string]string

	err    error
	source *file
}

type Entity struct {
	Identifier
	Name           string
	Fields         []*Field // the tree of struct fields (necessary for embedded structs)
	Properties     []*Property
	IdProperty     *Property
	LastPropertyId modelinfo.IdUid
	Relations      map[string]*StandaloneRelation
	Annotations    map[string]*Annotation

	binding          *Binding // parent
	uidRequest       bool
	propertiesByName map[string]bool
}

type Property struct {
	Identifier
	Name        string
	ObName      string
	Annotations map[string]*Annotation
	ObType      string
	ObFlags     []string
	GoType      string
	FbType      string
	Relation    *Relation
	Index       *Index
	Converter   *string

	// type casts for named types
	CastOnRead  string
	CastOnWrite string

	entity     *Entity
	uidRequest bool
	path       string // relative addressing path for embedded structs
}

type Relation struct {
	Target struct {
		Name string
	}
}

type StandaloneRelation struct {
	Identifier
	Target struct {
		Identifier
		Name      string
		IsPointer bool
	}
	Name string
}

type Index struct {
	Identifier
}

type Annotation struct {
	Value string
}

type Field struct {
	entity             *Entity // parent entity
	Name               string
	Type               string
	IsPointer          bool
	Property           *Property // nil if it's an embedded struct
	Fields             []*Field  // inner fields, nil if it's a property
	SimpleRelation     *Relation
	StandaloneRelation *StandaloneRelation
}

type Identifier struct {
	Id  id
	Uid uid
}

func newBinding() (*Binding, error) {
	return &Binding{}, nil
}

func (binding *Binding) createFromAst(f *file) (err error) {
	binding.source = f
	binding.Package = types.NewPackage(f.dir, f.f.Name.Name)
	binding.Imports = make(map[string]string)

	// this will hold the pointer to the latest GenDecl encountered (parent of the current struct)
	var prevDecl *ast.GenDecl

	// traverse the AST to process all structs
	f.walk(func(node ast.Node) bool {
		return binding.entityLoader(node, &prevDecl)
	})

	if binding.err != nil {
		return binding.err
	}

	return nil
}

// this function only processes structs and cuts-off on types that can't contain a struct
func (binding *Binding) entityLoader(node ast.Node, prevDecl **ast.GenDecl) bool {
	if binding.err != nil {
		return false
	}

	switch v := node.(type) {
	case *ast.TypeSpec:
		if strct, isStruct := v.Type.(*ast.StructType); isStruct {
			var name = v.Name.Name

			if name == "" {
				// NOTE this should probably not happen
				binding.err = fmt.Errorf("encountered a struct without a name")
				return false
			}

			var comments []*ast.Comment

			if v.Doc != nil && v.Doc.List != nil {
				// this will be defined in case the struct is inside a block of multiple types - `type (...)`
				comments = v.Doc.List

			} else if prevDecl != nil && *prevDecl != nil && (**prevDecl).Doc != nil && (**prevDecl).Doc.List != nil {
				// otherwise (`type A struct {`), use the docs from the parent GenDecl
				comments = (**prevDecl).Doc.List
			}

			binding.err = binding.createEntityFromAst(strct, name, comments)

			// no need to go any deeper in the AST
			return false
		}

		return true

	case *ast.GenDecl:
		// store the "parent" declaration - we need it to get the comments
		*prevDecl = v
		return true
	case *ast.File:
		return true
	}

	return false
}

func (binding *Binding) createEntityFromAst(strct *ast.StructType, name string, comments []*ast.Comment) error {
	entity := &Entity{
		binding:          binding,
		Name:             name,
		propertiesByName: make(map[string]bool),
		Relations:        make(map[string]*StandaloneRelation),
	}

	if comments != nil {
		if err := entity.setAnnotations(comments); err != nil {
			return fmt.Errorf("%s on entity %s", err, entity.Name)
		}
	}

	if entity.Annotations["uid"] != nil {
		if len(entity.Annotations["uid"].Value) == 0 {
			// in case the user doesn't provide `uid` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing entity
			entity.uidRequest = true
		} else if uid, err := strconv.ParseUint(entity.Annotations["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s on entity %s", err, entity.Name)
		} else {
			entity.Uid = uid
		}
	}

	if fields, err := entity.addFields(astStructFieldList{strct, binding.source}, entity.Name); err != nil {
		return err
	} else {
		entity.Fields = fields
	}

	if len(entity.Properties) == 0 {
		return fmt.Errorf("there are no properties in the entity %s", entity.Name)
	}

	if entity.IdProperty == nil {
		// try to find an ID property automatically based on it's name and type
		for _, property := range entity.Properties {
			if strings.ToLower(property.Name) == "id" &&
				(strings.ToLower(property.GoType) == "uint64" || strings.ToLower(property.GoType) == "string") {
				if entity.IdProperty == nil {
					entity.IdProperty = property
					property.addObFlag("ID")
				} else {
					// fail in case multiple fields match this condition
					return fmt.Errorf(
						"id field is missing on entity %s - annotate a field with `id` tag", entity.Name)
				}
			}
		}

		if entity.IdProperty == nil {
			return fmt.Errorf("id field is missing on entity %s - either annotate a field with `id` tag "+
				"or use an uint64 field named 'Id/id/ID'", entity.Name)
		}
	}

	// special handling for string IDs = they are transformed to uint64 in the binding
	if entity.IdProperty.GoType == "string" {
		if err := entity.IdProperty.setBasicType("uint64"); err != nil {
			return fmt.Errorf("%s on property %s, entity %s", err, entity.IdProperty.Name, entity.Name)
		}

		if entity.IdProperty.Annotations["converter"] == nil {
			var converter = "objectbox.StringIdConvert"
			entity.IdProperty.Converter = &converter
		}
	}

	binding.Entities = append(binding.Entities, entity)
	return nil
}

func (entity *Entity) addFields(fields fieldList, path string) ([]*Field, error) {
	var propertyError = func(err error, property *Property) error {
		return fmt.Errorf("%s on property %s, entity %s", err, property.Name, path)
	}

	var fieldsTree []*Field

	for i := 0; i < fields.Length(); i++ {
		f := fields.Field(i)

		var property = &Property{
			entity: entity,
			path:   path,
		}

		if name, err := f.Name(); err != nil {
			property.Name = strconv.FormatInt(int64(i), 10) // just for the error message
			return nil, propertyError(err, property)
		} else {
			property.Name = name
		}

		// this is used to correctly render embedded-structs initialization template
		var field = &Field{
			entity:   entity,
			Name:     property.Name,
			Property: property,
		}

		if tag := f.Tag(); tag != "" {
			if err := property.setAnnotations(tag); err != nil {
				return nil, propertyError(err, property)
			}
		}

		// transient properties are not stored, thus no need to use it in the binding
		if property.Annotations["transient"] != nil {
			continue
		}

		// if the embedded field is from a different package, check if it's available (starts with an upercase letter)
		if f.Package().Path() != entity.binding.Package.Path() {
			if len(field.Name) == 0 || field.Name[0] < 65 || field.Name[0] > 90 {
				log.Printf("Note - skipping unavailable field '%s' on entity %s", property.Name, path)
				continue
			}

			// import the package. Note that package aliases are not supported yet
			entity.binding.Imports[f.Package().Path()] = f.Package().Path()
		}

		fieldsTree = append(fieldsTree, field)

		if property.Annotations["type"] != nil {
			if err := property.setBasicType(property.Annotations["type"].Value); err != nil {
				return nil, propertyError(err, property)
			}
		} else if innerStructFields, err := field.processType(f); err != nil {
			return nil, propertyError(err, property)
		} else if innerStructFields != nil {
			// if it was recognized as a struct that should be embedded, add all the fields
			if innerFields, err := entity.addFields(innerStructFields, path+"."+property.Name); err != nil {
				return nil, err
			} else {
				// apply some struct-related settings to the field
				field.Property = nil
				field.Fields = innerFields
			}

			// this struct itself is not added, just the inner properties
			// so skip the the following steps of adding the property
			continue
		}

		if err := property.setObFlags(); err != nil {
			return nil, propertyError(err, property)
		}

		if property.Annotations["converter"] != nil {
			if property.Annotations["type"] == nil {
				// TODO this could probably be derived from the type-checker results - see getUnderlyingType()
				return nil, propertyError(errors.New("type annotation has to be specified when using converters"), property)
			}
			property.Converter = &property.Annotations["converter"].Value
		}

		// if this is an ID, set it as entity.IdProperty
		if property.Annotations["id"] != nil {
			if entity.IdProperty != nil {
				return nil, fmt.Errorf("struct %s has multiple ID properties - %s and %s",
					entity.Name, entity.IdProperty.Name, property.Name)
			}
			entity.IdProperty = property
		}

		if property.Annotations["nameindb"] != nil {
			if len(property.Annotations["nameindb"].Value) == 0 {
				return nil, propertyError(fmt.Errorf("nameInDb annotation value must not be empty"), property)
			} else {
				property.ObName = property.Annotations["nameindb"].Value
			}
		} else {
			property.ObName = property.Name
		}

		// ObjectBox core internally converts to lowercase so we should check it as this as well
		var realObName = strings.ToLower(property.ObName)
		if entity.propertiesByName[realObName] {
			return nil, propertyError(fmt.Errorf(
				"duplicate name (note that property names are case insensitive)"), property)
		} else {
			entity.propertiesByName[realObName] = true
		}

		if property.Annotations["uid"] != nil {
			if len(property.Annotations["uid"].Value) == 0 {
				// in case the user doesn't provide `uid` value, it's considered in-process of setting up UID
				// this flag is handled by the merge mechanism and prints the UID of the already existing property
				property.uidRequest = true
			} else if uid, err := strconv.ParseUint(property.Annotations["uid"].Value, 10, 64); err != nil {
				return nil, propertyError(fmt.Errorf("can't parse uid - %s", err), property)
			} else {
				property.Uid = uid
			}
		}

		entity.Properties = append(entity.Properties, property)
	}

	return fieldsTree, nil
}

func (field *Field) processType(f field) (fields fieldList, err error) {
	var typ = f.Type()
	var property = field.Property

	if err := property.setBasicType(typ.String()); err == nil {
		// if it's one of the basic supported types
		return nil, nil
	}

	// if not, get the underlying type and try again
	baseType, err := typ.UnderlyingOrError()
	if err != nil {
		return nil, err
	}

	// in case it's a pointer, get it's underlying type
	if pointer, isPointer := baseType.(*types.Pointer); isPointer {
		baseType = pointer.Elem().Underlying()
		field.IsPointer = true
	}

	if err := property.setBasicType(baseType.String()); err == nil {
		// if the baseType is one of the basic supported types

		// check if it needs a type cast (it is a named type, not an alias)
		if typ.IsNamed() {
			property.CastOnRead = baseType.String()
			property.CastOnWrite = typ.String()
		}

		return nil, nil
	}

	// try if it's a struct - it can be either embedded or a relation
	if strct, isStruct := baseType.(*types.Struct); isStruct {
		// fill in the field information
		field.fillInfo(f, typ)

		if property.Annotations["link"] != nil {
			// if it's a one-to-many relation
			if err := property.forceRelation(typeBaseName(typ.String()), false); err != nil {
				return nil, err
			}

			field.SimpleRelation = property.Relation
			return nil, nil

		} else {
			// otherwise inline all fields
			return structFieldList{strct}, nil
		}
	}

	// check if it's a slice of a non-base type
	if slice, isSlice := baseType.(*types.Slice); isSlice {
		var elementType = slice.Elem()

		// it's a many-to-many relation
		if err := property.forceRelation(typeBaseName(elementType.String()), true); err != nil {
			return nil, err
		}

		// add this as a standalone relation to the entity
		// TODO handle rename of the property (relation) using the uidRequest
		if field.entity.Relations[field.Name] != nil {
			return nil, fmt.Errorf("relation with the name %s already exists", field.Name)
		}

		rel := &StandaloneRelation{}
		rel.Name = field.Name
		rel.Target.Name = property.Annotations["link"].Value

		if _, isPointer := elementType.(*types.Pointer); isPointer {
			rel.Target.IsPointer = true
		}

		field.entity.Relations[field.Name] = rel

		// fill in the field information
		field.fillInfo(f, typ)
		field.StandaloneRelation = rel

		// we need to skip adding this field (it's not persisted in DB) so we add an empty list of fields
		return structFieldList{}, nil
	}

	return nil, fmt.Errorf("unknown type %s", typ.String())
}

func (field *Field) fillInfo(f field, typ typeErrorful) {
	if namedType, isNamed := f.TypeInternal().(*types.Named); isNamed {
		field.Type = namedType.Obj().Name()
	} else {
		field.Type = typ.String()
	}
	// strip the '*' if it's a pointer type
	if len(field.Type) > 1 && field.Type[0] == '*' {
		field.Type = field.Type[1:]
	}
	// get just the last component from `packagename.typename` for the field name
	var parts = strings.Split(field.Name, ".")
	field.Name = parts[len(parts)-1]
}

func (entity *Entity) setAnnotations(comments []*ast.Comment) error {
	lines := parseCommentsLines(comments)

	entity.Annotations = make(map[string]*Annotation)

	for _, tags := range lines {
		if err := parseAnnotations(tags, &entity.Annotations); err != nil {
			entity.Annotations = nil
			return err
		}
	}

	if len(entity.Annotations) == 0 {
		entity.Annotations = nil
	}

	return nil
}

func parseCommentsLines(comments []*ast.Comment) []string {
	var lines []string

	for _, comment := range comments {
		text := comment.Text
		text = strings.TrimSpace(text)

		// text is a single/multi line comment
		if strings.HasPrefix(text, "//") {
			text = strings.TrimPrefix(text, "//")
			lines = append(lines, strings.TrimSpace(text))

		} else if strings.HasPrefix(text, "/*") {
			text = strings.TrimPrefix(text, "/*")
			text = strings.TrimPrefix(text, "*")
			text = strings.TrimSuffix(text, "*/")
			text = strings.TrimSuffix(text, "*")
			text = strings.TrimSpace(text)
			for _, line := range strings.Split(text, "\n") {
				lines = append(lines, strings.TrimSpace(line))
			}
		} else {
			// unknown format, ignore
		}
	}

	return lines
}

func (property *Property) setAnnotations(tags string) error {
	property.Annotations = make(map[string]*Annotation)

	if err := parseAnnotations(tags, &property.Annotations); err != nil {
		property.Annotations = nil
		return err
	}

	if len(property.Annotations) == 0 {
		property.Annotations = nil
	}

	return nil
}

func (property *Property) forceRelation(target string, manyToMany bool) error {
	if property.Annotations == nil {
		property.Annotations = make(map[string]*Annotation)
	}

	if property.Annotations["link"] == nil {
		property.Annotations["link"] = &Annotation{}
	}

	if len(property.Annotations["link"].Value) == 0 {
		// set the relation target to the type of the target entity
		// TODO this doesn't respect nameInDb on the entity (but we don't support that at the moment)
		property.Annotations["link"].Value = target
	} else if property.Annotations["link"].Value != target {
		return fmt.Errorf("relation target mismatch, expected %s, got %s", target, property.Annotations["link"].Value)
	}

	if manyToMany {

	} else {
		// add this field as an ID field
		if err := property.setBasicType("uint64"); err != nil {
			return err
		}
	}

	return nil
}

func parseAnnotations(tags string, annotations *map[string]*Annotation) error {
	if len(tags) > 1 && tags[0] == tags[len(tags)-1] && (tags[0] == '`' || tags[0] == '"') {
		tags = tags[1 : len(tags)-1]
	} else {
		return nil
	}

	if tags == "" {
		return nil
	}

	// tags are space separated
	for _, tag := range strings.Split(tags, " ") {
		if len(tag) > 0 {
			var name string
			var value = &Annotation{}

			// if it contains a colon, it's a key:"value" pair
			if i := strings.IndexRune(tag, ':'); i >= 0 {
				name = tag[0:i]
				tag = tag[i+1:]

				if len(tag) > 1 && tag[0] == tag[len(tag)-1] && tag[0] == '"' {
					value.Value = strings.TrimSpace(tag[1 : len(tag)-1])
				} else {
					return fmt.Errorf("invalid annotation value %s for %s, expecting `name:\"value\"` format",
						tag, name)
				}
			} else {
				// otherwise there's no value
				name = tag
			}

			name = strings.ToLower(name)

			if (*annotations)[name] != nil {
				return fmt.Errorf("duplicate annotation %s", name)
			} else {
				(*annotations)[name] = value
			}
		}
	}

	return nil
}

func (property *Property) setBasicType(baseType string) error {
	property.GoType = baseType

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
	} else if ts == "int16" {
		property.ObType = "Short"
		property.FbType = "Int16"
	} else if ts == "uint16" {
		property.ObType = "Short"
		property.FbType = "Uint16"
	} else if ts == "int8" {
		property.ObType = "Byte"
		property.FbType = "Int8"
	} else if ts == "uint8" {
		property.ObType = "Byte"
		property.FbType = "Uint8"
	} else if ts == "byte" {
		property.ObType = "Byte"
		property.FbType = "Byte"
	} else if ts == "[]byte" {
		property.ObType = "ByteVector"
		property.FbType = "UOffsetT"
	} else if ts == "[]string" {
		property.ObType = "StringVector"
		property.FbType = "UOffsetT"
	} else if ts == "float64" {
		property.ObType = "Double"
		property.FbType = "Float64"
	} else if ts == "float32" {
		property.ObType = "Float"
		property.FbType = "Float32"
	} else if ts == "bool" {
		property.ObType = "Bool"
		property.FbType = "Bool"
	} else {
		return fmt.Errorf("unknown type %s", ts)
	}

	if property.Annotations["date"] != nil {
		if property.ObType != "Long" {
			return fmt.Errorf("invalid underlying type (%s) for date field", property.ObType)
		} else {
			property.ObType = "Date"
		}
	}

	if property.Annotations["link"] != nil {
		if property.ObType != "Long" {
			return fmt.Errorf("invalid underlying type (%s) for relation field", property.ObType)
		} else {
			property.ObType = "Relation"
		}
		property.Relation = &Relation{}
		property.Relation.Target.Name = property.Annotations["link"].Value
	}

	return nil
}

func (property *Property) addObFlag(flag string) {
	property.ObFlags = append(property.ObFlags, flag)
}

func (property *Property) setIndex() error {
	if property.Index != nil {
		return fmt.Errorf("index is already defined")
	} else {
		property.Index = &Index{}
		return nil
	}
}

func (property *Property) setObFlags() error {
	if property.Annotations["id"] != nil {
		property.addObFlag("ID")
	}

	if property.Annotations["index"] != nil {
		switch strings.ToLower(property.Annotations["index"].Value) {
		case "":
			// if the user doesn't define index type use the default based on the data-type
			if property.ObType == "String" {
				property.addObFlag("INDEX_HASH")
			} else {
				property.addObFlag("INDEXED")
			}
		case "value":
			property.addObFlag("INDEXED")
		case "hash":
			property.addObFlag("INDEX_HASH")
		case "hash64":
			property.addObFlag("INDEX_HASH64")
		default:
			return fmt.Errorf("unknown index type %s", property.Annotations["index"].Value)
		}

		if err := property.setIndex(); err != nil {
			return err
		}
	}

	if property.Annotations["unique"] != nil {
		property.addObFlag("UNIQUE")

		if err := property.setIndex(); err != nil {
			return err
		}
	}

	if property.Relation != nil {
		if err := property.setIndex(); err != nil {
			return err
		}
	}

	return nil
}

// called from the template
// avoid GO error "variable declared and not used"
func (entity *Entity) HasNonIdProperty() bool {
	for _, prop := range entity.Properties {
		if prop != entity.IdProperty {
			return true
		}
	}

	return false
}

func (entity *Entity) HasRelations() bool {
	for _, field := range entity.Fields {
		if field.StandaloneRelation != nil {
			return true
		}
		if field.SimpleRelation != nil {
			return true
		}
	}

	return false
}

// called from the template
func (field *Field) IsId() bool {
	return field.Property == field.entity.IdProperty
}

// calculates flatbuffers vTableOffset
// called from the template
func (property *Property) FbvTableOffset() uint16 {
	// derived from the FB generated code & https://google.github.io/flatbuffers/md__internals.html
	var result = 4 + 2*uint32(property.FbSlot())

	if uint32(uint16(result)) != result {
		panic(fmt.Errorf("can't calculate FlatBuffers VTableOffset: property %s ID %d is too large",
			property.Name, property.Id))
	}

	return uint16(result)
}

// calculates flatbuffers slot number
// called from the template
func (property *Property) FbSlot() int {
	return int(property.Id - 1)
}

// returns full path to the property (in embedded struct)
// called from the template
func (property *Property) Path() string {
	var parts = strings.Split(property.path, ".")

	// strip the first component
	parts = parts[1:]

	parts = append(parts, property.Name)
	return strings.Join(parts, ".")
}

func typeBaseName(name string) string {
	// strip the '*' if it's a pointer type
	name = strings.TrimPrefix(name, "*")

	// get just the last component from `packagename.typename` for the field name
	if strings.ContainsRune(name, '.') {
		name = strings.TrimPrefix(path.Ext(name), ".")
	}

	return name
}
