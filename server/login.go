package server

import (
	"fmt"
	"net/http"
	"reflect"
)

type LoginForm struct {
	name string;
	password string;
}


func (app *AppState) Login(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer,"invalid request", http.StatusBadRequest)
		return
	}

	form := &LoginForm{}
	ParseFormInto(request, form)
	fmt.Println(form)

	if len(form.name) == 0 || len(form.password) == 0{
		http.Error(writer,"invalid data", http.StatusBadRequest)
		return
	}

}

func ParseFormInto(r *http.Request, s interface{}) error {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := r.FormValue(field.Name)
		if value != "" {
			fieldValue := v.Field(i)
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Int:
				intValue := 0
				fmt.Sscanf(value, "%d", &intValue)
				fieldValue.SetInt(int64(intValue))
			}
		} else {
			return fmt.Errorf("lol")
		}
	}
	return nil
}
