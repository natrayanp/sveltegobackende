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

func CompanyCheck(app *application.Application, w http.ResponseWriter, r *http.Request) (*[]models.TblCompany, error) {
	fmt.Println("----------------- PACKAGE CHECK START -------------------")

	var data string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "COMPANYCHK-CHKCTX")
	if errs != nil {
		return &[]models.TblCompany{}, errs
	}

	const qry = `SELECT * FROM ac.company 
					WHERE companyid = $1
					AND companyStatus in ('A')`

	var myc []models.TblCompany

	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.Companyid),
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
			SlugCode:   "COMPANYCHK-SELECT",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.TblCompany{}, err
	}

	fmt.Println("----------------- PACKAGE CHECK END -------------------")

	return &myc, nil
}

func CompanySave(app *application.Application, w http.ResponseWriter, r *http.Request, cpy *models.Cpy) (*[]models.TblCompany, error) {
	fmt.Println("----------------- CompanySave CHECK START -------------------")

	var data string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "COMPANYCHK-CHKCTX")
	if errs != nil {
		return &[]models.TblCompany{}, errs
	}

	const qry = `INSERT INTO ac.company VALUES
					($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,
						$25,$26,$27,$28,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) RETURNING *;`

	//fmt.Println(cpy)

	var myc []models.TblCompany
	//var myc dbtran.Resultset

	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.Companyid, cpy.CompanyName, cpy.CompanyShortName, cpy.CompanyCategory,
			"A", "", "", "", cpy.CompanyIndustry, cpy.CompanyTaxID, cpy.CompanyAddLine1, cpy.CompanyAddLine2,
			cpy.CompanyCity, cpy.CompanyState, cpy.CompanyCountry, cpy.CompanyPinCode, cpy.CompanyPhone, cpy.CompanyFax,
			cpy.CompanyMobile, cpy.CompanyWebsite, cpy.CompanyEmail, cpy.CompanyStartDate, cpy.CompanyFiscalYear, cpy.CompanyTimeZone,
			cpy.CompanyBaseCurency, cpy.CompanysParent, "Y", userinfo.UUID),
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
			SlugCode:   "COMPANY-SAVE",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.TblCompany{}, err
	} else {
		fmt.Println(myc)
	}

	fmt.Println("----------------- CompanySave CHECK END -------------------")

	return &myc, nil
}

func Companyupdate(app *application.Application, w http.ResponseWriter, r *http.Request, cpynew *models.Cpy, cpyindb *models.Cpy) (*[]models.TblCompany, error) {
	fmt.Println("----------------- Company Update CHECK START -------------------")

	//var data string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "COMPANYCHK-CHKCTX")
	if errs != nil {
		return &[]models.TblCompany{}, errs
	}

	changelog, erre := diff.Diff(cpyindb, cpynew)
	if erre != nil {
		return &[]models.TblCompany{}, errs
	}

	qry := "UPDATE ac.company SET "
	if len(changelog) > 0 {
		for i, s := range changelog {
			if i != 0 {
				qry = qry + " , "
			}
			qry = qry + s.Path[0] + ` =  '` + fmt.Sprintf("%v", s.To) + `' `
		}

		qry = qry + ", lmtime = CURRENT_TIMESTAMP, lmuserid = $1 "

		qry = qry + "WHERE companyid = $2 RETURNING *;"
	}

	fmt.Println(qry)
	var myc []models.TblCompany
	//var myc dbtran.Resultset

	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, cpynew.CompanyId),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
		return err
	})
	data := "d"

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
			SlugCode:   "COMPANY-UPDATE",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.TblCompany{}, err
	} else {
		fmt.Println(myc)
	}

	fmt.Println("----------------- Company Update CHECK END -------------------")
	return &myc, nil
	//return &myc, nil
}
