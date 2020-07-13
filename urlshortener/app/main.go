package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/JordiGilTomas/coursecalhoun/urlshortener/urlshort"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var yamlFile string
var jsonFile string

func init() {
	flag.StringVar(&yamlFile, "yaml", "routes.yaml", "Carga un fichero YAML")
	flag.StringVar(&jsonFile, "json", "routes.json", "Carga un fichero YAML")

}
func main() {
	flag.Parse()
	dbRoutes := getFirestoreRoutes()

	file, err := os.Open(yamlFile)
	if err != nil {
		fmt.Println("No se pudo cargar el fichero yaml")
		os.Exit(3)
	}
	fileJSON, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("No se pudo cargar el fichero json")
		os.Exit(3)
	}
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
	dbHandler, err := urlshort.DBHandler(dbRoutes, jsonHandler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8087")
	http.ListenAndServe(":8087", dbHandler)
}

func getFirestoreRoutes() []map[string]interface{} {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("firebase.json"))
	db, err := app.Firestore(ctx)

	if err != nil {
		fmt.Println("Error al conectar a la base de datos", db, err)
		os.Exit(3)
	}

	routes := db.Collection("routes").Documents(ctx)
	var dbRoutes []map[string]interface{}

	for {
		doc, err := routes.Next()
		if err == iterator.Done {
			break
		}
		dbRoutes = append(dbRoutes, doc.Data())
	}

	for _, docs := range dbRoutes {
		for k, v := range docs {
			dbRoutes = append(dbRoutes, map[string]interface{}{k: v})
		}
	}
	return dbRoutes
}
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
