package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v, found := pathsToUrls[r.URL.String()]; found {
			http.Redirect(w, r, v, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
	return hf
}

type data struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...

	var ymlData []data
	err := yaml.Unmarshal(yml, &ymlData)

	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if index := getIndex(r.URL.String(), ymlData); index != -1 {
			http.Redirect(w, r, ymlData[index].URL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})

	return hf, err
}

//JSONHandler urlshortener from json file
func JSONHandler(jsondata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var data []data
	err := json.Unmarshal(jsondata, &data)
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if index := getIndex(r.URL.String(), data); index != -1 {
			http.Redirect(w, r, data[index].URL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})

	return hf, err
}

func getIndex(item string, slice []data) int {
	for i, v := range slice {
		if item == v.Path {
			return i
		}
	}
	return -1
}
