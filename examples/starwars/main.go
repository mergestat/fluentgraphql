package main

import (
	"fmt"

	fgql "github.com/mergestat/fluentgraphql"
)

// implements the queries list on this page: https://graphql.org/learn/queries/
func main() {
	var q string
	/*
		{
		  hero {
		    name
		  }
		}
	*/
	q = fgql.NewQuery().Selection("hero").Scalar("name").
		Root().String()
	fmt.Println(q)

	/*
		{
		  hero {
		    name
		    # Queries can have comments!
		    friends {
		      name
		    }
		  }
		}
	*/
	q = fgql.NewQuery().Selection("hero").Scalar("name").Selection("friends").Scalar("name").
		Root().String()
	fmt.Println(q)

	/*
		{
		  human(id: "1000") {
		    name
		    height
		  }
		}
	*/
	q = fgql.NewQuery().Selection("human", fgql.WithArguments(fgql.NewArgument("id", fgql.NewStringValue("1000")))).Scalar("name").Scalar("height").
		Root().String()
	fmt.Println(q)

	/*
		{
		  human(id: "1000") {
		    name
		    height(unit: FOOT)
		  }
		}
	*/
	q = fgql.NewQuery().
		Selection("human",
			fgql.WithArguments(fgql.NewArgument("id", fgql.NewStringValue("1000"))),
		).
		Scalar("name").
		Scalar("height",
			fgql.WithArguments(fgql.NewArgument("unit", fgql.NewEnumValue("FOOT"))),
		).
		Root().String()
	fmt.Println(q)

	/*
		{
		  empireHero: hero(episode: EMPIRE) {
		    name
		  }
		  jediHero: hero(episode: JEDI) {
		    name
		  }
		}
	*/
	q = fgql.NewQuery().
		Selection("hero",
			fgql.WithAlias("empireHero"),
			fgql.WithArguments(fgql.NewArgument("episode", fgql.NewEnumValue("EMPIRE"))),
		).Scalar("name").
		Parent(). // .Parent() moves us "up" one step in the builder tree
		Selection("hero",
			fgql.WithAlias("jediHero"),
			fgql.WithArguments(fgql.NewArgument("episode", fgql.NewEnumValue("JEDI"))),
		).Scalar("name").
		Root().String()
	fmt.Println(q)

	/*
		{
		  leftComparison: hero(episode: EMPIRE) {
		    ...comparisonFields
		  }
		  rightComparison: hero(episode: JEDI) {
		    ...comparisonFields
		  }
		}

		fragment comparisonFields on Character {
		  name
		  appearsIn
		  friends {
		    name
		  }
		}
	*/
	q = fgql.NewQuery().
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
		Scalar("appearsIn").
		Selection("friends").Scalar("name").
		Root().String()
	fmt.Println(q)

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

	/*
		query HeroNameAndFriends {
		  hero {
		    name
		    friends {
		      name
		    }
		  }
		}
	*/
	q = fgql.NewQuery(fgql.WithName("HeroNameAndFriends")).
		Selection("hero").Scalar("name").
		Selection("friends").Scalar("name").
		Root().String()
	fmt.Println(q)

	/*
		query HeroNameAndFriends($episode: Episode) {
		  hero(episode: $episode) {
		    name
		    friends {
		      name
		    }
		  }
		}
	*/
	q = fgql.NewQuery(
		fgql.WithName("HeroNameAndFriends"),
		fgql.WithVariableDefinitions(fgql.NewVariableDefinition("episode", "Episode", false, nil)),
	).
		Selection("hero", fgql.WithArguments(fgql.NewArgument("episode", fgql.NewVariableValue("episode")))).
		Scalar("name").
		Selection("friends").Scalar("name").
		Root().String()
	fmt.Println(q)

	/*
		query HeroNameAndFriends($episode: Episode = JEDI) {
		  hero(episode: $episode) {
		    name
		    friends {
		      name
		    }
		  }
		}
	*/
	q = fgql.NewQuery(
		fgql.WithName("HeroNameAndFriends"),
		fgql.WithVariableDefinitions(fgql.NewVariableDefinition("episode", "Episode", false, fgql.NewEnumValue("JEDI"))),
	).
		Selection("hero", fgql.WithArguments(fgql.NewArgument("episode", fgql.NewVariableValue("episode")))).
		Scalar("name").
		Selection("friends").Scalar("name").
		Root().String()
	fmt.Println(q)

	/*
		query Hero($episode: Episode, $withFriends: Boolean!) {
		  hero(episode: $episode) {
		    name
		    friends @include(if: $withFriends) {
		      name
		    }
		  }
		}
	*/
	// TODO(patrickdevivo) implement when directives are supported

	/*
		mutation CreateReviewForEpisode($ep: Episode!, $review: ReviewInput!) {
		  createReview(episode: $ep, review: $review) {
		    stars
		    commentary
		  }
		}
	*/
	q = fgql.NewMutation(
		fgql.WithName("CreateReviewForEpisode"),
		fgql.WithVariableDefinitions(
			fgql.NewVariableDefinition("ep", "Episode", true, nil),
			fgql.NewVariableDefinition("review", "ReviewInput", true, nil),
		),
	).
		Selection("createReview", fgql.WithArguments(
			fgql.NewArgument("episode", fgql.NewVariableValue("ep")),
			fgql.NewArgument("review", fgql.NewVariableValue("review")),
		)).
		Scalar("stars").Scalar("commentary").
		Root().String()
	fmt.Println(q)

	/*
		query HeroForEpisode($ep: Episode!) {
		  hero(episode: $ep) {
		    name
		    ... on Droid {
		      primaryFunction
		    }
		    ... on Human {
		      height
		    }
		  }
		}
	*/
	q = fgql.NewQuery(
		fgql.WithName("HeroForEpisode"),
		fgql.WithVariableDefinitions(
			fgql.NewVariableDefinition("ep", "Episode", true, nil),
		),
	).
		Selection("hero", fgql.WithArguments(
			fgql.NewArgument("episode", fgql.NewVariableValue("ep")),
		)).
		Scalar("name").
		InlineFragment("Droid").Scalar("primaryFunction").
		Parent().
		InlineFragment("Human").Scalar("height").
		Root().String()
	fmt.Println(q)

	/*
		{
		  search(text: "an") {
		    __typename
		    ... on Human {
		      name
		    }
		    ... on Droid {
		      name
		    }
		    ... on Starship {
		      name
		    }
		  }
		}
	*/
	q = fgql.NewQuery().
		Selection("search", fgql.WithArguments(
			fgql.NewArgument("text", fgql.NewStringValue("an")),
		)).
		Scalar("__typename").
		InlineFragment("Human").Scalar("name").
		Parent().
		InlineFragment("Droid").Scalar("name").
		Parent().
		InlineFragment("Starship").Scalar("name").
		Root().String()
	fmt.Println(q)
}
