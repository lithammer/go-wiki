package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
)

var options struct {
	Dir      string `short:"d" long:"dir" description:"Path to wiki directory" default:"wiki"`
	Template string `short:"t" long:"base-template" description:"Path to base HTML template" default:"templates/base.html"`
	Port     int    `short:"p" long:"port" description:"Port to listen on" default:"8080"`

	template *template.Template
	git      bool
}

func main() {
	_, err := flags.Parse(&options)
	if err != nil {
		os.Exit(0)
	}

	log.Println("Serving wiki from", options.Dir)
	log.Println("Using base template", options.Template)

	// Parse base template
	options.template, err = template.ParseFiles(options.Template)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	// Trim trailing slash from root path
	if strings.HasSuffix(options.Dir, "/") {
		options.Dir = options.Dir[:len(options.Dir)-1]
	}

	// Verify that the wiki folder exists
	_, err = os.Stat(options.Dir)
	if os.IsNotExist(err) {
		log.Fatalln("ERROR", err)
	}

	// Check if the wiki folder is a Git repository
	options.git = IsGitRepository(options.Dir)
	if options.git {
		log.Println("Git repository found in directory")
	} else {
		log.Println("No git repository found in directory")
	}

	r := mux.NewRouter()

	// API endpoints
	r.HandleFunc("/api/diff/{hash}/{file}", DiffHandler)

	// Static endpoints
	r.HandleFunc("/{filepath}", WikiHandler)
	r.HandleFunc("/", HomeHandler)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("public")))

	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(r)

	n.Run(fmt.Sprintf(":%d", options.Port))
}
