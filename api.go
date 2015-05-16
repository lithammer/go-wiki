package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
)

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[len("/api/diff/"):], "/")
	if len(parts) != 2 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	hash := parts[0]
	file := parts[1] + ".md"

	diff, err := Diff(file, hash)
	if err != nil {
		log.Println("ERROR", "Failed to get commit hash", hash)
	}

	// XXX: This could probably be done in a nicer way
	wrappedDiff := []byte("```diff\n" + string(diff) + "```")
	// md := blackfriday.MarkdownCommon(wrappedDiff)
	md := github_flavored_markdown.Markdown(wrappedDiff)

	w.Header().Set("Content-Type", "text/html")
	w.Write(md)
}
