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
		cpy := &models.Cpy{}
		fmt.Println(p)
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

		/*

			if cmpy, errs = commonfuncs.CompanyCheck(app, w, r); errs != nil {
				return
			}
			cmpycp = *cmpy

			if len(*cmpy) == 1 {
				havcpydetail = true
				status = "FAILURE"
			} else {*/
		havcpydetail = false
		status = "SUCCESS"
		cmpycp = []models.TblCompany{}

		if cmpy, errs = commonfuncs.CompanySave(app, w, r, cpy); errs != nil {
			return
		}
		cmpycp = *cmpy

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
