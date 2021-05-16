package middleware

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/fireauth"
)

func ParseHeadMiddleware(next http.Handler, app application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sess := r.Header.Get("session")
		sid := r.Header.Get("siteid")

		userinfo, err := fireauth.UserFromCtxs(ctx)
		userinfo.session = sess
		userinfo.siteid = sid

		const qry = `SELECT companyid FROM ac.domainmap WHERE 
		domainmapid = (SELECT domainmapid FROM ac.userlogin WHERE userid = $1 AND siteid = $2)
		AND status = 'A'`

		type cid struct {
			companyid string
		}
		var myc cid

		stmts := []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.siteid),
		}

		_, err = dbtran.WithTransaction(dbtran.TranTypeNoTran, app.DB, nil, func(typ dbtran.TranType, db *sqlx.DB, ttx dbtran.Transaction) error {
			_, err := dbtran.RunPipeline(typ, db, ttx, stmts...)
			return err
		})

		if err != nil {
			fmt.Println(err)
			//TODO: implement error response
		}

		fmt.Println(myc)
		userinfo.companyid = myc.companyid
		fireauth.SetUserInCtx(userinfo, r)
		next.ServeHTTP(w, r)
	})

}
