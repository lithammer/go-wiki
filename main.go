package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	flag "github.com/ogier/pflag"
)

const Usage = `Usage: gowiki [options...] <path>

Positional arguments:
  path                  directory to serve wiki pages from

Optional arguments:
  -h, --help            show this help message and exit
  -p PORT, --port=PORT  listen port (default 8080)
  --custom-css=PATH     path to custom CSS file
`

var options struct {
	Dir       string
	Port      int
	CustomCSS string

	template *template.Template
	git      bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, Usage)
	}

	flag.IntVarP(&options.Port, "port", "p", 8080, "")
	flag.StringVar(&options.CustomCSS, "custom-css", "", "")

	flag.Parse()

	options.Dir = flag.Arg(0)

	if options.Dir == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Serving wiki from", options.Dir)

	// Parse base template
	var err error
	options.template, err = template.New("base").Parse(Template)
	if err != nil {
		log.Fatalln("Error parsing HTML template:", err)
	}

	// Trim trailing slash from root path
	if strings.HasSuffix(options.Dir, "/") {
		options.Dir = options.Dir[:len(options.Dir)-1]
	}

	// Verify that the wiki folder exists
	_, err = os.Stat(options.Dir)
	if os.IsNotExist(err) {
		log.Fatalln("Directory not found")
	}

	// Check if the wiki folder is a Git repository
	options.git = IsGitRepository(options.Dir)
	if options.git {
		log.Println("Git repository found in directory")
	} else {
		log.Println("No git repository found in directory")
	}

	http.Handle("/api/diff/", commonHandler(DiffHandler))
	http.Handle("/", commonHandler(WikiHandler))

	log.Println("Listening on:", options.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", options.Port), nil)
}
