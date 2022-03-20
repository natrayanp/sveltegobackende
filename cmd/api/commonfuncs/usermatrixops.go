package commonfuncs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
)

// PackageFetch returns menus and access rights for it for a user.
// Parameters:
// packfuncid --> If you want only the packsfuncs sent. Send the PACKID form AC.PACKS as array.
//					This forcefully sent only those packs if it exists at company and user level(refer query)
// companyid  --> Always send company id here
//					Front don't send company id send userinfo.Companyid from calling side
//					else send whatever company id received from front end
// It returns PacksResp struct which is self explanatory and error.
func UserMatrixFetch(app *application.Application, w http.ResponseWriter, r *http.Request, rolereq *models.UserMatrixReq) (*[]models.Usermatrix, error) {

	var data string
	var err error
	var qry string
	moreusr := false
	var stmts []*dbtran.PipelineStmt

	ctx := r.Context()

	useridf := ""
	if len(rolereq.Matrixuserid) > 0 {
		moreusr = true
		for _, n := range rolereq.Matrixuserid {
			useridf = useridf + n
		}
	}

	userinfo, errs := FetchUserinfoFromcontext(w, r, "PACKAGE-CHKCTX")
	if errs != nil {
		return &[]models.Usermatrix{}, errs
	}
	var matrixlist []models.Usermatrix

	if rolereq.Optype == "detail" {

		qry = `SELECT sd.usrprof_userid AS userid,sd.usrprof_firstname AS firstname,sd.usrprof_lastname AS lastname,sd.usrprof_department AS department,
					sd.usrprof_designation AS designation,sd.usrprof_gender AS gender, sd.usrprof_AddLine1 AS addressline1, sd.usrprof_AddLine2 AS addressline2,
					sd.usrprof_city AS city, usrprof_state  AS state,usrprof_country  AS country, usrprof_pinCode AS  pincode,sd.usrprof_dob AS dob, 
					sd.usrprof_mobile AS mobile,sd.usrprof_email AS email, sd.usrprof_joiningdate AS joiningdate, sd.usrprof_lastdate AS lastdate,
					sd.usrprof_taxid AS taxid, sd.usrprof_companyid AS companyid,sd.usrprof_imagelink AS imagelink, c.userstatus, 'TRUE' AS fulldetails,
				(SELECT json_agg(b)
					FROM (	SELECT a.usrrole_branchidaccess AS branchid,d.branchname,
								(SELECT json_agg(g)
									FROM ( 	SELECT f.rolemasterid,f.rmdisplayname FROM ac.rolemaster f where
							  				f.rolemasterid = ANY(a.usrrole_rolemasterid) AND f.rmstatus = 'A'								
										  )AS g
								) AS roleaccess						 
							FROM ac.userrole a
							LEFT JOIN ac.branch d ON d.companyid = a.usrrole_companyid AND d.branchid = a.usrrole_branchidaccess AND d.branchstatus = 'A'						
							WHERE  a.usrrole_companyid = $1 AND a.usrrole_status = 'A' AND branchname IS NOT NULL`
		if moreusr {
			qry = qry + " AND a.usrrole_userid in ($2) "
		}

		qry = qry + ` ) AS b
		 		)  AS accessmatrix
			FROM ac.userprofile sd 
			LEFT JOIN ac.userlogin c ON sd.usrprof_userid = c.userid AND sd.usrprof_companyid = c.companyid AND c.userstatus = 'A'			
			WHERE sd.usrprof_companyid = $1 AND userstatus IS NOT NULL`

	} else {

		qry = `SELECT sd.usrprof_userid AS userid,sd.usrprof_firstname AS firstname,sd.usrprof_lastname AS lastname,sd.usrprof_department AS department,
					sd.usrprof_designation AS designation,sd.usrprof_gender AS gender, sd.usrprof_AddLine1 AS addressline1, sd.usrprof_AddLine2 AS addressline2,
					sd.usrprof_city AS city, usrprof_state  AS state,usrprof_country  AS country, usrprof_pinCode AS  pincode,sd.usrprof_dob AS dob, 
					sd.usrprof_mobile AS mobile,sd.usrprof_email AS email, sd.usrprof_joiningdate AS joiningdate, sd.usrprof_lastdate AS lastdate,
					sd.usrprof_taxid AS taxid, sd.usrprof_companyid AS companyid,sd.usrprof_imagelink AS imagelink, c.userstatus, 'FALSE' AS fulldetails,
					(SELECT array[]::text[] AS accessmatrix)
			FROM ac.userprofile sd 
			LEFT JOIN ac.userlogin c ON sd.usrprof_userid = c.userid AND sd.usrprof_companyid = c.companyid AND c.userstatus = 'A'			
			WHERE sd.usrprof_companyid = $1 AND userstatus IS NOT NULL`

	}

	if moreusr {
		stmts = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &matrixlist, userinfo.Companyid, useridf),
		}
	} else {
		stmts = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &matrixlist, userinfo.Companyid),
		}
	}

	_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
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
			SlugCode:   "USRMATRIX-FETCH",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &[]models.Usermatrix{}, errs
	}
	fmt.Println("print results")
	fmt.Println(matrixlist)
	return &matrixlist, errs
}
