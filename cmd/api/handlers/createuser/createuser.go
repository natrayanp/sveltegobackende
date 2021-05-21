package createuser

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func createUser(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		user := &models.User{}
		json.NewDecoder(r.Body).Decode(user)

		if err := user.Create(r.Context(), app); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Oops")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(user)
		w.Write(response)
	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(createUser(app), mymiddleware.LogRequest)
	//return createUser(app)
}
