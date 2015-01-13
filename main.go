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
	flag "github.com/ogier/pflag"
)

var options struct {
	Dir      string
	Template string
	Static   string
	Port     int

	template *template.Template
	git      bool
}

var usage = `Usage: gowiki [options...] <path>

Positional arguments:
  path                  directory to serve wiki pages from

Optional arguments:
  -h, --help            show this help message and exit
  -p PORT, --port=PORT  listen port (default 8080)
  -t FILE, --base-template=FILE
                        base HTML template (default /usr/local/share/gowiki/templates/base.html)
  -s PATH, --static=PATH
                        static files folder (default /usr/local/share/gowiki/public)
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.StringVarP(&options.Template, "base-template", "t", "/usr/local/share/gowiki/templates/base.html", "")
	flag.StringVarP(&options.Static, "static", "s", "/usr/local/share/gowiki/templates/base.html", "")
	flag.IntVarP(&options.Port, "port", "p", 8080, "")

	flag.Parse()

	options.Dir = flag.Arg(0)

	if options.Dir == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Serving wiki from", options.Dir)
	log.Println("Using base template", options.Template)

	// Parse base template
	var err error
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

	r.HandleFunc("/api/diff/{hash}/{file}", DiffHandler)
	r.HandleFunc("/{filepath}", WikiHandler)
	r.HandleFunc("/", HomeHandler)

	n := negroni.New()

	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewStatic(http.Dir(options.Static)))
	n.UseHandler(r)

	n.Run(fmt.Sprintf(":%d", options.Port))
}
