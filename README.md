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
A query or mutation is started like so:

```golang
    q := fgql.NewQuery() // a new query builder
    m := fgql.NewMutation() // a new mutation builder
```

A query can be constructed with calls to builder methods, such as in the following example.
See [this file](https://github.com/mergestat/fluentgraphql/blob/main/examples/starwars/main.go) for more thorough examples.

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
Note the call to `.Root().String()`.
`Root()` traverses the builder tree back to the root, so that when `String()` is called, the *entire* query is printed as a string.

### Batching Requests
A use case where a fluent interface is valuable is when dynamically generating a "batch" of queries to make to a GraphQL API.
For instance, in the [`github-batch-request` example](https://github.com/mergestat/fluentgraphql/blob/main/examples/github-batch-request/main.go), we can build a query that retrieves the `stargazerCount` field of multiple, arbitrary repositories at once.
This allows us to batch multiple lookups into a single HTTP request, avoiding multiple round-trip requests.

```golang
q := fgql.NewQuery()

// iterate over the list of repos and add a selection to the query for each one
for i, repo := range repoList {
    split := strings.Split(repo, "/")
    owner := split[0]
    name := split[1]

    q.Selection("repository", fgql.WithAlias(fmt.Sprintf("repo_%d", i)), fgql.WithArguments(
        fgql.NewArgument("owner", fgql.NewStringValue(owner)),
        fgql.NewArgument("name", fgql.NewStringValue(name)),
    )).
        Selection("owner").Scalar("login").Parent().
        Scalar("name").Scalar("stargazerCount")
}
```

Produces a query that looks something like:

```graphql
{
  repo_0: repository(owner: "marko-js", name: "marko") {
    owner {
      login
    }
    name
    stargazerCount
  }
  repo_1: repository(owner: "mithriljs", name: "mithril.js") {
    owner {
      login
    }
    name
    stargazerCount
  }
  repo_2: repository(owner: "angular", name: "angular") {
    owner {
      login
    }
    name
    stargazerCount
  }
  ...
}
```
