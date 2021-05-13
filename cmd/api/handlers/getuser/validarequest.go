package getuser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sveltegobackende/cmd/api/models"
)

func validateRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := r.URL.Query()["id"]

		//id, err := strconv.Atoi(uid)
		if !err || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, "malformed id")
			return
		}
		id := keys[0]

		ctx := context.WithValue(r.Context(), models.CtxKey("userid"), id)
		r = r.WithContext(ctx)
		//next(w, r)
	}
}
