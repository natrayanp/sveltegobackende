package roles

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

func FetchRoles(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n Fetch Role Start \n-------------------")

		defer r.Body.Close()

		roledata := &models.RoleReq{}
		json.NewDecoder(r.Body).Decode(roledata)
		fmt.Println(roledata)

		rlRes, err := commonfuncs.RoleFetch(app, w, r, roledata)

		fmt.Println(err)

		fmt.Println(rlRes)
		status := "SUCCESS"
		lgmsg := "Role Fetch successful"
		ssd := map[string]interface{}{"message": lgmsg, "roledata": rlRes}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "ROLE-FETCH",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()

	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(FetchRoles(app))
}
