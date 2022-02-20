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
		if err != nil {
			return
		}

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

func SaveRoles(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n Save Role Start \n-------------------")
		var rlRes *models.RoleResp
		defer r.Body.Close()
		roledata := &models.RolesaveReq{}
		json.NewDecoder(r.Body).Decode(roledata)

		err := commonfuncs.Roleupdate(app, w, r, roledata)
		fmt.Println(err)

		fmt.Println("989899")

		if err != nil {
			return
		}

		fmt.Println(err)

		fmt.Println("-------------------\n Save Role END starting fetch \n-------------------")

		rolefetcdata := &models.RoleReq{}
		rolefetcdata.Companyid = roledata.Companyid
		rolefetcdata.Branchid = roledata.Branchid

		rlRes, err = commonfuncs.RoleFetch(app, w, r, rolefetcdata)
		if err != nil {
			return
		}

		//rlRes := &models.RoleResp{}

		fmt.Println(roledata)

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
		fmt.Println("-------------------\n Save fetch in Save Role END \n-------------------")

		cc.HttpRespond()
		return

	}
}

func DoGet(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(FetchRoles(app))
}

func DoSave(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(SaveRoles(app))
}
