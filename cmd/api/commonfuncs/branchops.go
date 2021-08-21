package commonfuncs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/r3labs/diff/v2"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func BranchCheck(app *application.Application, w http.ResponseWriter, r *http.Request, companyid string, branchidlist []string) (*[]models.TblBranch, error) {
	fmt.Println("----------------- PACKAGE CHECK START -------------------")

	var data string
	var myc []models.TblBranch

	ctx := r.Context()

	userinfo, errs := FetchUserinfoFromcontext(w, r, "BRANCHCHK-CHKCTX")
	if errs != nil {
		return &[]models.TblBranch{}, errs
	}

	if companyid == "DEFAULT" {
		companyid = userinfo.Companyid
	}

	qry := `SELECT a.*,b.companyname FROM ac.branch a
						FULL OUTER JOIN ac.company b
						ON a.companyid = b.companyid
						WHERE a.companyid = $1
						AND b.companystatus in ('A')
						AND a.branchStatus in ('A') `

	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, companyid),
	}

	if branchidlist[0] != "all" {
		qry = qry + "AND a.branchid = ANY($2) "

		stmts1 = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &myc, companyid, branchidlist),
		}
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
		return err
	})

	if err != nil {
		//https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError

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
			SlugCode:   "BRANCHCHK-SELECT",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.TblBranch{}, err
	}

	fmt.Println("----------------- PACKAGE CHECK END -------------------")
	fmt.Println(myc)
	return &myc, nil
}

func BranchSave(app *application.Application, w http.ResponseWriter, r *http.Request, cpy *models.TblBranch) (*[]models.TblBranch, error) {
	fmt.Println("----------------- Branch Save CHECK START -------------------")

	var data string
	var lgmsg string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)
	userinfo, errs := FetchUserinfoFromcontext(w, r, "BRANCHCHK-CHKCTX")
	fmt.Println(userinfo)
	fmt.Println(userinfo.Companyid)

	if errs != nil {
		return &[]models.TblBranch{}, errs
	}

	const qry = `INSERT INTO ac.branch VALUES
						($1,DEFAULT,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,
							CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)  RETURNING *;`

	var myc []models.TblBranch
	//var myc dbtran.Resultset
	fmt.Println(cpy)
	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, cpy.Companyid, cpy.Branchname, cpy.Branchshortname, "ca",
			"A", "desc", "imgur", cpy.Branchaddline1, cpy.Branchaddline2, cpy.Branchcity, cpy.Branchstate, cpy.Branchcountry, cpy.Branchpincode,
			cpy.Branchphone, cpy.Branchfax, cpy.Branchmobile, cpy.Branchwebsite, cpy.Branchemail, cpy.Branchstartdate, cpy.Isdefault, userinfo.UUID),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
		return err
	})
	fmt.Println(myc)
	//if err != nil || myc.RowsAffected < 1 {
	if err != nil || len(myc) < 1 {
		//https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		data = "Technical Error.  Please contact support"

		if errors.As(err, &pgErr) {
			lgmsg = pgErr.Error()
		}

		//	if myc.RowsAffected < 1 {
		if len(myc) < 1 {
			err = errors.New("no rows inserted")
			lgmsg = "No data saved by successful INSERT query"
		}

		//		dd := errors.SlugError{
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "BRANCH-SAVE",
			LogMsg:     lgmsg,
		}
		dd.HttpRespond()
		return &[]models.TblBranch{}, err
	} else {
		fmt.Println(myc)
	}

	fmt.Println("-----------------  Branch Save CHECK END -------------------")

	return &myc, nil
}

func Branchupdate(app *application.Application, w http.ResponseWriter, r *http.Request, cpynew *models.TblBranch, cpyindb *models.TblBranch) (*[]models.TblBranch, error) {
	fmt.Println("----------------- Branch Update CHECK START -------------------")

	var data string
	var lgmsg string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "BRANCHCHK-CHKCTX")
	if errs != nil {
		return &[]models.TblBranch{}, errs
	}

	changelog, erre := diff.Diff(cpyindb, cpynew)
	if erre != nil {
		return &[]models.TblBranch{}, errs
	}

	qry := "UPDATE ac.branch SET "
	if len(changelog) > 0 {
		for i, s := range changelog {
			if i != 0 {
				qry = qry + " , "
			}
			qry = qry + s.Path[0] + ` =  '` + fmt.Sprintf("%v", s.To) + `' `
		}

		qry = qry + ", lmtime = CURRENT_TIMESTAMP, lmuserid = $1 "

		qry = qry + "WHERE companyid = $2 AND branchid = $3 RETURNING *;"
	}

	fmt.Println(qry)
	var myc []models.TblBranch
	//var myc dbtran.Resultset

	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, cpynew.Companyid, cpynew.Branchid),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
		return err
	})

	//if err != nil || myc.RowsAffected < 1 {
	if err != nil || len(myc) < 1 {
		//https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		data = "Technical Error.  Please contact support"

		if errors.As(err, &pgErr) {
			lgmsg = pgErr.Error()
		}

		//if myc.RowsAffected < 1 {
		if len(myc) < 1 {
			err = errors.New("no data updated")
			lgmsg = "No data updated by succesful UPDATE query"
		}

		//		dd := errors.SlugError{
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "BRANCH-UPDATE",
			LogMsg:     lgmsg,
		}
		dd.HttpRespond()
		return &[]models.TblBranch{}, err
	} else {
		fmt.Println(myc)
	}

	fmt.Println("----------------- Branch Update CHECK END -------------------")
	return &myc, nil
	//return &myc, nil
}
