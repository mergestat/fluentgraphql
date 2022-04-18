package main

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	fgl "github.com/mergestat/fluentgraphql"
)

func main() {
	s := fgl.NewQuery(fgl.WithName("MyQuery"), fgl.WithVariableDefinitions(
		fgl.NewVariableDefinition("var1", "string", false, fgl.NewBooleanValue(true)),
	)).
		Scalar("hello", fgl.WithAlias("alias")).
		Scalar("world").
		Selection("object", fgl.WithArguments(
			fgl.NewArgument("int", fgl.NewIntValue(123)),
			fgl.NewArgument("float", fgl.NewFloatValue(123.102434)),
			fgl.NewArgument("string", fgl.NewStringValue("string-arg")),
			fgl.NewArgument("boolean", fgl.NewBooleanValue(false)),
			fgl.NewArgument("enum", fgl.NewEnumValue("SOME_ENUM")),
			fgl.NewArgument("list", fgl.NewListValue(
				fgl.NewStringValue("str1"),
				fgl.NewStringValue("str2"),
				fgl.NewStringValue("str3"),
			)),
			fgl.NewArgument("object", fgl.NewObjectValue(map[string]*fgl.Value{
				"field1": fgl.NewStringValue("str1"),
				"field2": fgl.NewIntValue(123),
			})),
			fgl.NewArgument("var", fgl.NewVariableValue("someVar")),
		)).
		Scalar("patrick").Parent().
		Selection("world").Scalar("hello").InlineFragment("User").Scalar("fragmentField").Selection("patrick").Scalar("devivo").
		Root().Fragment("comparisonFields", "Character").Scalar("name").Selection("friendsConnection", fgl.WithArguments(
		fgl.NewArgument("first", fgl.NewVariableValue("first")),
	)).Scalar("totalCount").Selection("edges").Selection("node").Scalar("name").
		Root().String()

	fmt.Println(s)

	a := ast.NewOperationDefinition(&ast.OperationDefinition{
		Operation: ast.OperationTypeMutation,
		Name:      ast.NewName(&ast.Name{Value: ""}),
		SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
			Selections: []ast.Selection{
				ast.NewField(&ast.Field{
					Name: ast.NewName(&ast.Name{Value: "hello"}),
				}),
				ast.NewField(&ast.Field{
					Name: ast.NewName(&ast.Name{Value: "world"}),
					Arguments: []*ast.Argument{
						ast.NewArgument(&ast.Argument{
							Name:  ast.NewName(&ast.Name{Value: "arg1"}),
							Value: ast.NewIntValue(&ast.IntValue{Value: "123"}),
						}),
					},
				}),
				ast.NewField(&ast.Field{
					Name: ast.NewName(&ast.Name{Value: "author"}),
					SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
						Selections: []ast.Selection{
							ast.NewInlineFragment((&ast.InlineFragment{
								TypeCondition: ast.NewNamed(&ast.Named{Name: ast.NewName(&ast.Name{Value: "User"})}),
								SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
									Selections: []ast.Selection{
										ast.NewField(&ast.Field{
											Name: ast.NewName(&ast.Name{Value: "userField"}),
										}),
									},
								}),
							})),
						},
					}),
				}),
			},
		}),
	})

	printer.Print(a)
	// fmt.Println(printer.Print(a))

}
