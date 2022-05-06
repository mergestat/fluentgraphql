package fluentgraphql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func queryMatchesTree(t *testing.T, query string, root ast.Node) string {
	t.Helper()
	opts := parser.ParseOptions{
		NoSource:   true,
		NoLocation: true,
	}
	params := parser.ParseParams{
		Source:  query,
		Options: opts,
	}
	expectedDocument, err := parser.Parse(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	document := ast.NewDocument(&ast.Document{
		Definitions: []ast.Node{root},
	})

	return cmp.Diff(expectedDocument, document)
}

func TestSelections(t *testing.T) {
	for name, testCase := range map[string]struct {
		wanted    string
		selection *Selection
	}{
		"SingleScalar": {
			wanted:    `{ hello }`,
			selection: NewQuery().Scalar("hello"),
		},
		"SingleScalarWithIntArgument": {
			wanted:    `{ hello(arg1: 123) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewIntValue(123)))),
		},
		"SingleScalarWithFloatArgument": {
			wanted:    `{ hello(arg1: 123.12) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewFloatValue(123.12)))),
		},
		"SingleScalarWithStringArgument": {
			wanted:    `{ hello(arg1: "string") }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewStringValue("string")))),
		},
		"SingleScalarWithBooleanArgument": {
			wanted:    `{ hello(arg1: true) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewBooleanValue(true)))),
		},
		"SingleScalarWithEnumArgument": {
			wanted:    `{ hello(arg1: SOME_ENUM) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewEnumValue("SOME_ENUM")))),
		},
		"SingleScalarWithListArgument": {
			wanted:    `{ hello(arg1: ["some-string", "some-other-string"]) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewListValue(NewStringValue("some-string"), NewStringValue("some-other-string"))))),
		},
		"SingleScalarWithObjectArgument": {
			wanted: `{ hello(arg1: { val1: "a", val2: "b"}) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewObjectValue(
				NewObjectValueField("val1", NewStringValue("a")),
				NewObjectValueField("val2", NewStringValue("b")),
			)))),
		},
		"SingleScalarWithVariableArgument": {
			wanted:    `{ hello(arg1: $var1) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewVariableValue("var1")))),
		},
		"SingleScalarWithQueryName": {
			wanted:    `query SomeName { hello }`,
			selection: NewQuery(WithName("SomeName")).Scalar("hello"),
		},
		"SingleScalarWithAlias": {
			wanted:    `{ someAlias: hello }`,
			selection: NewQuery().Scalar("hello", WithAlias("someAlias")),
		},
		"SingleScalarWithVariable": {
			wanted:    `query($a: string) { hello }`,
			selection: NewQuery(WithVariableDefinitions(NewVariableDefinition("a", "string", false, nil))).Scalar("hello"),
		},
		"SingleScalarWithRequiredVariable": {
			wanted:    `query($a: string!) { hello }`,
			selection: NewQuery(WithVariableDefinitions(NewVariableDefinition("a", "string", true, nil))).Scalar("hello"),
		},
		"SingleScalarWithVariableAndName": {
			wanted:    `query NamedQuery($a: string) { hello }`,
			selection: NewQuery(WithName("NamedQuery"), WithVariableDefinitions(NewVariableDefinition("a", "string", false, nil))).Scalar("hello"),
		},
		"SubSelection": {
			wanted:    `{ hello { world } }`,
			selection: NewQuery().Selection("hello").Scalar("world"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := testCase.selection.Root()
			if diff := queryMatchesTree(t, testCase.wanted, root.node); diff != "" {
				t.Log("produced GraphQL query does not match what's wanted", diff)
				t.Fatal()
			}
		})
	}
}

func TestMutations(t *testing.T) {
	for name, testCase := range map[string]struct {
		wanted    string
		selection *Selection
	}{
		"SingleScalar": {
			wanted:    `mutation { hello }`,
			selection: NewMutation().Scalar("hello"),
		},
		"SingleScalarWithIntArgument": {
			wanted:    `mutation { hello(arg1: 123) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewIntValue(123)))),
		},
		"SingleScalarWithFloatArgument": {
			wanted:    `mutation { hello(arg1: 123.12) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewFloatValue(123.12)))),
		},
		"SingleScalarWithStringArgument": {
			wanted:    `mutation { hello(arg1: "string") }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewStringValue("string")))),
		},
		"SingleScalarWithBooleanArgument": {
			wanted:    `mutation { hello(arg1: true) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewBooleanValue(true)))),
		},
		"SingleScalarWithEnumArgument": {
			wanted:    `mutation { hello(arg1: SOME_ENUM) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewEnumValue("SOME_ENUM")))),
		},
		"SingleScalarWithListArgument": {
			wanted:    `mutation { hello(arg1: ["some-string", "some-other-string"]) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewListValue(NewStringValue("some-string"), NewStringValue("some-other-string"))))),
		},
		"SingleScalarWithObjectArgument": {
			wanted: `mutation { hello(arg1: { val1: "a", val2: "b"}) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewObjectValue(
				NewObjectValueField("val1", NewStringValue("a")),
				NewObjectValueField("val2", NewStringValue("b")),
			)))),
		},
		"SingleScalarWithVariableArgument": {
			wanted:    `mutation { hello(arg1: $var1) }`,
			selection: NewMutation().Scalar("hello", WithArguments(NewArgument("arg1", NewVariableValue("var1")))),
		},
		"SingleScalarWithQueryName": {
			wanted:    `mutation SomeName { hello }`,
			selection: NewMutation(WithName("SomeName")).Scalar("hello"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := testCase.selection.Root()
			if diff := queryMatchesTree(t, testCase.wanted, root.node); diff != "" {
				t.Log("produced GraphQL query does not match what's wanted", diff)
				t.Fatal()
			}
		})
	}
}
