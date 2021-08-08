package branch

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

func fetchBranch(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n fetchBranch Start \n-------------------")

		var brch *[]models.TblBranch
		var brchcp []models.TblBranch
		havcpydetail := false
		var errs error
		var status string

		if brch, errs = commonfuncs.BranchCheck(app, w, r); errs != nil {
			return
		}
		brchcp = *brch

		if len(*brch) == 1 {
			havcpydetail = true
			status = "SUCCESS"
		} else {
			havcpydetail = false
			status = "FAILURE"
			brchcp = []models.TblBranch{}
		}

		dd := []models.RefDatReq{
			{Reftype: "group", Refname: "branch"},
		}

		ddf := models.RefDatReqFinal{
			Refs: dd,
		}

		fmt.Println("-------------------\n fetchBranch Start 1 \n-------------------")

		if err := commonfuncs.RefDataFetch1(app, w, r, &ddf); err != nil {
			return
		}

		fmt.Println(ddf.RefResult)

		fmt.Println("-------------------\n fetchBranch Start 2 \n-------------------")

		lgmsg := "Branch Fetch successful.  But havecpy detail? = " + strconv.FormatBool(havcpydetail)
		ssd := map[string]interface{}{"message": lgmsg, "branch": brchcp, "refdata": ddf.RefResult}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "BRANCH-FETCH",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()
		fmt.Println("-------------------\n fetchBranch Stop \n-------------------")
		//return

	}
}

func DoBrFetch(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(fetchBranch(app))
}
