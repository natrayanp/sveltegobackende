package mymiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/errors"
	"github.com/sveltegobackend/pkg/errors/httperr"
	"github.com/sveltegobackend/pkg/fireauth"
)

func ParseHeadMiddleware(app *application.Application) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			sess := r.Header.Get("session")
			sid := r.Header.Get("siteid")
			fmt.Println(sess)
			fmt.Println(sid)

			userinfo, ctxfetchok := ctx.Value(fireauth.UserContextKey).(fireauth.User)
			if !ctxfetchok {
				dd := errors.SlugError{
					ErrType:    errors.ErrorTypeDatabase,
					RespWriter: w,
					Request:    r,
					Data:       map[string]interface{}{"message": "Technical issue. Please contact support"},
					SlugCode:   "PARSEHEADER-CTXFETCHFAIL",
					LogMsg:     "Context fetch Failed",
				}
				dd.HttpRespondWithError()
			}

			userinfo.Session = sess
			userinfo.Siteid = sid

			const qry = `SELECT companyid FROM ac.domainmap WHERE
				domainmapid = (SELECT domainmapid FROM ac.userlogin WHERE userid = $1 AND siteid = $2)
				AND status = 'A'`

			type Cid struct {
				Companyid string
			}
			var myc []Cid

			stmts := []*dbtran.PipelineStmt{
				dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Siteid),
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
				httperr.InternalError("Database error", "", err, w, r)
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
			fmt.Println(userinfo)
			userinfo.Companyid = ""
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
