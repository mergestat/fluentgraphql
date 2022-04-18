package fluentgraphql

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func queryMatchesTree(t *testing.T, query string, root ast.Node) []string {
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

	return deep.Equal(document, expectedDocument)
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
			selection: NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewObjectValue(map[string]*Value{
				"val1": NewStringValue("a"),
				"val2": NewStringValue("b"),
			})))),
		},
		"SingleScalarWithVariableArgument": {
			wantedQuery: `{ hello(arg1: $var1) }`,
			selection:   NewQuery().Scalar("hello", WithArguments(NewArgument("arg1", NewVariableValue("var1")))),
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := testCase.selection.Root()
			if diffs := queryMatchesTree(t, testCase.wantedQuery, root.node); diffs != nil {
				for _, diff := range diffs {
					t.Log(diff)
				}
				t.Fatalf("produced GraphQL query does not match what's wanted: found %d diffs, should be none", len(diffs))
			}
		})
	}
}
