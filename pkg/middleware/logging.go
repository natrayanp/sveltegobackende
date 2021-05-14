package middleware

import (
	"net/http"
	//logs "github.com/sveltegobackend/pkg/logger"
)

func LogRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//logs.Info.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		//f(w, r)
	}
}
