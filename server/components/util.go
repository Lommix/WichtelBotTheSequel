package components

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

// load snippet from json file
func LoadSnippets(lang string, path string) ( map[string]interface{}, error ) {

	var out = make(map[string] interface{})
	// read a file from disc
	snippets, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	s := string(snippets)
	var data map[string]map[string]string
	json.Unmarshal([]byte(s), &data)

	for key, snippet := range data {
		out[key] = snippet[string(lang)]
	}

	return out, nil
}


func LangFromRequest(r *http.Request) Language {
	if r.Header.Get("Accept-Language") == "de" {
		return German
	}
	return English
}


// load from data from request into any interface respecting required attributes: `required:"true"`
func FromFormData(r *http.Request, s interface{}) error {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		formValue := r.FormValue(t.Field(i).Name)
		if formValue != "" && v.CanSet(){
			fieldValue := v.Field(i)
			switch fieldValue.Kind() {
			case reflect.String:
				v.Field(i).SetString(formValue)
			case reflect.Int:
				intValue := 0
				fmt.Sscanf(formValue, "%d", &intValue)
				v.Field(i).SetInt(int64(intValue))
			}
		} else {
			if t.Field(i).Tag.Get("required") == "true" {
				return fmt.Errorf("missing required field %s", t.Field(i).Name)
			}
		}
	}
	return nil
}


// reading the frist slug from an url
func GetFirstSlug(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return ""
	}

	host := parsedURL.Host
	if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
	}
	// Split the path into slugs
	pathSlugs := strings.Split(parsedURL.Path, "/")
	if len(pathSlugs) == 0 {
		return ""
	}
	return pathSlugs[0]
}

