package commonfuncs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vgarvardt/gue/v2"

	//"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func DomRegis(app *application.Application, w http.ResponseWriter, r *http.Request, dom string) error {

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	if !ok {
		err := fmt.Errorf("Empty context")
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": "Technical Error.  Please contact support"},
			SlugCode:   "SESSION-CHKCTX",
			LogMsg:     "Context fetch error",
		}
		dd.HttpRespond()
		return err
	}

	const qry1 = `UPDATE ac.userlogin SET companyid = 'CPYID'||nextval('companyid_seq'::regclass), lmtime = CURRENT_TIMESTAMP, hostname = $1
					WHERE userid = $2`

	var myc1 dbtran.Resultset

	stmts1 := []*dbtran.PipelineStmt{

		dbtran.NewPipelineStmt("update", qry1, &myc1, dom, userinfo.UUID),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
		return err
	})

	if err != nil {
		return err
	}
	type assingrole struct {
		UUID string
	}
	args, err := json.Marshal(assingrole{UUID: userinfo.UUID})

	if err != nil {
		return err
	}

	err = app.Que.Enquejob(&gue.Job{Type: "AssignRole", Args: args})

	if err != nil {
		return err
	}
	return nil

}
