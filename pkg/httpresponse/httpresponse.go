package httpresponse

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/render"
	"github.com/sveltegobackend/pkg/fireauth"
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

type finalResponse struct {
	Data       map[string]interface{} `json:"data"`
	Status     string                 `json:"status"`
	Session    string                 `json:"session"`
	SlugCode   string                 `json:"slugcode"`
	httpStatus int
}

type SlugResponse struct {
	Err        error
	ErrType    ErrorType
	RespWriter http.ResponseWriter
	Request    *http.Request
	Userinfo   fireauth.User
	Data       interface{}
	Status     string
	SlugCode   string
	LogMsg     string
}

func (s SlugResponse) ErrorType() ErrorType {
	return s.ErrType
}

//s.Data can be stuct pointer or a map[string]interface}{} type for it to work
func (s SlugResponse) HttpRespond() {
	fmt.Println(s.Err)
	fmt.Println(s.Err != nil)
	if s.Err != nil {
		log.GetLogEntry(s.Request).WithError(s.Err).WithField("error-slug", map[string]interface{}{"error": s.Data, "slugcode": s.SlugCode}).Warn(s.LogMsg)
		if s.Status == "" {
			s.Status = "ERROR"
		}
	} else {
		log.GetLogEntry(s.Request).WithField("slugcode", s.SlugCode).Info(s.LogMsg)
	}
	fmt.Println(s)
	fmt.Println("s.Status")
	fmt.Println(s.Status)
	fmt.Println("s.Status")
	resp := finalResponse{structToMap(s.Data), s.Status, s.Userinfo.Session, s.SlugCode, s.getStatucode()}

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
	fmt.Println("type", v)
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
	fmt.Println("return from println")
	fmt.Println(res)
	return res
}
