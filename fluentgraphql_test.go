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
		wantedQuery string
		selection   *Selection
	}{
		"SingleScalar": {
			wantedQuery: `{ hello }`,
			selection:   NewQuery().Scalar("hello"),
		},
		"SingleScalarWithIntArgument": {
			wantedQuery: `{ hello(arg1: 123) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewIntValue(123)))),
		},
		"SingleScalarWithFloatArgument": {
			wantedQuery: `{ hello(arg1: 123.12) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewFloatValue(123.12)))),
		},
		"SingleScalarWithStringArgument": {
			wantedQuery: `{ hello(arg1: "string") }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewStringValue("string")))),
		},
		"SingleScalarWithBooleanArgument": {
			wantedQuery: `{ hello(arg1: true) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewBooleanValue(true)))),
		},
		"SingleScalarWithEnumArgument": {
			wantedQuery: `{ hello(arg1: SOME_ENUM) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewEnumValue("SOME_ENUM")))),
		},
		"SingleScalarWithListArgument": {
			wantedQuery: `{ hello(arg1: ["some-string", "some-other-string"]) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewListValue(NewStringValue("some-string"), NewStringValue("some-other-string"))))),
		},
		"SingleScalarWithObjectArgument": {
			wantedQuery: `{ hello(arg1: { val1: "a", val2: "b"}) }`,
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewObjectValue(
				NewObjectValueField("val1", NewStringValue("a")),
				NewObjectValueField("val2", NewStringValue("b")),
			)))),
		},
		"SingleScalarWithVariableArgument": {
			wantedQuery: `{ hello(arg1: $var1) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewVariableValue("var1")))),
		},
		"SingleScalarWithQueryName": {
			wantedQuery: `query SomeName { hello }`,
			selection:   NewQuery(WithName("SomeName")).Scalar("hello"),
		},
		"SingleScalarWithAlias": {
			wantedQuery: `{ someAlias: hello }`,
			selection:   NewQuery().Scalar("hello", WithAlias("someAlias")),
		},
		"SingleScalarWithVariable": {
			wantedQuery: `query($a: string) { hello }`,
			selection:   NewQuery(WithVariableDefinitions(NewVariableDefinition("a", "string", false, nil))).Scalar("hello"),
		},
		"SingleScalarWithVariableAndName": {
			wantedQuery: `query NamedQuery($a: string) { hello }`,
			selection:   NewQuery(WithName("NamedQuery"), WithVariableDefinitions(NewVariableDefinition("a", "string", false, nil))).Scalar("hello"),
		},
		"SubSelection": {
			wantedQuery: `{ hello { world } }`,
			selection:   NewQuery().Selection("hello").Scalar("world"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := testCase.selection.Root()
			if diff := queryMatchesTree(t, testCase.wantedQuery, root.node); diff != "" {
				t.Log("produced GraphQL query does not match what's wanted", diff)
				t.Fatal()
			}
		})
	}
}

func TestMutations(t *testing.T) {
	for name, testCase := range map[string]struct {
		wantedQuery string
		selection   *Selection
	}{
		"SingleScalar": {
			wantedQuery: `mutation { hello }`,
			selection:   NewMutation().Scalar("hello"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := testCase.selection.Root()
			if diff := queryMatchesTree(t, testCase.wantedQuery, root.node); diff != "" {
				t.Log("produced GraphQL query does not match what's wanted", diff)
				t.Fatal()
			}
		})
	}
}
