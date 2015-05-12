package main

import (
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/renstrom/go-wiki/vendor/_nuts/github.com/russross/blackfriday"
)

type Wiki struct {
	Body      template.HTML
	Markdown  []byte
	Commits   []Commit
	CustomCSS string

	template *template.Template
	filepath string
}

func (w Wiki) Title() string {
	_, file := path.Split(w.filepath)
	file = strings.Replace(file, "_", " ", -1)
	file = strings.Title(file)
	return file
}

func (w *Wiki) Write(rw http.ResponseWriter) {
	w.Body = template.HTML(blackfriday.MarkdownCommon(w.Markdown))
	err := w.template.Execute(rw, w)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
