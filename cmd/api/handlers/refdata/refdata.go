package refdata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func getRefdata(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n getRefdata Start \n-------------------")

		defer r.Body.Close()
		var p []models.RefDatReq

		err := json.NewDecoder(r.Body).Decode(&p)
		fmt.Println(p)

		if err != nil {
			return
		}

		fmt.Println("-------------------\n getRefdata Stop \n-------------------")
		return

	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(getRefdata(app))
}
