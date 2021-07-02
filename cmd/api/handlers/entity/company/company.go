package company

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func fetchCompany(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n fetchCompany Start \n-------------------")

		var cmpy *[]models.TblCompany
		var cmpycp []models.TblCompany
		havcpydetail := false
		var errs error
		var status string

		if cmpy, errs = commonfuncs.CompanyCheck(app, w, r); errs != nil {
			return
		}
		cmpycp = *cmpy

		if len(*cmpy) == 1 {
			havcpydetail = true
			status = "SUCCESS"
		} else {
			havcpydetail = false
			status = "FAILURE"
			cmpycp = []models.TblCompany{}
		}

		dd := []models.RefDatReq{
			{Reftype: "group", Refname: "company"},
		}

		ddf := models.RefDatReqFinal{
			Refs: dd,
		}

		fmt.Println("-------------------\n fetchCompany Start 1 \n-------------------")

		if err := commonfuncs.RefDataFetch1(app, w, r, &ddf); err != nil {
			return
		}

		fmt.Println(ddf.RefResult)

		fmt.Println("-------------------\n fetchCompany Start 2 \n-------------------")

		lgmsg := "Company Fetch successful.  But havecpy detail? = " + strconv.FormatBool(havcpydetail)
		ssd := map[string]interface{}{"message": lgmsg, "company": cmpycp, "refdata": ddf.RefResult}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "COMAPNY-FETCH",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()
		fmt.Println("-------------------\n fetchCompany Stop \n-------------------")
		return

	}
}

func DoFetch(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(fetchCompany(app))
}

func DoSave(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(fetchCompany(app))
}
