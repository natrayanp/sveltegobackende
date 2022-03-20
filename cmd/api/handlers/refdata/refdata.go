package refdata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/httpresponse"
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

		ddf := models.RefDatReqFinal{
			Refs: p,
		}

		fmt.Println("-------------------\n fetchCompany Start 1 ref \n-------------------")

		if err := commonfuncs.RefDataFetch1(app, w, r, &ddf); err != nil {
			return
		}

		fmt.Println(ddf.RefResult)
		status := "SUCCESS"
		lgmsg := "Ref Fetch successful"
		ssd := map[string]interface{}{"message": lgmsg, "refdata": ddf.RefResult}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "COMAPNY-FETCH",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()
		fmt.Println("-------------------\n fetchCompany Stop 1 ref\n-------------------")

	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(getRefdata(app))
}
