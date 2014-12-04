package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

const imageTypes = ".jpg .jpeg .png .gif"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}

func WikiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Deny requests trying to traverse up the directory structure using
	// relative paths
	if strings.Contains(vars["filepath"], "..") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Path to the file as it is on the the local file system
	realpath := fmt.Sprintf("%s/%s", options.Dir, vars["filepath"])

	// Serve (accepted) images
	for _, filext := range strings.Split(imageTypes, " ") {
		if path.Ext(r.URL.Path) == filext {
			http.ServeFile(w, r, realpath)
			return
		}
	}

	md, err := ioutil.ReadFile(realpath + ".md")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	wiki := Wiki{
		Markdown: md,
		filepath: realpath,
		template: options.template,
	}

	wiki.Commits, err = Commits(vars["filepath"]+".md", 5)
	if err != nil {
		log.Println("ERROR", "Failed to get commits")
	}

	wiki.Write(w)
}
