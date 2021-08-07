package httperr

import (
	"net/http"

	"github.com/go-chi/render"
	logs "github.com/sveltegobackend/pkg/logger"
)

func InternalError(slug string, logMSg string, err error, w http.ResponseWriter, r *http.Request) {
	if slug == "" {
		slug = "Internal server error"
	}
	if logMSg == "" {
		logMSg = "Internal server error"
	}
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Error(logMSg)
	httpRespondWithError(err, slug, w, r, logMSg, http.StatusInternalServerError)
}

func Unauthorised(slug string, logMSg string, err error, w http.ResponseWriter, r *http.Request) {
	if slug == "" {
		slug = "AUTH-FAIL"
	}
	if logMSg == "" {
		logMSg = "AUTH-FAIL"
	}
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMSg)
	httpRespondWithError(err, "AUTH-FAIL", w, r, logMSg, http.StatusUnauthorized)
}

func BadRequest(slug string, logMSg string, err error, w http.ResponseWriter, r *http.Request) {
	if slug == "" {
		slug = "Bad request"
	}
	if logMSg == "" {
		logMSg = "Bad request"
	}
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMSg)
	httpRespondWithError(err, "Bad request", w, r, logMSg, http.StatusBadRequest)
}

func DataBaseError(slug string, logMSg string, err error, w http.ResponseWriter, r *http.Request) {
	if logMSg == "" {
		logMSg = "Database Error"
	}
	InternalError(slug, logMSg, err, w, r)
}

/*
func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Println("******* indise slug")
	slugError, ok := err.(errors.SlugError)
	fmt.Println(slugError)
	fmt.Println("err type:   ", slugError.ErrorType())
	fmt.Println(ok)
	if !ok {
		InternalError("internal-server-error", "", err, w, r)
		return
	}
	fmt.Println("after err type:   ")
	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorised(slugError.Slug(), "", slugError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), "", slugError, w, r)
	default:
		InternalError(slugError.Slug(), "", slugError, w, r)
	}
}
*/

func httpRespondWithError(err error, slugcode string, w http.ResponseWriter, r *http.Request, logMSg string, status int) {
	//logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMSg)
	resp := ErrorResponse{slugcode, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)

	}
	return
}

type ErrorResponse struct {
	Slugcode   string `json:"slugcode"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
