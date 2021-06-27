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
		havcpydetail := false
		var errs error

		if cmpy, errs = commonfuncs.CompanyCheck(app, w, r); errs != nil {
			return
		}

		if len(*cmpy) == 1 {
			havcpydetail = true
		} else {
			havcpydetail = false
			cmpy = &[]models.TblCompany{}
		}

		//data = "Company Fetch successful"
		lgmsg := "Company fetch successful.  But havecpy detail? = " + strconv.FormatBool(havcpydetail)
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       cmpy,
			Status:     "SUCCESS",
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
