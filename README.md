[![Go Reference](https://pkg.go.dev/badge/github.com/mergestat/fluentgraphql.svg)](https://pkg.go.dev/github.com/mergestat/fluentgraphql)
[![Go Report Card](https://goreportcard.com/badge/github.com/mergestat/fluentgraphql)](https://goreportcard.com/report/github.com/mergestat/fluentgraphql)


# fluentgraphql

This package wraps the [graphql-go/graphql](https://github.com/graphql-go/graphql) implementation to provide a "fluent" pattern for constructing GraphQL queries in Go.
This can be valuable in situations where *dynamic* queries are desired: when the fields of a GraphQL query (or mutation) are not known until runtime.
For most other use cases, plain query strings or a helper library such as [this](https://github.com/shurcooL/graphql) should be sufficient.

```golang
package main

import (
    fgql "github.com/mergestat/fluentgraphql"
)

func main() {
    fgql.NewQuery().Scalar("hello").Root().String() // { hello }
}
```

## Basic Usage

`go get github.com/mergestat/fluentgraphql`

```golang
import (
    fgql "github.com/mergestat/fluentgraphql"
)
```

The package name is `fluentgraphql`, but we alias it here to `fgql` which is more concise.

```golang
    q := fgql.NewQuery() // a new query builder
    m := fgql.NewMutation() // a new mutation builder
```

```golang
/*
    query HeroComparison($first: Int = 3) {
        leftComparison: hero(episode: EMPIRE) {
        ...comparisonFields
        }
        rightComparison: hero(episode: JEDI) {
        ...comparisonFields
        }
    }

    fragment comparisonFields on Character {
        name
        friendsConnection(first: $first) {
        totalCount
        edges {
            node {
            name
            }
        }
        }
    }
*/
q = fgql.NewQuery(
    fgql.WithName("HeroComparison"),
    fgql.WithVariableDefinitions(
        fgql.NewVariableDefinition("first", "Int", false, fgql.NewIntValue(3)),
    ),
).
    Selection("hero",
        fgql.WithAlias("leftComparison"),
        fgql.WithArguments(fgql.NewArgument("episode", fgql.NewEnumValue("EMPIRE"))),
    ).FragmentSpread("comparisonFields").
    Parent().
    Selection("hero",
        fgql.WithAlias("rightComparison"),
        fgql.WithArguments(fgql.NewArgument("episode", fgql.NewEnumValue("JEDI"))),
    ).FragmentSpread("comparisonFields").
    Root().
    Fragment("comparisonFields", "Character").
    Scalar("name").
    Selection("friendsConnection", fgql.WithArguments(fgql.NewArgument("first", fgql.NewVariableValue("first")))).
    Scalar("totalCount").
    Selection("edges").Selection("node").Scalar("name").
    Root().String()
fmt.Println(q)
```
