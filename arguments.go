package fluentgraphql

import (
	"github.com/graphql-go/graphql/language/ast"
)

// argument represents an argument to a GraphQL selection
type argument struct {
	astArg *ast.Argument
}

// NewArgument constructs a new argument with a value
func NewArgument(name string, val *Value) *argument {
	return &argument{
		astArg: ast.NewArgument(&ast.Argument{
			Name:  ast.NewName(&ast.Name{Value: name}),
			Value: val.astValue,
		}),
	}
}

// WithArguments is a selection option for specifying arguments
func WithArguments(args ...*argument) selectionOption {
	return func(s *Selection) {
		for _, arg := range args {
			switch n := s.node.(type) {
			case *ast.Field:
				n.Arguments = append(n.Arguments, arg.astArg)
			}
		}
	}
}
