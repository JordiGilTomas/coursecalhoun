package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

//MapHandler urlshortener from variable
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	hf := func(w http.ResponseWriter, r *http.Request) {
		if v, found := pathsToUrls[r.URL.String()]; found {
			http.Redirect(w, r, v, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
	return hf
}

type data struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

//YAMLHandler urlshortener from yaml file
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var ymlData []data
	err := yaml.Unmarshal(yml, &ymlData)
	if err != nil {
		return nil, err
	}
	mapData := makeMap(ymlData)
	return MapHandler(mapData, fallback), nil
}

//JSONHandler urlshortener from json file
func JSONHandler(jsondata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var data []data
	err := json.Unmarshal(jsondata, &data)
	if err != nil {
		return nil, err
	}
	mapData := makeMap(data)
	return MapHandler(mapData, fallback), nil
}

//DBHandler urlshortener from Database Firestore
func DBHandler(db []map[string]interface{}, fallback http.Handler) (http.HandlerFunc, error) {
	data := dataToJSON(db)
	mapData := makeMap(data)
	return MapHandler(mapData, fallback), nil
}

func dataToJSON(routes []map[string]interface{}) []data {
	var data []data
	dbdata, err := json.Marshal(routes)
	if err != nil {
		fmt.Println("Error al convertir en Marhsal")
		os.Exit(3)
	}
	err = json.Unmarshal(dbdata, &data)
	if err != nil {
		fmt.Println("Error al convertir en Unmarshal")
		os.Exit(3)
	}
	return data
}

func makeMap(data []data) map[string]string {
	mapData := make(map[string]string)
	for _, v := range data {
		mapData[v.Path] = v.URL
	}
	return mapData
}
