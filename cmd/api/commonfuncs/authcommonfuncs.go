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

func CheckUserRegistered(app *application.Application, w http.ResponseWriter, r *http.Request) (bool, error) {
	//Check user registered Start

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	if !ok {
		err := fmt.Errorf("Empty context")
		dd := errors.SlugError{
			Err:        err,
			ErrType:    errors.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": "Technical Error.  Please contact support"},
			SlugCode:   "AUTH-CHKCTX",
			LogMsg:     "Context fetch error",
		}
		dd.HttpRespondWithError()
		return false, err
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
			Data:       map[string]interface{}{"message": "Database error"},
			SlugCode:   "AUTH-INT",
			LogMsg:     "Database error",
		}
		dd.HttpRespondWithError()

		return false, err
	}
	fmt.Println("ddsds")
	fmt.Println(myc)

	if len(myc) > 1 {
		dd := errors.SlugError{
			Err:        err,
			ErrType:    errors.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": "Invalid Company Profile Setup Exists.  Contact Support"},
			SlugCode:   "AUTH-NOMULCPY",
			LogMsg:     "Company Details Not set or Have multiple Company; sql:" + qry,
		}
		dd.HttpRespondWithError()
		return false, fmt.Errorf("Invalid Company Profile Setup Exists.  Contact Support")
	} else if len(myc) == 0 {
		fmt.Println("no record db success")
		return false, nil
	}

	//Check user registered end
	return true, nil
}

func RegisterUser(app *application.Application, w http.ResponseWriter, r *http.Request) (bool, error) {

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)
	if !ok {
		err := fmt.Errorf("Empty context")
		dd := errors.SlugError{
			Err:        err,
			ErrType:    errors.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": "Technical Error.  Please contact support"},
			SlugCode:   "AUTH-REGCTX",
			LogMsg:     "Context fetch error",
		}
		dd.HttpRespondWithError()
		return false, err
	}

	const qry = `INSERT INTO ac.userlogin (userid, username, useremail, userpassword, userstatus, emailverified, siteid, userstatlstupdt, octime, lmtime) 
	VALUES ($1, $2, $3, $4, $5,$6,$7, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);`

	//uspass := ""

	var myc dbtran.Resultset

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("insert", qry, &myc, userinfo.UUID, userinfo.DisplayName, userinfo.Email, "", "A", userinfo.EmailVerified, userinfo.Siteid),
	}
	fmt.Println("calling tran")
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
			Data:       map[string]interface{}{"message": "Database error"},
			SlugCode:   "ERROR",
			LogMsg:     "Database error",
		}
		//	dd.HttpRespondWithError()
		fmt.Println(dd)
		return false, err
	}
	fmt.Println("calling tran end")

	fmt.Println(myc)

	return true, nil
}
