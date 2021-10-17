package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

/*
import (
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func getPacks(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n getPacks Start \n-------------------")
		defer r.Body.Close()

		var data string
		//var myc []models.TblMytree
		fmt.Println(r.Body)
		//isregis := false
		//havdom := false

		//Check user registered END
		userinfo, errs := commonfuncs.FetchUserinfoFromcontext(w, r, "PACKAGE-CHKCTX")
		if errs != nil {
			return
		}
		var myc *[]models.TblMytree
		myc, errs = commonfuncs.PackageFetch(app, w, r)
		if errs != nil {
			return
		}

		data = "Your domain registration successful. Login with your url - " + hostname
		ssd := map[string]interface{}{"message": data}
		//&nat{"nat1", "nat2"},
		fmt.Println("domain registration completed ss sent")
		ss := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     "SUCCESS",
			SlugCode:   "DOMAIN-REG",
			LogMsg:     "testing",
		}

		ss.HttpRespond()
		fmt.Println("-------------------\n getPacks Stop \n-------------------")
		return

	}
}



func DoPacks(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(getPacks(app))
	//return createUser(app)
}
*/

func setPacks(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n setPacks Start \n-------------------")
		defer r.Body.Close()
		ctx := r.Context()
		userinfo, erre := commonfuncs.FetchUserinfoFromcontext(w, r, "COMPANYCHK-CHKCTX")
		if erre != nil {
			return
		}
		var p models.PackSelect
		data := "Subdomain or Domain already in use.  Please select a new value"

		err := json.NewDecoder(r.Body).Decode(&p)
		fmt.Println(p)

		if err != nil {
			return
		}

		fmt.Println(userinfo.Companyid)
		fmt.Println(p.Planid)

		const qry = `INSERT INTO ac.companypacks (COMPANYID,PLANID,PACKFUNCID,STARTDATE,EXPIRYDATE,USERROLELiMIT,USERLIMIT,BRANCHLIMIT,STATUS,OCTIME,LMTIME)
						SELECT $1::varchar AS COMPANYID,PLANID,packid,CURRENT_DATE,CURRENT_DATE+  make_interval(days => durationdays) ,userrolelimit,userlimit,branchlimit,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP 
						FROM ac.planpacks WHERE planid = $2
						AND NOT EXISTS (SELECT 1 FROM ac.companypacks WHERE companyid= $1 AND   planid = $2 )`

		var myc dbtran.Resultset

		stmts := []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("insert", qry, &myc, userinfo.Companyid, p.Planid),
		}

		_, errs := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
			err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
			return err
		})

		fmt.Println("myrows", myc.RowsAffected)

		if errs != nil || myc.RowsAffected < 1 {
			//https://github.com/jackc/pgx/issues/474

			data = "Technical Error.  Please contact support"
			lgmsg := "No rows updated for set plan action"
			var pgErr *pgconn.PgError
			if errs != nil {

				if errors.As(err, &pgErr) {
					lgmsg = pgErr.Error()
				}
			}
			fmt.Println(data)
			//		dd := errors.SlugError{
			dd := httpresponse.SlugResponse{
				Err:        err,
				ErrType:    httpresponse.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Data:       map[string]interface{}{"message": data},
				SlugCode:   "PACKREGIS-INSERT",
				LogMsg:     lgmsg,
			}
			//dd.HttpRespondWithError()
			dd.HttpRespond()
			return
		}

		data = "Plan registered successfully"
		dd := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "PACKREGIS-INSERT",
			LogMsg:     "setPacks End SUCCESSFULLY",
		}
		//dd.HttpRespondWithError()
		dd.HttpRespond()

		fmt.Println("-------------------\n setPacks End \n-------------------")
		return
	}
}

func DoPacks(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(setPacks(app))
	//return createUser(app)
}
