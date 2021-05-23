package errors

import (
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
	Slug       string
	SlugCode   string
	LogMsg     string
}

func (s SlugError) Error() error {
	return s.Err
}

func (s SlugError) GSlug() string {
	return s.Slug
}

func (s SlugError) ErrorType() ErrorType {
	return s.ErrType
}

func (s SlugError) HttpRespondWithError() {
	log.GetLogEntry(s.r).WithError(s.error).WithField("error-slug", s.slug).Warn(s.logMsg)
	resp := ErrorResponse{s.slug, s.slugCode, s.getStatucode()}

	if err := render.Render(s.w, s.r, resp); err != nil {
		panic(err)

	}
}

func NewSlugError(err error, errorType ErrorType, w http.ResponseWriter, r *http.Request, slug string, slugcode string, logmsg string) SlugError {
	return SlugError{
		error:     err,
		errorType: errorType,
		w:         w,
		r:         r,
		slug:      slug,
		slugCode:  slugcode,
		logMsg:    logmsg,
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Slugcode   string `json:"slugcode"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
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
