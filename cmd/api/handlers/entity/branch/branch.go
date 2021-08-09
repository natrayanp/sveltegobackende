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

func saveBranch(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n save Company Start \n-------------------")
		/*
			defer r.Body.Close()

			user := &models.User{}
			json.NewDecoder(r.Body).Decode(user)

			var cmpy *[]models.TblCompany
			var cmpycp []models.TblCompany
			havcpydetail := false
			var errs error
			var status string
			cmpycp = []models.TblCompany{}

			if cpyop.Optype == "update" {
				if cmpy, errs = commonfuncs.CompanyCheck(app, w, r); errs != nil {
					return
				}
				cmpycp = *cmpy
				/*
					const layoutISO = "2006-01-02"
					cpydd, _ := time.Parse(layoutISO, cpy.CompanyStartDate)
					cpy.CompanyStartDate = cpydd.String()
		*/
		/*			if len(*cmpy) == 1 {
				havcpydetail = true
				status = "SUCCESS"

				cpy1 := models.Cpy{
					CompanyId:          cmpycp[0].Companyid.String,
					CompanyName:        cmpycp[0].Companyname.String,
					CompanyShortName:   cmpycp[0].Companyshortname.String,
					CompanyAddLine1:    cmpycp[0].Companyaddline1.String,
					CompanyAddLine2:    cmpycp[0].Companyaddline2.String,
					CompanyCategory:    cmpycp[0].Companycategory.String,
					CompanyStatus:      cmpycp[0].Companystatus.String,
					CompanyLogoUrl:     cmpycp[0].Companyimageurl.String,
					CompanyLogo:        cmpycp[0].Companylogo.String,
					CompanyIndustry:    cmpycp[0].Companyindustry.String,
					CompanyTaxID:       cmpycp[0].Companytaxid.String,
					CompanyStartDate:   cmpycp[0].Companystartdate.Time.String(),
					CompanyCountry:     cmpycp[0].Companycountry.String,
					CompanyCity:        cmpycp[0].Companycity.String,
					CompanyState:       cmpycp[0].Companystate.String,
					CompanyPinCode:     cmpycp[0].Companypincode.String,
					CompanyPhone:       cmpycp[0].Companyphone.String,
					CompanyFax:         cmpycp[0].Companyfax.String,
					CompanyMobile:      cmpycp[0].Companymobile.String,
					CompanyEmail:       cmpycp[0].Companyemail.String,
					CompanyWebsite:     cmpycp[0].Companywebsite.String,
					CompanyFiscalYear:  cmpycp[0].Companyfiscalyear.String,
					CompanyTimeZone:    cmpycp[0].Companytimezone.String,
					CompanyBaseCurency: cmpycp[0].Companybasecurency.String,
					CompanysParent:     cmpycp[0].Companysparent.String,
				}

				if cmpy, errs = commonfuncs.Companyupdate(app, w, r, cpy, &cpy1); errs != nil {
					return
				}
				cmpycp = *cmpy

			} else {
				havcpydetail = true
				status = "ERROR"
				//TODO: send error response.
			}

		} else if cpyop.Optype == "save" {
			havcpydetail = false
			status = "SUCCESS"

			if cmpy, errs = commonfuncs.CompanySave(app, w, r, cpy); errs != nil {
				return
			}
			cmpycp = *cmpy

		}
		//	}

		fmt.Println("-------------------\n fetchCompany in save company Start 1 \n-------------------")

		//cmpycp = []models.TblCompany{}
		//havcpydetail = true
		//status = "success"

		fmt.Println("-------------------\n fetchCompany in save company  Start 2 \n-------------------")

		lgmsg := "Company Save successful.  But havecpy detail? = " + strconv.FormatBool(havcpydetail)
		//ssd := map[string]interface{}{"message": lgmsg, "company": cmpycp, "refdata": ddf.RefResult}
		ssd := map[string]interface{}{"message": lgmsg, "company": cmpycp, "refdata": ""}
		cc := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     status,
			SlugCode:   "COMPANY-SAVE",
			LogMsg:     lgmsg,
		}
		cc.HttpRespond()
		*/
		fmt.Println("-------------------\n save Company Stop \n-------------------")
		return

	}
}

func DoBrFetch(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(fetchBranch(app))
}

func DoBrSave(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(saveBranch(app))
}
