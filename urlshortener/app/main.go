package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/JordiGilTomas/coursecalhoun/urlshortener/urlshort"
)

var yamlFile string
var jsonFile string

func init() {
	flag.StringVar(&yamlFile, "yaml", "routes.yaml", "Carga un fichero YAML")
	flag.StringVar(&jsonFile, "json", "routes.json", "Carga un fichero YAML")
}
func main() {
	flag.Parse()
	file, err := os.Open(yamlFile)
	fileJSON, err := os.Open(jsonFile)
	yaml, err := ioutil.ReadAll(file)
	jsonData, err := ioutil.ReadAll(fileJSON)
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// if yaml == "" {
	// 	yaml = `
	// 	- path: /urlshort
	// 	  url: https://github.com/gophercises/urlshort
	// 	- path: /urlshort-final
	// 	  url: https://github.com/gophercises/urlshort/tree/solution
	// 	  `
	// }
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)

	jsonHandler, err := urlshort.JSONHandler([]byte(jsonData), yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8086")
	http.ListenAndServe(":8086", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
