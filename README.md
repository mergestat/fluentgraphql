# fluentgraphql

This package wraps the [graphql-go/graphql](https://github.com/graphql-go/graphql) implementation to provide a "fluent" pattern for constructing GraphQL queries in Go.
This can be valuable in situations where *dynamic* queries are desired: when the fields of a query are not known until runtime.
For most other use cases, plain query strings or a helper library such as [this](https://github.com/shurcooL/graphql) should be sufficient.
