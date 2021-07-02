package commonfuncs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func refCountry(app *application.Application, w http.ResponseWriter, r *http.Request, packfuncid []string) (*[]models.TblRefdata, error) {
	fmt.Println("----------------- refCountry Fetch START -------------------")

	var data string
	packfuncidf := ""
	if len(packfuncid) == 1 && packfuncid[0] == "ALL" {
		packfuncidf = "ALL"
	} else {
		for _, n := range packfuncid {
			packfuncidf = packfuncidf + n
		}
	}

	ctx := r.Context()

	var qry string
	var myc []models.TblRefdata
	var stmts []*dbtran.PipelineStmt

	qry = `WITH RECURSIVE MyTree AS (
			SELECT id,refvalcat,refvalue,parent FROM ac.refdata WHERE refcode = 'country'
UNION
SELECT m.id,m.refvalcat,m.refvalue,m.parent as open FROM ac.refdata AS m JOIN MyTree AS t ON t.id  = ANY(m.parent)
)
SELECT * FROM MyTree ORDER BY refvalcat,refvalue;`

	stmts = []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
		return err
	})

	fmt.Println("---------------$$$end5")
	//fmt.Println(err.Error())
	if err != nil {
		//https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		fmt.Println("---------------$$$end5a")
		if errors.As(err, &pgErr) {
			data = "Technical Error.  Please contact support"
		}

		//		dd := errors.SlugError{
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "DOMAINREG-UPDATE",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.TblRefdata{}, err
	}

	fmt.Println("---------------$$$end6")
	dd, _ := json.Marshal(myc)
	fmt.Println(string(dd))
	fmt.Println("---------------$$$end6a")
	fmt.Printf("&myc is: %p\n", &myc)
	createDataTree1(&myc)
	fmt.Println("---------------$$$end6b")
	dd1, _ := json.Marshal(myc)
	fmt.Printf("&myc is: %p\n", &myc)
	fmt.Println(string(dd1))
	fmt.Println("---------------$$$end7")
	fmt.Println("----------------- refCountry FETCH END -------------------")

	return &myc, nil
}

func createDataTree1(mnodes *[]models.TblRefdata) {
	nodes := *mnodes
	fmt.Printf("&mnodes is: %p\n", mnodes)
	var newnodes []models.TblRefdata
	fmt.Println("---------------$$$end6a1")
	m := make(map[pgtype.Text]*models.TblRefdata)
	for i, _ := range nodes {
		//fmt.Printf("Setting m[%d] = <node with ID=%d>\n", n.ID, n.ID)
		m[nodes[i].Id] = &nodes[i]
	}
	fmt.Println(m)
	fmt.Println("---------------$$$end6a2")
	for i, n := range nodes {
		//fmt.Printf("Setting <node with ID=%d>.Child to <node with ID=%d>\n", n.ID, m[n.ParentID].ID)
		fmt.Println(n)
		fmt.Println(n.Parent.Dimensions[0].Length)
		fmt.Println("---------------$$$end6a2a")
		if n.Parent.Dimensions[0].Length > 0 {
			fmt.Println(n.Parent.Status)

			for _, t := range n.Parent.Elements {
				fmt.Println(t.Status)
				fmt.Println(t.Status == pgtype.Null)
				fmt.Println("---------------$$$end6a2a1")
				if t.Status != pgtype.Null {
					m[t].Submenu = append(m[t].Submenu, m[nodes[i].Id])
				}
			}
		}
	}
	fmt.Println("---------------$$$end6a3")
	for _, n := range m {
		fmt.Println(n)
		fmt.Println(n.Parent.Elements[0].Status)
		fmt.Println(n.Parent.Elements[0].Status == pgtype.Null)
		if n.Parent.Elements[0].Status == pgtype.Null {
			fmt.Println(n)
			fmt.Println(newnodes)
			newnodes = append(newnodes, *n)
			fmt.Println(newnodes)
		}
	}
	fmt.Println("---------------$$$end6a4")
	fmt.Printf("&mnodes is: %p\n", mnodes)
	fmt.Printf("&newnodes is: %p\n", &newnodes)
	*mnodes = newnodes
	fmt.Printf("&mnodes is: %p\n", mnodes)
	fmt.Println("---------------$$$end6a5")
}

func refIndusType(app *application.Application, w http.ResponseWriter, r *http.Request) (*[]models.TblRefdata, error) {
	fmt.Println("----------------- refCountry Fetch START -------------------")

	var data string

	ctx := r.Context()

	var qry string
	var myc []models.TblRefdata
	var stmts []*dbtran.PipelineStmt

	qry = `SELECT id,refvalcat,refvalue,parent FROM ac.refdata WHERE refcode = 'industype' ORDER BY refvalcat,refvalue;`

	stmts = []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
		return err
	})

	fmt.Println("---------------$$$end5")
	//fmt.Println(err.Error())
	if err != nil {
		//https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		fmt.Println("---------------$$$end5a")
		if errors.As(err, &pgErr) {
			data = "Technical Error.  Please contact support"
		}

		//		dd := errors.SlugError{
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "DOMAINREG-UPDATE",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.TblRefdata{}, err
	}

	return &myc, nil
}

func RefDataFetch1(app *application.Application, w http.ResponseWriter, r *http.Request, refdata *models.RefDatReqFinal) error {
	var rf []string
	m := make(map[string]interface{})
	for i, s := range *&refdata.Refs {
		fmt.Println("inside range")
		fmt.Println(s.Refname)
		fmt.Println(s.Reftype)
		fmt.Println(i)
		if s.Reftype == "group" {
			rf = getRefIndividualItems(s.Refname)
		} else {
			rf = append(rf, s.Refname)
		}

		for _, t := range rf {
			switch t {
			case "country":
				//dd := []string{"india", "singapore"}
				d := []string{"ALL"}
				ss, e := refCountry(app, w, r, d)
				if e != nil {
					return e
				}
				fmt.Println(ss)
				m["country"] = ss
				fmt.Println(m)

			case "industype":
				ss, _ := refIndusType(app, w, r)
				m["industype"] = ss
			}
		}

	}

	*&(*&refdata).RefResult = m

	fmt.Println(*&refdata.RefResult)
	return nil
}

func getRefIndividualItems(group string) []string {

	switch group {
	case "company":
		return []string{"country", "industype"}
	default:
		return []string{}
	}

}

/*

Example to accept any data

package main

import (
	"fmt"
		"encoding/json"

)

type  EmailPostData struct{
Name string
Body string
}


func main() {
	fmt.Println("Hello, playground")
	    var postData = EmailPostData{}
	    ConvertRequestJsonToJson( &postData)
	     fmt.Println(postData)

}


func ConvertRequestJsonToJson( model interface{}) {
	postContent := []byte(`{"Name":"Alice","Body":"Hello"}`)

   	 json.Unmarshal(postContent, model)
	fmt.Println("completed")
	//json.Unmarshal stores the result in the value pointed to by model
}


*/
