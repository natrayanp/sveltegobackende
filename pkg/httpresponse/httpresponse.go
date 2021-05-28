package httpresponse

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/render"

	log "github.com/sveltegobackend/pkg/logger"
)

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown         = ErrorType{"unknown"}
	ErrorTypeDatabase        = ErrorType{"database"}
	ErrorTypeAuthorization   = ErrorType{"authorization"}
	ErrorTypeIncorrectInput  = ErrorType{"incorrect-input"}
	ErrorTypeContexFetchFail = ErrorType{"usercontext-fetch-fail"}
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

type finalResponse struct {
	Data       map[string]interface{} `json:"data"`
	Status     string                 `json:"status"`
	httpStatus int
}

type SlugResponse struct {
	Err        error
	ErrType    ErrorType
	RespWriter http.ResponseWriter
	Request    *http.Request
	Data       interface{}
	Status     string
	SlugCode   string
	LogMsg     string
}

func (s SlugResponse) ErrorType() ErrorType {
	return s.ErrType
}

//s.Data can be stuct pointer or a map[string]interface}{} type for it to work
/*
func (s SlugResponse) RespData() (FinalResponse, int) {
	if s.Err != nil {
		log.GetLogEntry(s.Request).WithError(s.Err).WithField("error-slug", map[string]interface{}{"error": s.Data, "slugcode": s.SlugCode}).Warn(s.LogMsg)
	} else {
		log.GetLogEntry(s.Request).WithField("slugcode", s.SlugCode).Info(s.LogMsg)
	}

	resp := FinalResponse{structToMap(s.Data), s.Status, s.getStatucode()}
	fmt.Println(resp)
	return resp, http.StatusOK
}
*/

func (s SlugResponse) HttpRespond() {
	if s.Err != nil {
		log.GetLogEntry(s.Request).WithError(s.Err).WithField("error-slug", map[string]interface{}{"error": s.Data, "slugcode": s.SlugCode}).Warn(s.LogMsg)
	} else {
		log.GetLogEntry(s.Request).WithField("slugcode", s.SlugCode).Info(s.LogMsg)
	}

	resp := finalResponse{structToMap(s.Data), s.Status, s.getStatucode()}
	fmt.Println("chek data resp")
	fmt.Println(resp)
	if err := render.Render(s.RespWriter, s.Request, &resp); err != nil {
		panic(err)
	}
}

func (s SlugResponse) getStatucode() int {
	if s.Err != nil {

		switch s.ErrorType() {
		case ErrorTypeAuthorization:
			return http.StatusUnauthorized
		case ErrorTypeIncorrectInput:
			return http.StatusBadRequest
		case ErrorTypeDatabase:
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
		}
	} else {
		return http.StatusOK
	}
}

func (e finalResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	//w.Header().Set("Content-Type", "application/json")
	return nil
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
