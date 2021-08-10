package branch

import (
	"encoding/json"
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

		if len(*brch) > 0 {
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
		fmt.Println(brchcp)
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
		return

	}
}

func saveBranch(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n save Branch Start \n-------------------")

		defer r.Body.Close()

		brncdata := &models.BrnResp{}
		json.NewDecoder(r.Body).Decode(brncdata)
		fmt.Println(brncdata)

		var brnc *[]models.TblBranch
		var brnccp []models.TblBranch
		havcpydetail := false
		var errs error
		var status string
		brnccp = []models.TblBranch{}
		//brn := &models.Brn{}
		brn := &brncdata.Branchdata

		fmt.Println(brn)
		fmt.Println(brncdata.Optype)

		if brncdata.Optype == "Update" {
			if brnc, errs = commonfuncs.BranchCheck(app, w, r); errs != nil {
				return
			}
			brnccp = *brnc
			/*
				const layoutISO = "2006-01-02"
				brdd, _ := time.Parse(layoutISO, cpy.CompanyStartDate)
				brnc.BranchStartDate = cpydd.String()
			*/
			if len(*brnc) == 1 {
				havcpydetail = true
				status = "SUCCESS"

				brn1 := models.Brn{
					CompanyId:         brnccp[0].Companyid.String,
					BranchId:          brnccp[0].Branchid.String,
					BranchName:        brnccp[0].Branchname.String,
					BranchShortName:   brnccp[0].Branchshortname.String,
					BranchCategory:    brnccp[0].Branchcategory.String,
					BranchStatus:      brnccp[0].Branchstatus.String,
					BranchDescription: brnccp[0].Branchdescription.String,
					BranchImageUrl:    brnccp[0].Branchimageurl.String,
					BranchAddLine1:    brnccp[0].Branchaddline1.String,
					BranchAddLine2:    brnccp[0].Branchaddline2.String,
					BranchCountry:     brnccp[0].Branchcountry.String,
					BranchState:       brnccp[0].Branchstate.String,
					BranchCity:        brnccp[0].Branchcity.String,
					BranchPinCode:     brnccp[0].Branchpincode.String,
					BranchPhone:       brnccp[0].Branchphone.String,
					BranchFax:         brnccp[0].Branchfax.String,
					BranchMobile:      brnccp[0].Branchmobile.String,
					BranchEmail:       brnccp[0].Branchemail.String,
					BranchWebsite:     brnccp[0].Branchwebsite.String,
					BranchStartDate:   brnccp[0].Branchstartdate.Time.String(),
					Isdefault:         brnccp[0].Isdefault.String,
				}

				if brnc, errs = commonfuncs.Branchupdate(app, w, r, brn, &brn1); errs != nil {
					return
				}
				brnccp = *brnc

			} else {
				havcpydetail = true
				status = "ERROR"
				//TODO: send error response.
			}

		} else if brncdata.Optype == "Save" {
			havcpydetail = false
			status = "SUCCESS"

			if brnc, errs = commonfuncs.BranchSave(app, w, r, brn); errs != nil {
				fmt.Println("I am returning/n")
				return
			}
			brnccp = *brnc
			fmt.Println("after fetch/n")

		}
		//	}

		fmt.Println("-------------------\n fetchBranch in save company Start 1 \n-------------------")
		fmt.Println(status)
		if status == "SUCCESS" {
			if brnc, errs = commonfuncs.BranchCheck(app, w, r); errs != nil {
				return
			}
			brnccp = *brnc
		}
		fmt.Println("-------------------\n fetchBranch in save company  Start 2 \n-------------------")

		lgmsg := "Branch Save successful.  But havecpy detail? = " + strconv.FormatBool(havcpydetail)
		//ssd := map[string]interface{}{"message": lgmsg, "company": brnccp, "refdata": ddf.RefResult}
		ssd := map[string]interface{}{"message": lgmsg, "branch": brnccp, "refdata": "s"}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "BRANCH-SAVE",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()
		fmt.Println("-------------------\n save Branch Stop \n-------------------")
		return

	}
}

func DoBrFetch(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(fetchBranch(app))
}

func DoBrSave(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(saveBranch(app))
}
