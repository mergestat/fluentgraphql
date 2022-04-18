package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	fgql "github.com/mergestat/fluentgraphql"
)

var (
	githubToken = os.Getenv("GITHUB_TOKEN")
	// repoList taken from here: https://github.com/github/explore/blob/main/collections/front-end-javascript-frameworks/index.md
	repoList = []string{
		"marko-js/marko",
		"mithriljs/mithril.js",
		"angular/angular",
		"emberjs/ember.js",
		"knockout/knockout",
		"tastejs/todomvc",
		"spine/spine",
		"vuejs/vue",
		"Polymer/polymer",
		"facebook/react",
		"finom/seemple",
		"aurelia/framework",
		"optimizely/nuclear-js",
		"jashkenas/backbone",
		"dojo/dojo",
		"jorgebucaran/hyperapp",
		"riot/riot",
		"daemonite/material",
		"polymer/lit-element",
		"aurelia/aurelia",
		"sveltejs/svelte",
		"neomjs/neo",
		"preactjs/preact",
	}
)

func main() {
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

	fmt.Println(q.Root().String())

	body := map[string]interface{}{
		"query": q.Root().String(),
	}

	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", githubToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resBody))
}
