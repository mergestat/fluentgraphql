package main

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	fgql "github.com/mergestat/fluentgraphql"
)

func main() {
	s := fgql.NewQuery(fgql.WithName("MyQuery"), fgql.WithVariableDefinitions(
		fgql.NewVariableDefinition("var1", "string", false, fgql.NewBooleanValue(true)),
	)).
		Scalar("hello", fgql.WithAlias("alias")).
		Scalar("world").
		Selection("object", fgql.WithArguments(
			fgql.NewArgument("int", fgql.NewIntValue(123)),
			fgql.NewArgument("float", fgql.NewFloatValue(123.102434)),
			fgql.NewArgument("string", fgql.NewStringValue("string-arg")),
			fgql.NewArgument("boolean", fgql.NewBooleanValue(false)),
			fgql.NewArgument("enum", fgql.NewEnumValue("SOME_ENUM")),
			fgql.NewArgument("list", fgql.NewListValue(
				fgql.NewStringValue("str1"),
				fgql.NewStringValue("str2"),
				fgql.NewStringValue("str3"),
			)),
			fgql.NewArgument("object", fgql.NewObjectValue(
				fgql.NewObjectValueField("field1", fgql.NewStringValue("str1")),
				fgql.NewObjectValueField("field1", fgql.NewIntValue(123)),
			)),
			fgql.NewArgument("var", fgql.NewVariableValue("someVar")),
		)).
		Scalar("patrick").Parent().
		Selection("world").Scalar("hello").InlineFragment("User").Scalar("fragmentField").Selection("patrick").Scalar("devivo").
		Root().Fragment("comparisonFields", "Character").Scalar("name").Selection("friendsConnection", fgql.WithArguments(
		fgql.NewArgument("first", fgql.NewVariableValue("first")),
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

	// printer.Print(a)
	fmt.Println(printer.Print(a))

}
