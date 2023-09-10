package server

import (
	"fmt"
	"net/http"
	"reflect"
)

type LoginForm struct {
	Username string;
	Password string;
}


func (app *AppState) Login(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer,"invalid request", http.StatusBadRequest)
		return
	}

	form := &LoginForm{}
	ParseFormInto(request, form)

	fmt.Println(form.Username)
	fmt.Println(form.Password)

	if len(form.Username) == 0 || len(form.Password) == 0{
		http.Error(writer,"invalid data", http.StatusBadRequest)
		return
	}

}

func ParseFormInto(r *http.Request, s interface{}) error {
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
			return fmt.Errorf("lol")
		}
	}
	return nil
}
