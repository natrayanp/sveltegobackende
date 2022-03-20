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

func FetchUserMatrix(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n Fetch Role Start \n-------------------")
		var datosend models.UserMatrixResp

		defer r.Body.Close()

		matrixreq := &models.UserMatrixReq{}
		json.NewDecoder(r.Body).Decode(matrixreq)
		fmt.Println(matrixreq)

		matxRes, err := commonfuncs.UserMatrixFetch(app, w, r, matrixreq)
		if err != nil {
			return
		}

		datosend.Listmatrix = matxRes
		datosend.Resptype = matrixreq.Optype

		status := "SUCCESS"
		lgmsg := "Role Fetch successful"
		ssd := map[string]interface{}{"message": lgmsg, "matrixdata": datosend}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "USRMATRIX-FETCH",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()
	}
}

func SaveUserMatrix(app *application.Application) http.HandlerFunc {
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

func DoGetmatrix(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(FetchUserMatrix(app))
}

func DoSavematrix(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(SaveUserMatrix(app))
}
