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
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
)

type uid = uint64
type id = uint32

var supportedAnnotations = map[string]bool{
	"-":         true,
	"converter": true,
	"date":      true,
	"id":        true,
	"index":     true,
	"inline":    true,
	"lazy":      true,
	"link":      true,
	"name":      true,
	"type":      true,
	"uid":       true,
	"unique":    true,
}

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
	BaseName    string // name in the containing struct (might be embedded)
	Name        string // prefixed name (unique)
	ObName      string // name of the field in DB
	Annotations map[string]*Annotation
	ObType      int
	ObFlags     []int
	GoType      string
	FbType      string
	IsPointer   bool
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
	Name       string
	uidRequest bool
}

type Index struct {
	Identifier
}

type Annotation struct {
	Value string
}

type Field struct {
	Entity             *Entity // parent entity
	Name               string
	Type               string
	IsPointer          bool
	Property           *Property // nil if it's an embedded struct
	Fields             []*Field  // inner fields, nil if it's a property
	SimpleRelation     *Relation
	StandaloneRelation *StandaloneRelation // to-many relation stored as a standalone relation in the model
	IsLazyLoaded       bool                // only standalone (to-many) relations currently support lazy loading
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
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing entity
			entity.uidRequest = true
		} else if uid, err := strconv.ParseUint(entity.Annotations["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s on entity %s", err, entity.Name)
		} else {
			entity.Uid = uid
		}
	}

	if fields, err := entity.addFields(astStructFieldList{strct, binding.source}, entity.Name, ""); err != nil {
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
					property.addObFlag(PropertyFlagId)
				} else {
					// fail in case multiple fields match this condition
					return fmt.Errorf(
						"id field is missing or multiple fields match the automatic detection condition on "+
							"entity %s - annotate a field with `objectbox:\"id\"` tag", entity.Name)
				}
			}
		}

		if entity.IdProperty == nil {
			return fmt.Errorf("id field is missing on entity %s - either annotate a field with `objectbox:\"id\"` "+
				"tag or use an uint64 field named 'Id/id/ID'", entity.Name)
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

	// fmt.Errorf is called in GetRelated()
	if entity.HasLazyLoadedRelations() {
		binding.Imports["fmt"] = "fmt"
	}
	return nil
}

func (entity *Entity) addFields(fields fieldList, fieldPath, prefix string) ([]*Field, error) {
	var propertyError = func(err error, property *Property) error {
		return fmt.Errorf("%s on property %s, entity %s", err, property.Name, fieldPath)
	}

	var fieldsTree []*Field

	for i := 0; i < fields.Length(); i++ {
		f := fields.Field(i)

		var property = &Property{
			entity: entity,
			path:   fieldPath,
		}

		if name, err := f.Name(); err != nil {
			property.Name = strconv.FormatInt(int64(i), 10) // just for the error message
			return nil, propertyError(err, property)
		} else {
			property.Name = name
		}

		// this is used to correctly render embedded-structs initialization template
		var field = &Field{
			Entity:   entity,
			Name:     property.Name,
			Property: property,
		}

		if tag := f.Tag(); tag != "" {
			if err := property.setAnnotations(tag); err != nil {
				return nil, propertyError(err, property)
			}
		}

		// skip fields with `objectbox:"-"` tag
		if property.Annotations["-"] != nil {
			if len(property.Annotations) != 1 || property.Annotations["-"].Value != "" {
				return nil, propertyError(errors.New("to ignore the property, use only `objectbox:\"-\"` as a tag"), property)
			}
			continue
		}

		// if the embedded field is from a different package, check if it's available (starts with an upercase letter)
		if f.Package().Path() != entity.binding.Package.Path() {
			if len(field.Name) == 0 || field.Name[0] < 65 || field.Name[0] > 90 {
				log.Printf("Note - skipping unavailable field '%s' on entity %s", property.Name, fieldPath)
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

			var innerPrefix = prefix
			if property.Annotations["inline"] == nil {
				// if NOT inline, use prefix based on the field name
				if len(innerPrefix) == 0 {
					innerPrefix = field.Name
				} else {
					innerPrefix = innerPrefix + "_" + field.Name
				}
			}

			if innerFields, err := entity.addFields(innerStructFields, fieldPath+"."+property.Name, innerPrefix); err != nil {
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

		if property.Annotations["name"] != nil {
			if len(property.Annotations["name"].Value) == 0 {
				return nil, propertyError(fmt.Errorf("name annotation value must not be empty - it's the field name in DB"), property)
			} else {
				property.ObName = property.Annotations["name"].Value
			}
		} else {
			property.ObName = property.Name
		}

		property.BaseName = property.Name
		if len(prefix) != 0 {
			property.ObName = prefix + "_" + property.ObName
			property.Name = prefix + "_" + property.Name
		}

		// ObjectBox core internally converts to lowercase so we should check it as this as well
		var realObName = strings.ToLower(property.ObName)
		if entity.propertiesByName[realObName] {
			return nil, propertyError(fmt.Errorf(
				"duplicate name (note that property names are case insensitive)"), property)
		} else {
			entity.propertiesByName[realObName] = true
		}

		if err := property.handleUid(); err != nil {
			return nil, propertyError(err, property)
		}

		entity.Properties = append(entity.Properties, property)
	}

	return fieldsTree, nil
}

// processType analyzes field type information and configures it.
// It might result in setting a field.Type (in case it's one of the basic types),
// field.StandaloneRelation (in case of many-to-many relations) or field.SimpleRelation (one-to-many relations).
// It also updates (fixes) the field.Name on embedded fields
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

	// check if it needs a type cast (it is a named type, not an alias)
	var isNamed bool

	// in case it's a pointer, get it's underlying type
	if pointer, isPointer := baseType.(*types.Pointer); isPointer {
		baseType = pointer.Elem().Underlying()
		field.IsPointer = true
		field.Property.IsPointer = true
		isNamed = typesTypeErrorful{Type: baseType}.IsNamed()
	} else {
		isNamed = typ.IsNamed()
	}

	if err := property.setBasicType(baseType.String()); err == nil {
		// if the baseType is one of the basic supported types

		// check if it needs a type cast (it is a named type, not an alias)
		if isNamed {
			property.CastOnRead = baseType.String()
			property.CastOnWrite = path.Base(typ.String()) // sometimes, it may contain a full import path
		}

		return nil, nil
	}

	// try if it's a struct - it can be either embedded or a relation
	if strct, isStruct := baseType.(*types.Struct); isStruct {
		// fill in the field information
		field.fillInfo(f, typ)

		if property.Annotations["link"] != nil {
			// if it's a one-to-many relation
			if err := property.setRelation(typeBaseName(typ.String()), false); err != nil {
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
		if err := property.setRelation(typeBaseName(elementType.String()), true); err != nil {
			return nil, err
		}

		if err := property.handleUid(); err != nil {
			return nil, err
		}

		// add this as a standalone relation to the entity
		if field.Entity.Relations[field.Name] != nil {
			return nil, fmt.Errorf("relation with the name %s already exists", field.Name)
		}

		rel := &StandaloneRelation{}
		rel.Name = field.Name
		rel.Target.Name = property.Annotations["link"].Value
		rel.uidRequest = property.uidRequest
		rel.Uid = property.Uid

		if _, isPointer := elementType.(*types.Pointer); isPointer {
			rel.Target.IsPointer = true
		}

		field.Entity.Relations[field.Name] = rel
		field.StandaloneRelation = rel
		if field.Property.Annotations["lazy"] != nil {
			// relations only
			field.IsLazyLoaded = true
		}

		// fill in the field information
		field.fillInfo(f, typesTypeErrorful{elementType})
		if rel.Target.IsPointer {
			field.Type = "[]*" + field.Type
		} else {
			field.Type = "[]" + field.Type
		}

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

	// strip leading dots (happens sometimes, I think it's for local types from type-checked package)
	field.Type = strings.TrimLeft(field.Type, ".")

	// if the package path is specified (happens for embedded fields), check whether it's current package
	if strings.ContainsRune(strings.Replace(field.Type, "\\", "/", -1), '/') {
		// if the package is the current package, strip the path & name
		var parts = strings.Split(field.Type, ".")

		if len(parts) == 2 && parts[0] == field.Entity.binding.Package.Path() {
			field.Type = parts[len(parts)-1]
		}

	}
	// get just the last component from `packagename.typename` for the field name
	var parts = strings.Split(field.Name, ".")
	field.Name = parts[len(parts)-1]
}

func (entity *Entity) setAnnotations(comments []*ast.Comment) error {
	lines := parseCommentsLines(comments)

	entity.Annotations = make(map[string]*Annotation)

	for _, tags := range lines {
		// only handle comments in the form of:   // `tags`
		if len(tags) > 1 && tags[0] == tags[len(tags)-1] && tags[0] == '`' {
			if err := parseAnnotations(tags, &entity.Annotations); err != nil {
				entity.Annotations = nil
				return err
			}
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

// setRelation sets a relation on the property.
// If the user has previously defined a relation manually, it must match the arguments (relation target)
func (property *Property) setRelation(target string, manyToMany bool) error {
	if property.Annotations == nil {
		property.Annotations = make(map[string]*Annotation)
	}

	if property.Annotations["link"] == nil {
		property.Annotations["link"] = &Annotation{}
	}

	if len(property.Annotations["link"].Value) == 0 {
		// set the relation target to the type of the target entity
		// TODO this doesn't respect `objectbox:"name:entity"` on the entity (but we don't support that at the moment)
		property.Annotations["link"].Value = target
	} else if property.Annotations["link"].Value != target {
		return fmt.Errorf("relation target mismatch, expected %s, got %s", target, property.Annotations["link"].Value)
	}

	if manyToMany {
		// nothing to do here, it's handled as a standalone relation so this "property" is skipped completely

	} else {
		// add this field as an ID field
		if err := property.setBasicType("uint64"); err != nil {
			return err
		}
	}

	return nil
}

func (property *Property) handleUid() error {
	if property.Annotations["uid"] != nil {
		if len(property.Annotations["uid"].Value) == 0 {
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing property
			property.uidRequest = true
		} else if uid, err := strconv.ParseUint(property.Annotations["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s", err)
		} else {
			property.Uid = uid
		}
	}
	return nil
}

func parseAnnotations(tags string, annotations *map[string]*Annotation) error {
	if len(tags) > 1 && tags[0] == tags[len(tags)-1] && (tags[0] == '`' || tags[0] == '"') {
		tags = tags[1 : len(tags)-1]
	}

	if tags == "" {
		return nil
	}

	// if it's a top-level call, i.e. tags is something like `objectbox:"tag1 tag2:value2" irrelevant:"value"`
	var tag = reflect.StructTag(tags)
	if contents, found := tag.Lookup("objectbox"); found {
		tags = contents
	} else if contents, found := tag.Lookup("ObjectBox"); found {
		tags = contents
	} else {
		return nil
	}

	// tags are space separated
	for _, tag := range strings.Split(tags, " ") {
		if len(tag) > 0 {
			var name string
			var value = &Annotation{}

			// if it contains a colon, it's a key:value pair
			if i := strings.IndexRune(tag, ':'); i >= 0 {
				name = tag[0:i]
				value.Value = tag[i+1:]
			} else {
				// otherwise there's no value
				name = tag
			}

			// names are case insensitive
			name = strings.ToLower(name)

			if (*annotations)[name] != nil {
				return fmt.Errorf("duplicate annotation %s", name)
			} else if !supportedAnnotations[name] {
				return fmt.Errorf("unknown annotation %s", name)
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
		property.ObType = PropertyTypeString
		property.FbType = "UOffsetT"
	} else if ts == "int" || ts == "int64" {
		property.ObType = PropertyTypeLong
		property.FbType = "Int64"
	} else if ts == "uint" || ts == "uint64" {
		property.ObType = PropertyTypeLong
		property.FbType = "Uint64"
		property.addObFlag(PropertyFlagUnsigned)
	} else if ts == "int32" || ts == "rune" {
		property.ObType = PropertyTypeInt
		property.FbType = "Int32"
	} else if ts == "uint32" {
		property.ObType = PropertyTypeInt
		property.FbType = "Uint32"
		property.addObFlag(PropertyFlagUnsigned)
	} else if ts == "int16" {
		property.ObType = PropertyTypeShort
		property.FbType = "Int16"
	} else if ts == "uint16" {
		property.ObType = PropertyTypeShort
		property.FbType = "Uint16"
		property.addObFlag(PropertyFlagUnsigned)
	} else if ts == "int8" {
		property.ObType = PropertyTypeByte
		property.FbType = "Int8"
	} else if ts == "uint8" {
		property.ObType = PropertyTypeByte
		property.FbType = "Uint8"
		property.addObFlag(PropertyFlagUnsigned)
	} else if ts == "byte" {
		property.ObType = PropertyTypeByte
		property.FbType = "Byte"
	} else if ts == "[]byte" {
		property.ObType = PropertyTypeByteVector
		property.FbType = "UOffsetT"
	} else if ts == "[]string" {
		property.ObType = PropertyTypeStringVector
		property.FbType = "UOffsetT"
	} else if ts == "float64" {
		property.ObType = PropertyTypeDouble
		property.FbType = "Float64"
	} else if ts == "float32" {
		property.ObType = PropertyTypeFloat
		property.FbType = "Float32"
	} else if ts == "bool" {
		property.ObType = PropertyTypeBool
		property.FbType = "Bool"
	} else {
		return fmt.Errorf("unknown type %s", ts)
	}

	if property.Annotations["date"] != nil {
		if property.ObType != PropertyTypeLong {
			return fmt.Errorf("invalid underlying type (PropertyType %v) for date field; expecting long", property.ObType)
		} else {
			property.ObType = PropertyTypeDate
		}
	}

	if property.Annotations["link"] != nil {
		if property.ObType != PropertyTypeLong {
			return fmt.Errorf("invalid underlying type (PropertyType %v) for relation field; expecting long", property.ObType)
		} else {
			property.ObType = PropertyTypeRelation
		}
		property.Relation = &Relation{}
		property.Relation.Target.Name = property.Annotations["link"].Value
	}

	return nil
}

func (property *Property) addObFlag(flag int) {
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
		property.addObFlag(PropertyFlagId)
	}

	if property.Annotations["index"] != nil {
		switch strings.ToLower(property.Annotations["index"].Value) {
		case "":
			// if the user doesn't define index type use the default based on the data-type
			if property.ObType == PropertyTypeString {
				property.addObFlag(PropertyFlagIndexHash)
			} else {
				property.addObFlag(PropertyFlagIndexed)
			}
		case "value":
			property.addObFlag(PropertyFlagIndexed)
		case "hash":
			property.addObFlag(PropertyFlagIndexHash)
		case "hash64":
			property.addObFlag(PropertyFlagIndexHash64)
		default:
			return fmt.Errorf("unknown index type %s", property.Annotations["index"].Value)
		}

		if err := property.setIndex(); err != nil {
			return err
		}
	}

	if property.Annotations["unique"] != nil {
		property.addObFlag(PropertyFlagUnique)

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

func (property *Property) ObTypeString() string {
	switch property.ObType {
	case PropertyTypeBool:
		return "Bool"
	case PropertyTypeByte:
		return "Byte"
	case PropertyTypeShort:
		return "Short"
	case PropertyTypeChar:
		return "Char"
	case PropertyTypeInt:
		return "Int"
	case PropertyTypeLong:
		return "Long"
	case PropertyTypeFloat:
		return "Float"
	case PropertyTypeDouble:
		return "Double"
	case PropertyTypeString:
		return "String"
	case PropertyTypeDate:
		return "Date"
	case PropertyTypeRelation:
		return "Relation"
	case PropertyTypeByteVector:
		return "ByteVector"
	case PropertyTypeStringVector:
		return "StringVector"
	default:
		panic(fmt.Errorf("unrecognized type %v", property.ObType))
	}
}

func (property *Property) ObFlagsCombined() int {
	var result = 0

	for _, flag := range property.ObFlags {
		result = result | flag
	}

	return result
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
		if field.HasRelations() {
			return true
		}
	}

	return false
}

func (entity *Entity) HasLazyLoadedRelations() bool {
	for _, field := range entity.Fields {
		if field.HasLazyLoadedRelations() {
			return true
		}
	}

	return false
}

func (field *Field) HasRelations() bool {
	if field.StandaloneRelation != nil || field.SimpleRelation != nil {
		return true
	}

	for _, inner := range field.Fields {
		if inner.HasRelations() {
			return true
		}
	}

	return false
}

func (field *Field) HasLazyLoadedRelations() bool {
	if field.StandaloneRelation != nil && field.IsLazyLoaded {
		return true
	}

	for _, inner := range field.Fields {
		if inner.HasLazyLoadedRelations() {
			return true
		}
	}

	return false
}

// called from the template
func (field *Field) IsId() bool {
	return field.Property == field.Entity.IdProperty
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

	parts = append(parts, property.BaseName)
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
