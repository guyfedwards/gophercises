package urlshort

import (
	"fmt"
	"net/http"

	"encoding/json"

	bolt "github.com/coreos/bbolt"
	"gopkg.in/yaml.v2"
)

type url struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		url := pathsToUrls[path]

		if url == "" {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, url, 302)
	}
}

func createMap(y []url) map[string]string {
	m := make(map[string]string)

	for _, value := range y {
		m[value.Path] = value.URL
	}

	return m
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
	var y []url

	err := yaml.Unmarshal(yml, &y)
	if err != nil {
		return nil, err
	}

	m := createMap(y)
	return MapHandler(m, fallback), nil
}

// JSONHandler is same as YAMLHandler but for JSON
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var j []url
	err := json.Unmarshal(jsn, &j)
	if err != nil {
		return nil, err
	}

	m := createMap(j)
	return MapHandler(m, fallback), nil
}

// DBHandler reads key/values as paths/urls from the DB
func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	var path []byte

	return func(w http.ResponseWriter, r *http.Request) {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("mappings"))
			v := b.Get([]byte(r.URL.Path))
			path = make([]byte, len(v))

			copy(path, v)

			return nil
		})

		p := string(path)
		if p == "" {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, p, 302)

		fmt.Println("derp", string(path))
	}
}
