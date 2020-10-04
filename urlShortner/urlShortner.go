package urlShortner

import (
	"fmt"
	"net/http"
	"os"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func StartHttpServer(port string, yamlFileName string) {
	defaultMux := defaultMux()
	// pathsMap := map[string]string{
	// 	"/youtube": "https://www.youtube.com",
	// }
	pathsMap := yamlParser(yamlFileName)
	fmt.Println("pathsMap: ", pathsMap)
	urlShortnerHandler := urlShortnerHandler(pathsMap, defaultMux)
	err := http.ListenAndServe(port, urlShortnerHandler)
	if err != nil {
		fmt.Println("error while starting the server: ", err)
	}
}

func urlShortnerHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fmt.Println("path is: ", path)
		if dest, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func yamlParser(yamlFilename string) map[string]string {
	yamlBytes := readYamlFile(yamlFilename)
	var paths []path
	yaml.Unmarshal(yamlBytes, &paths)
	pathsMap := buildMap(paths)
	return pathsMap
}

func buildMap(paths []path) map[string]string {
	pathsMap := make(map[string]string)
	for _, curPath := range paths {
		pathsMap[curPath.Path] = curPath.Url
	}
	return pathsMap
}

func readYamlFile(yamlFilename string) []byte {
	yamlBytes, err := ioutil.ReadFile(yamlFilename)
	if err != nil {
		fmt.Println("error while reading yaml file: ", err)
		os.Exit(1)
	}
	return yamlBytes
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

type path struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
