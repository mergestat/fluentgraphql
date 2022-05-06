package fluentgraphql

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
)

type Selection struct {
	parent *Selection
	node   ast.Node
}

// NewQuery returns a selection builder for a new GraphQL query.
// query { ... }
func NewQuery(options ...operationOption) *Selection {
	s := &Selection{
		parent: nil,
		node: ast.NewOperationDefinition(&ast.OperationDefinition{
			Operation: ast.OperationTypeQuery,
			// VariableDefinitions: make([]*ast.VariableDefinition, 0),
			Directives:   make([]*ast.Directive, 0),
			SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{}),
		}),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// NewMutation returns a selection builder for a new GraphQL mutation.
// mutation { ... }
func NewMutation(options ...operationOption) *Selection {
	s := &Selection{
		parent: nil,
		node: ast.NewOperationDefinition(&ast.OperationDefinition{
			Operation:           ast.OperationTypeMutation,
			VariableDefinitions: make([]*ast.VariableDefinition, 0),
			Directives:          make([]*ast.Directive, 0),
			SelectionSet:        ast.NewSelectionSet(&ast.SelectionSet{}),
		}),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

type operationOption selectionOption

// WithName specifies a name for the operation
func WithName(name string) operationOption {
	return func(s *Selection) {
		switch n := s.node.(type) {
		case *ast.OperationDefinition:
			n.Name = ast.NewName(&ast.Name{
				Value: name,
			})
			n.VariableDefinitions = make([]*ast.VariableDefinition, 0)
		}
	}
}

type variableDefinition struct {
	astVarDef *ast.VariableDefinition
}

// NewVariableDefinition defines a new variable definition
func NewVariableDefinition(name string, varType string, required bool, defaultVal *Value) *variableDefinition {
	varDef := &variableDefinition{
		astVarDef: ast.NewVariableDefinition(&ast.VariableDefinition{
			Variable: NewVariableValue(name).astValue.(*ast.Variable),
		}),
	}

	if required {
		varDef.astVarDef.Type = ast.NewNonNull(&ast.NonNull{
			Type: ast.NewNamed(&ast.Named{
				Name: ast.NewName(&ast.Name{
					Value: varType,
				}),
			}),
		})
	} else {
		varDef.astVarDef.Type = ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{
				Value: varType,
			}),
		})
	}

	if defaultVal != nil {
		varDef.astVarDef.DefaultValue = defaultVal.astValue
	}

	return varDef
}

// WithVariableDefinitions is an operation option for declaring variable definitions
func WithVariableDefinitions(vars ...*variableDefinition) operationOption {
	return func(s *Selection) {
		switch n := s.node.(type) {
		case *ast.OperationDefinition:
			varDefs := make([]*ast.VariableDefinition, 0, len(vars))
			for _, v := range vars {
				varDefs = append(varDefs, v.astVarDef)
			}
			n.VariableDefinitions = varDefs
		}
	}
}

// Scalar adds a scalar field to the current selection
func (s *Selection) Scalar(fieldName string, options ...selectionOption) *Selection {
	newS := &Selection{
		node: ast.NewField(&ast.Field{
			Name:       ast.NewName(&ast.Name{Value: fieldName}),
			Arguments:  make([]*ast.Argument, 0),
			Directives: make([]*ast.Directive, 0),
		}),
	}
	switch n := s.node.(type) {
	case *ast.OperationDefinition:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	case *ast.Field:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	case *ast.InlineFragment:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	case *ast.FragmentDefinition:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	}

	for _, option := range options {
		option(newS)
	}

	return s
}

// Selection adds a subselection to the current selection
func (s *Selection) Selection(fieldName string, options ...selectionOption) *Selection {
	newS := &Selection{
		parent: s,
		node: ast.NewField(&ast.Field{
			Name:         ast.NewName(&ast.Name{Value: fieldName}),
			SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{}),
			Arguments:    make([]*ast.Argument, 0),
			Directives:   make([]*ast.Directive, 0),
		}),
	}
	switch n := s.node.(type) {
	case *ast.OperationDefinition:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	case *ast.Field:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	case *ast.InlineFragment:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	case *ast.FragmentDefinition:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.Field))
	}

	for _, option := range options {
		option(newS)
	}

	return newS
}

// InlineFragment adds an inline fragment to the current selection
func (s *Selection) InlineFragment(typeCondition string) *Selection {
	newFrag := ast.NewInlineFragment((&ast.InlineFragment{
		TypeCondition: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: typeCondition})}),
		SelectionSet:  ast.NewSelectionSet(&ast.SelectionSet{}),
	}))
	newS := &Selection{
		parent: s,
		node:   newFrag,
	}
	switch n := s.node.(type) {
	case *ast.Field:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.InlineFragment))
	}

	return newS
}

// Fragment adds a fragement definition
func (s *Selection) Fragment(name, typeCondition string) *Selection {
	newFrag := ast.NewFragmentDefinition((&ast.FragmentDefinition{
		Name:          ast.NewName(&ast.Name{Value: name}),
		TypeCondition: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: typeCondition})}),
		SelectionSet:  ast.NewSelectionSet(&ast.SelectionSet{}),
	}))
	newS := &Selection{
		parent: s,
		node:   newFrag,
	}
	switch n := s.node.(type) {
	case *ast.Field:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.FragmentDefinition))
	case *ast.OperationDefinition:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.FragmentDefinition))
	}

	return newS
}

// FragmentSpread adds a fragement spread
func (s *Selection) FragmentSpread(name string) *Selection {
	newFrag := ast.NewFragmentSpread((&ast.FragmentSpread{
		Name: ast.NewName(&ast.Name{Value: name}),
	}))
	newS := &Selection{
		parent: s,
		node:   newFrag,
	}
	switch n := s.node.(type) {
	case *ast.Field:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.FragmentSpread))
	case *ast.OperationDefinition:
		n.SelectionSet.Selections = append(n.SelectionSet.Selections, newS.node.(*ast.FragmentSpread))
	}

	return s
}

// Parent returns the parent of this selection. If it's the root, will return nil.
func (s *Selection) Parent() *Selection {
	return s.parent
}

// Root traverses all parents of the current selection until the root
func (s *Selection) Root() *Selection {
	if s.parent == nil {
		return s
	}
	current := s.parent
	for {
		if current.parent == nil {
			return current
		} else {
			current = current.parent
		}
	}
}

// String returns the selection as a GraphQL query string
func (s *Selection) String() string {
	return printer.Print(s.node).(string)
}

// selectionOption enables options for a selection
type selectionOption func(*Selection)

// WithAlias is an option for specifying a selection alias
func WithAlias(alias string) selectionOption {
	return func(s *Selection) {
		switch n := s.node.(type) {
		case *ast.Field:
			n.Alias = ast.NewName(&ast.Name{
				Value: alias,
			})
		}
	}
}
