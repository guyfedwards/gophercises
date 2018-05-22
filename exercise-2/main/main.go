package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	bolt "github.com/coreos/bbolt"
	"github.com/guyfedwards/gophercises/exercise-2/urlshort"
)

func main() {
	filePtr := flag.String("file", "", "Path to config file")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	var yaml string

	if *filePtr != "" {
		dat, err := ioutil.ReadFile(*filePtr)
		if err != nil {
			panic(err)
		}

		yaml = string(dat)
	} else {
		yaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	}
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json := `
[{
	"path": "/nerd",
	"url": "https://google.com"
}]
	`
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DBHandler := urlshort.DBHandler(db, jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", DBHandler)
}

func defaultMux() *http.ServeMux {
	fmt.Println("hi")
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FCIA")
	fmt.Fprintln(w, "Hello, world!")
}
