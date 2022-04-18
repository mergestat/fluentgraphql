package main

import (
	"fmt"

	fgql "github.com/mergestat/fluentgraphql"
)

func main() {
	a := fgql.NewQuery().Scalar("hello").Root().String()

	/*
		{
			hello
		}
	*/
	fmt.Println(a)

	b := fgql.NewQuery().Scalar("hello").Scalar("world").Root().String()

	/*
		{
			hello
			world
		}
	*/
	fmt.Println(b)

	c := fgql.NewQuery(
		fgql.WithName("MyQuery"), fgql.WithVariableDefinitions(fgql.NewVariableDefinition("myVar", "string", true, nil))).
		Scalar("hello", fgql.WithAlias("myAlias"), fgql.WithArguments(
			fgql.NewArgument("anArg", fgql.NewVariableValue("myVar")),
		)).
		Scalar("world").
		Root().String()

	/*
		query MyQuery($myVar: string!) {
			myAlias: hello(anArg: $myVar)
			world
		}
	*/
	fmt.Println(c)
}
