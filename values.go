package fluentgraphql

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql/language/ast"
)

// Value represents a GraphQL value
type Value struct {
	astValue ast.Value
}

// NewIntValue returns an integer value
func NewIntValue(val int) *Value {
	return &Value{
		astValue: ast.NewIntValue(&ast.IntValue{
			Value: strconv.Itoa(val),
		}),
	}
}

// NewFloatValue returns a float value
func NewFloatValue(val float64) *Value {
	return &Value{
		astValue: ast.NewFloatValue(&ast.FloatValue{
			Value: fmt.Sprintf("%v", val),
		}),
	}
}

// NewStringValue returns a string value
func NewStringValue(val string) *Value {
	return &Value{
		astValue: ast.NewStringValue(&ast.StringValue{
			Value: val,
		}),
	}
}

// NewBooleanValue returns a boolean value
func NewBooleanValue(val bool) *Value {
	return &Value{
		astValue: ast.NewBooleanValue(&ast.BooleanValue{
			Value: val,
		}),
	}
}

// NewEnumValue returns an enum value
func NewEnumValue(val string) *Value {
	return &Value{
		astValue: ast.NewEnumValue(&ast.EnumValue{
			Value: val,
		}),
	}
}

// NewListValue returns a list value
func NewListValue(values ...*Value) *Value {
	vals := make([]ast.Value, 0, len(values))
	for _, v := range values {
		vals = append(vals, v.astValue)
	}

	return &Value{
		astValue: ast.NewListValue(&ast.ListValue{
			Values: vals,
		}),
	}
}

// NewObjectValue returns an object value
func NewObjectValue(values map[string]*Value) *Value {
	fields := make([]*ast.ObjectField, 0, len(values))

	for n, v := range values {
		fields = append(fields, ast.NewObjectField(&ast.ObjectField{
			Name:  ast.NewName(&ast.Name{Value: n}),
			Value: v.astValue,
		}))
	}

	return &Value{
		astValue: ast.NewObjectValue(&ast.ObjectValue{
			Fields: fields,
		}),
	}
}

// NewVariableValue returns a variable value
func NewVariableValue(name string) *Value {
	return &Value{
		astValue: ast.NewVariable(&ast.Variable{
			Name: ast.NewName(&ast.Name{Value: name}),
		}),
	}
}
