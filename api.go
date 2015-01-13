package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shurcooL/go/github_flavored_markdown"
)

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	file := vars["file"] + ".md"

	diff, err := Diff(file, vars["hash"])
	if err != nil {
		log.Println("ERROR", "Failed to get commit hash", vars["hash"])
	}

	// XXX: This could probably be done in a nicer way
	wrappedDiff := []byte("```diff\n" + string(diff) + "```")
	// md := blackfriday.MarkdownCommon(wrappedDiff)
	md := github_flavored_markdown.Markdown(wrappedDiff)

	w.Header().Set("Content-Type", "text/html")
	w.Write(md)
}
