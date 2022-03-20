package company

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		var cpyd string

		defer r.Body.Close()
		var p models.ReqEntityTree

		err := json.NewDecoder(r.Body).Decode(&p)
		fmt.Println(p)

		if err != nil {
			return
		}

		if p.EntityType == "company" && p.Entityid[0] != "null" {
			cpyd = p.Entityid[0]
		} else {
			cpyd = "DEFAULT"
		}

		if cmpy, errs = commonfuncs.CompanyCheck(app, w, r, cpyd); errs != nil {
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

		ddf := models.RefDatReqFinal{
			Refs: []models.RefDatReq{
				{Reftype: "group", Refname: "company"},
			},
		}

		fmt.Println("-------------------\n fetchCompany Start 1  comp\n-------------------")

		if err := commonfuncs.RefDataFetch1(app, w, r, &ddf); err != nil {
			return
		}

		fmt.Println(ddf.RefResult)

		fmt.Println("-------------------\n fetchCompany Start 2 comp \n-------------------")

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

func saveCompany(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n save Company Start \n-------------------")

		r.Body = http.MaxBytesReader(w, r.Body, 32<<20+1024)
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// parse text field
		//	text := make([]byte, 512)
		p, err := reader.NextPart()
		// one more field to parse, EOF is considered as failure here
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if p.FormName() != "text_field" {
			http.Error(w, "text_field is expected", http.StatusBadRequest)
			return
		}
		//cpy := &models.Cpy{}
		cpy := &models.TblCompany{}
		jsonDecoder := json.NewDecoder(p)
		err = jsonDecoder.Decode(&cpy)

		//_, err = p.Read(text)

		if err != nil {
			//&& err != io.EOF {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		/*

			out, _ := json.Marshal(text)
			out1, _ := json.Marshal(out)

			cpy := &models.Cpy{}
			err = json.Unmarshal([]byte(text), &cpy)

			fmt.Printf("%s", out1)
			fmt.Printf(err.Error())
		*/
		fmt.Println("----------")
		fmt.Println(cpy)

		p, err = reader.NextPart()
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if p.FormName() != "text_action" {
			http.Error(w, "text_action is expected", http.StatusBadRequest)
			return
		}
		cpyop := &models.Cpyops{}
		fmt.Println(p)
		jsonDecoder = json.NewDecoder(p)
		err = jsonDecoder.Decode(&cpyop)

		//_, err = p.Read(text)

		if err != nil {
			//&& err != io.EOF {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("----------cpyop\n----------")
		fmt.Println(cpyop)
		fmt.Println("----------cpyop\n----------")

		// parse file field
		p, err = reader.NextPart()
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if p.FormName() != "file_field" {
			http.Error(w, "file_field is expected", http.StatusBadRequest)
			return
		}
		buf := bufio.NewReader(p)
		sniff, _ := buf.Peek(512)
		contentType := http.DetectContentType(sniff)
		fmt.Println(contentType)

		/*
			if contentType != "application/zip" {
				http.Error(w, "file type not allowed", http.StatusBadRequest)
				return
			}
		*/
		fmt.Println("cleared this")

		f, err := ioutil.TempFile("", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var maxSize int64 = 32 << 20
		lmt := io.MultiReader(buf, io.LimitReader(p, maxSize-511))
		written, err := io.Copy(f, lmt)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if written > maxSize {
			os.Remove(f.Name())
			http.Error(w, "file size over limit", http.StatusBadRequest)
			return
		}

		fmt.Println("written - :", written, "maxSize - ", maxSize)
		var cmpy *[]models.TblCompany
		var cmpycp []models.TblCompany
		havcpydetail := false
		var errs error
		var status string
		cmpycp = []models.TblCompany{}

		if cpyop.Optype == "Update" {
			if cmpy, errs = commonfuncs.CompanyCheck(app, w, r, cpy.Companyid); errs != nil {
				return
			}
			cmpycp = *cmpy
			//layout := "2006-01-02T15:04:05.000Z"
			/*
				const layoutISO = "2006-01-02"
				cpydd, _ := time.Parse(layoutISO, cpy.CompanyStartDate)
				cpy.CompanyStartDate = cpydd.String()


				fmt.Println(cmpycp[0].Companystartdate)
			*/
			if len(*cmpy) == 1 {
				havcpydetail = true
				status = "SUCCESS"

				cpy1 := cmpycp[0]

				/*
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
				*/
				if cmpy, errs = commonfuncs.Companyupdate(app, w, r, cpy, &cpy1); errs != nil {
					return
				}
				cmpycp = *cmpy

			} else {
				havcpydetail = true
				status = "ERROR"
				//TODO: send error response.
			}

		} else if cpyop.Optype == "Save" {
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
		fmt.Println("-------------------\n save Company Stop \n-------------------")
		return

	}
}

func DoFetch(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(fetchCompany(app))
}

func DoSave(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(saveCompany(app))
}
