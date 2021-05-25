package httpresponse

import (
	"net/http"
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
	Data       string `json:"data"`
	Respcode   string `json:"respcode"`
	httpStatus int
}

func (e SucessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
