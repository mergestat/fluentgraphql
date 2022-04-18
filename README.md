[![Go Reference](https://pkg.go.dev/badge/github.com/mergestat/fluentgraphql.svg)](https://pkg.go.dev/github.com/mergestat/fluentgraphql)
[![Go Report Card](https://goreportcard.com/badge/github.com/mergestat/fluentgraphql)](https://goreportcard.com/report/github.com/mergestat/fluentgraphql)


# fluentgraphql

This package wraps the [graphql-go/graphql](https://github.com/graphql-go/graphql) implementation to provide a "fluent" pattern for constructing GraphQL queries in Go.
This can be valuable in situations where *dynamic* queries are desired: when the fields of a query are not known until runtime.
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
