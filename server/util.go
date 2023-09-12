package server

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)



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


func getFirstSlug(urlString string) string {
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

