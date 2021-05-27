package httpresponse

import (
	"fmt"
	"net/http"
	"reflect"

	log "github.com/sveltegobackend/pkg/logger"
)

/*
type SlugResp struct {
	RespWriter http.ResponseWriter
	Request    *http.Request
	Respdata   interface{}
	RespCode   string
	Statuscode int
}

func (s SlugResp) HttpRespondWithError() {
	if s.statuscode == 0 {
		s.statuscode = http.StatusOK
	}

	resp := SucessResponse{s.Respdata, s.RespCode, s.statuscode}

	if err := render.Render(s.RespWriter, s.Request, resp); err != nil {
		panic(err)
	}
}
*/

type SucessResponse struct {
	Data   map[string]interface{} `json:"data"`
	Status string                 `json:"status"`
}

type SlugSuccess struct {
	RespWriter http.ResponseWriter
	Request    *http.Request
	Data       interface{}
	Status     string
	SlugCode   string
	LogMsg     string
}

//s.Data can be stuct pointer or a map[string]interface}{} type for it to work
func (s SlugSuccess) RespData() (SucessResponse, int) {
	log.GetLogEntry(s.Request).WithField("slugcode", s.SlugCode).Info(s.LogMsg)

	resp := SucessResponse{structToMap(s.Data), s.Status}
	fmt.Println(resp)
	return resp, http.StatusOK
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}

	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()

		fmt.Println(v.Kind())

		for i := 0; i < v.NumField(); i++ {
			tag := v.Field(i).Tag.Get("json")
			field := reflectValue.Field(i).Interface()
			if tag != "" && tag != "-" {
				if v.Field(i).Type.Kind() == reflect.Struct {
					res[tag] = structToMap(field)
				} else {
					res[tag] = field
				}
			}
		}
	} else if v.Kind() == reflect.Map {
		var ok bool
		res, ok = item.(map[string]interface{})
		if !ok {
			return map[string]interface{}{}
		}
	}

	fmt.Println(res)
	return res
}
