package mymiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"

	//"github.com/sveltegobackend/pkg/errors/httperr"
	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func ParseHeadMiddleware(app *application.Application) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			org := r.Header.Get("Origin")
			sess := r.Header.Get("session")
			sid := r.Header.Get("siteid")

			fmt.Println("============= chk session ===================")
			fmt.Println(sess)
			fmt.Println(sid)

			userinfo, ctxfetchok := ctx.Value(fireauth.UserContextKey).(fireauth.User)
			if !ctxfetchok {
				dd := httpresponse.SlugResponse{
					ErrType:    httpresponse.ErrorTypeDatabase,
					RespWriter: w,
					Request:    r,
					Data:       map[string]interface{}{"message": "Technical issue. Please contact support"},
					SlugCode:   "PARSEHEADER-CTXFETCHFAIL",
					LogMsg:     "Context fetch Failed",
				}
				//dd.HttpRespondWithError()
				dd.HttpRespond()
				return
			}
			dd := strings.SplitAfter(org, "://")
			//ss:=strings.Split(dd[1], ".")
			fmt.Println(dd[1])

			userinfo.Hostname = dd[1]
			userinfo.Session = sess
			userinfo.Siteid = sid

			const qry = `SELECT * FROM ac.userlogin WHERE				
				userid = $1
				AND siteid = $2
				AND hostname = $3
				AND userstatus = 'A'`

			var myc []models.TblUserlogin

			stmts := []*dbtran.PipelineStmt{
				dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Siteid, userinfo.Hostname),
				//dbtran.NewPipelineStmt("select", qry, &myc),
				//dbtran.NewPipelineStmt("delete", qry, nil),
			}

			_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeNoTran, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
				err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
				return err
			})

			if err != nil {
				fmt.Println(err)
				//TODO: implement error response
				//httperr.InternalError("Database error", "", err, w, r)
				dd := httpresponse.SlugResponse{
					ErrType:    httpresponse.ErrorTypeDatabase,
					RespWriter: w,
					Request:    r,
					Data:       map[string]interface{}{"message": "Technical issue. Please contact support"},
					SlugCode:   "PARSEHEADER-DOMAINFETCH",
					LogMsg:     "Domain fetch Failed",
				}
				//dd.HttpRespondWithError()
				dd.HttpRespond()
				return
			}
			/*
				fmt.Println("+++++++++++++++++++++res")
				fmt.Println(myc)
				fmt.Println(len(myc))
				fmt.Println("+++++++++++++++++++++res")

				fmt.Println("+++++++++++++++++++++res non select")
				fmt.Println(stmts)
				fmt.Println((*stmts[0]).Resultstruct)
				fmt.Println("+++++++++++++++++++++res non select")
			*/

			if len(myc) == 1 {
				userinfo.Companyid = *myc[0].Companyid
			} else {
				userinfo.Companyid = "PUBLIC"
			}

			userinfo.Entityid = myc[0].Entityid
			//myc[0].Entityid.AssignTo(userinfo.Entityid)

			errs := commonfuncs.SessionOps(app, w, r, &userinfo)
			if errs != nil {
				return
			}

			fmt.Println(userinfo)
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@userinfo@@@@@@@@@@@@@@@@@@")
			ctx = context.WithValue(ctx, fireauth.GetUserCtxKey(), userinfo)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

/*
func ParseHeadMiddleware1(next http.Handler, app application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sess := r.Header.Get("session")
		sid := r.Header.Get("siteid")

		userinfo, err := fireauth.UserFromCtxs(ctx)
		userinfo.Session = sess
		userinfo.Siteid = sid

		const qry = `SELECT companyid FROM ac.domainmap WHERE
		domainmapid = (SELECT domainmapid FROM ac.userlogin WHERE userid = $1 AND siteid = $2)
		AND status = 'A'`

		type cid struct {
			companyid string
		}
		var myc []map[string]interface{}

		stmts := []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, myc, userinfo.UUID, userinfo.Siteid),
		}

		_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeNoTran, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
			err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
			return err
		})

		if err != nil {
			fmt.Println(err)
			//TODO: implement error response
		}

		fmt.Println(myc)

		if val, ok := myc[0]["companyid"]; ok {
			//do something here
			userinfo.Companyid = fmt.Sprintf("%v", val)
		}

		fireauth.SetUserInCtx(userinfo, r)
		next.ServeHTTP(w, r)
	})

}
*/
