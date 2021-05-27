package errors

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	log "github.com/sveltegobackend/pkg/logger"
)

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeDatabase       = ErrorType{"database"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)

type SlugError struct {
	Err        error
	ErrType    ErrorType
	RespWriter http.ResponseWriter
	Request    *http.Request
	Data       map[string]interface{}
	SlugCode   string
	LogMsg     string
}

type SlugError1 struct {
	Err        error
	ErrType    ErrorType
	RespWriter http.ResponseWriter
	Request    *http.Request
	Data       string
	SlugCode   string
	LogMsg     string
}

func (s SlugError) Error() error {
	return s.Err
}

func (s SlugError) GData() map[string]interface{} {
	return s.Data
}

func (s SlugError) ErrorType() ErrorType {
	return s.ErrType
}

func (s SlugError1) ErrorType() ErrorType {
	return s.ErrType
}

func (s SlugError) HttpRespondWithError() {
	log.GetLogEntry(s.Request).WithError(s.Err).WithField("error-slug", map[string]interface{}{"error": s.Data, "slugcode": s.SlugCode}).Warn(s.LogMsg)
	resp := ErrorResponse{s.Data, "ERROR", s.getStatucode()}
	fmt.Println("chek data resp")
	fmt.Println(resp)
	if err := render.Render(s.RespWriter, s.Request, &resp); err != nil {
		panic(err)
	}
}

func (s SlugError1) HttpRespondWithError() {
	log.GetLogEntry(s.Request).WithError(s.Err).WithField("error-slug", map[string]interface{}{"error": s.Data, "slugcode": s.SlugCode}).Warn(s.LogMsg)
	resp := ErrorResponse1{s.Data, "ERROR", s.getStatucode()}
	fmt.Println("chek data resp")
	fmt.Println(resp)
	if err := render.Render(s.RespWriter, s.Request, &resp); err != nil {
		panic(err)
	}
}

func NewSlugError(err error, errorType ErrorType, w http.ResponseWriter, r *http.Request, slug map[string]interface{}, slugcode string, logmsg string) SlugError {
	return SlugError{
		Err:        err,
		ErrType:    errorType,
		RespWriter: w,
		Request:    r,
		Data:       slug,
		SlugCode:   slugcode,
		LogMsg:     logmsg,
	}
}

type ErrorResponse struct {
	Data       map[string]interface{} `json:"data"`
	Status     string                 `json:"status"`
	httpStatus int
}

type ErrorResponse1 struct {
	Data       string `json:"data"`
	Status     string `json:"status"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}

func (e ErrorResponse1) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}

func (s SlugError) getStatucode() int {
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
}

func (s SlugError1) getStatucode() int {
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
}

/*
func NewAuthorizationError(error string, slug string) SlugError {
	return SlugError{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeAuthorization,
	}
}

func NewIncorrectInputError(error string, slug string) SlugError {
	return SlugError{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeIncorrectInput,
	}
}
*/
