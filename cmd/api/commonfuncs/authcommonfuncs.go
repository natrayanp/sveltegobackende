package commonfuncs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/errors"
	"github.com/sveltegobackend/pkg/fireauth"
)

func CheckUserRegistered(app *application.Application, w http.ResponseWriter, r *http.Request) bool {
	//Check user registered Start

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)
	if !ok {
		return ok
	}

	qry := `SELECT * FROM ac.userlogin
			WHERE userid = $1
			AND  siteid = $2;`

	var myc []models.TblUserlogin

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Siteid),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeNoTran, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
		return err
	})

	if err != nil {
		dd := errors.SlugError{
			Err:        err,
			ErrType:    errors.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Slug:       "Database error",
			SlugCode:   "AUTH-INT",
			LogMsg:     "Database error",
		}
		dd.HttpRespondWithError()
		return false
	}
	fmt.Println("ddsds")
	fmt.Println(myc)

	if len(myc) > 1 {
		dd := errors.SlugError{
			Err:        err,
			ErrType:    errors.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Slug:       "Invalid Company Profile Setup Exists.  Contact Support",
			SlugCode:   "AUTH-NOMULCPY",
			LogMsg:     "Company Details Not set or Have multiple Company; sql:" + qry,
		}
		dd.HttpRespondWithError()
		return false
	} else if len(myc) == 0 {
		fmt.Println("no record db success")
		return false
	}

	//Check user registered end
	return true
}

func RegisterUser(app *application.Application, w http.ResponseWriter, r *http.Request) bool {

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)
	if !ok {
		return ok
	}

	const qry = `INSERT INTO ac.userlogin (userid, username, useremail, userpassword, userstatus, emailverified, siteid, userstatlstupdt, octime, lmtime) 
	VALUES ($1, $2, $3, $4, $5,$6,$7, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);`

	//uspass := ""

	var myc []dbtran.Resultset

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.DisplayName, userinfo.Email, "", "A", userinfo.EmailVerified, userinfo.Siteid),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeNoTran, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
		return err
	})

	if err != nil {
		dd := errors.SlugError{
			Err:        err,
			ErrType:    errors.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Slug:       "Database error",
			SlugCode:   "AUTH-INT",
			LogMsg:     "Database error",
		}
		dd.HttpRespondWithError()
		return false
	}

	fmt.Println(myc)

	return true
}
