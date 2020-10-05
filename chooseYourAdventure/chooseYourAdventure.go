package chooseYourAdventure

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Option struct {
	Text string
	Arc  string
}

type Adventure struct {
	Title   string
	Story   []string
	Options []Option
}

func (adventure Adventure) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, adventure)
	adventureTemplate := template.Must(template.New("adventure.html").ParseFiles("adventure.html"))
	fmt.Println("Defined templates: ", adventureTemplate.DefinedTemplates())
	executeErr := adventureTemplate.ExecuteTemplate(w, "adventure.html", adventure)
	fmt.Println("am i coming here?")
	if executeErr != nil {
		fmt.Println(executeErr)
		fmt.Fprintln(w, "Unable to render: ", executeErr)
	}
}

func ParseJson(jsonFileName string) map[string]Adventure {
	jsonBytes, jsonRdErr := ioutil.ReadFile(jsonFileName)
	if jsonRdErr != nil {
		fmt.Println("error while reading json file: ", jsonRdErr)
		os.Exit(1)
	}
	adventures := make(map[string]Adventure)
	unmarshalErr := json.Unmarshal(jsonBytes, &adventures)
	if unmarshalErr != nil {
		fmt.Println("error while unmarshalling json file: ", unmarshalErr)
		os.Exit(1)
	}
	return adventures
}

func StartServer(address string, jsonFileName string) {
	adventuresJson := ParseJson(jsonFileName)
	for advKey, advValue := range adventuresJson {
		http.Handle("/"+advKey, advValue)
	}
	startServerErr := http.ListenAndServe(address, nil)
	if startServerErr != nil {
		panic(startServerErr)
	}
}
