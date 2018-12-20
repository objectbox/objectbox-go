package generator

import (
	"fmt"
	"go/ast"
	"go/types"
)

// these interfaces are used in the binding to iterate over fields coming from multiple sources (AST & type checker)
type fieldList interface {
	Length() int
	Field(int) field
}

type field interface {
	Name() (string, error)
	Tag() string
	Type() typeErrorful
	TypeInternal() types.Type
}

type typeErrorful interface {
	String() string
	UnderlyingOrError() (types.Type, error)
}

//region ast.StructType wrappers

type astStructFieldList struct {
	*ast.StructType
	source *file
}

func (fields astStructFieldList) Length() int {
	return len(fields.Fields.List)
}

func (fields astStructFieldList) Field(i int) field {
	return &astStructField{fields.Fields.List[i], fields.source}
}

type astStructField struct {
	*ast.Field
	source *file
}

func (field astStructField) Name() (string, error) {
	if len(field.Names) == 0 {
		// in case of an unnamed field, use the type name
		var typ = types.ExprString(field.Field.Type)
		if len(typ) >= 1 && typ[0] == '*' {
			typ = typ[1:] // strip the '*' if it's a pointer type
		}
		return typ, nil
	} else if len(field.Names) == 1 {
		return field.Names[0].Name, nil
	} else {
		return "", fmt.Errorf("the field has too many names: %v", len(field.Names))
	}
}

func (field astStructField) Tag() string {
	if field.Field.Tag != nil {
		return field.Field.Tag.Value
	}
	return ""
}

func (field astStructField) Type() typeErrorful {
	return astTypeExpr{Expr: field.Field.Type, source: field.source}
}
func (field astStructField) TypeInternal() types.Type {
	return astTypeExpr{Expr: field.Field.Type, source: field.source}
}

type astTypeExpr struct {
	ast.Expr
	source *file
}

func (expr astTypeExpr) String() string {
	return types.ExprString(expr.Expr)
}

func (expr astTypeExpr) Underlying() types.Type {
	if t, err := expr.source.getUnderlyingType(expr.Expr); err != nil {
		panic(err)
	} else {
		return t
	}
}

func (expr astTypeExpr) UnderlyingOrError() (types.Type, error) {
	if t, err := expr.source.getUnderlyingType(expr.Expr); err != nil {
		return nil, err
	} else {
		return t, nil
	}
}

//endregion

//region types.Struct wrappers

type structFieldList struct {
	*types.Struct
}

func (fields structFieldList) Length() int {
	return fields.Struct.NumFields()
}

func (fields structFieldList) Field(i int) field {
	return structField{fields.Struct.Field(i), fields.Tag(i)}
}

type structField struct {
	*types.Var
	tag string
}

func (field structField) Name() (string, error) {
	return field.Var.Name(), nil
}

func (field structField) Tag() string {
	return field.tag
}

func (field structField) Type() typeErrorful {
	return typesTypeErrorful{field.Var.Type()}
}

func (field structField) TypeInternal() types.Type {
	return field.Var.Type()
}

type typesTypeErrorful struct {
	types.Type
}

func (typ typesTypeErrorful) String() string {
	return typ.Type.String()
}

func (typ typesTypeErrorful) UnderlyingOrError() (types.Type, error) {
	return typ.Type.Underlying(), nil
}

//endregion
