package components

import (
	"encoding/json"
	"os"
)

const SnippetPath = "snippets.json"

type Language string

const (
	German  Language = "de"
	English Language = "en"
)

type Snippets struct {
	data *map[Language]map[string]string
}

func (snippets *Snippets) Load() error {
	var err error
	snippets.data, err = LoadSnippetsFromDisk(SnippetPath)
	return err
}

// get a single snippet
func (snippets *Snippets) Get(name string, lang Language) (string,error) {
	if snippets.data == nil {
		err := snippets.Load()
		if err != nil {
			return "", err
		}
	}
	return (*snippets.data)[lang][name], nil
}

// get a list for a specific language
func (snippets *Snippets) GetList(lang Language) *map[string]string {
	list := (*snippets.data)[lang]
	return &list
}

func LoadSnippetsFromDisk(path string) (*map[Language]map[string]string, error ) {
	out := make(map[Language]map[string]string)
	in := make(map[string]map[string]string)

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return &out, err
	}

	err = json.Unmarshal(fileContent, &in)
	if err != nil {
		return &out, err
	}

	for key, snip := range(in){
		for	lang, value := range(snip){
			if out[Language(lang)] == nil {
				out[Language(lang)] = make(map[string]string)
			}
			out[Language(lang)][key]= value
		}
	}

	return &out, nil
}

