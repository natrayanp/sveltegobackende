package commonfuncs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mitchellh/mapstructure"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/vgarvardt/gue/v2"
)

// PackageFetch returns menus and access rights for it for a user.
// Parameters:
// packfuncid --> If you want only the packsfuncs sent. Send the PACKID form AC.PACKS as array.
//					This forcefully sent only those packs if it exists at company and user level(refer query)
// companyid  --> Always send company id here
//					Front don't send company id send userinfo.Companyid from calling side
//					else send whatever company id received from front end
// It returns PacksResp struct which is self explanatory and error.
func RoleFetch(app *application.Application, w http.ResponseWriter, r *http.Request, rolereq *models.RoleReq) (*models.RoleResp, error) {
	fmt.Println("----------------- ROLE Fetch START -------------------")

	var data string
	var qry string
	var datosend models.RoleResp
	var err error

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "ROLE-FETCH")
	if errs != nil {
		return &models.RoleResp{}, errs
	}

	//This will give all the modules available for role creation for the company -- START
	//This is always has same value irrespective of company has role or not.

	var availmod []models.TtblMytree
	/*
		qry = `WITH MYAA AS(
				SELECT c.COMPANYID,'PUBLIC' AS BRANCHID,'' AS Roledetailid,'' AS ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,
				c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,
				array_fill(FALSE, ARRAY[array_length(c.allowedops,1)])  AS allowedopsval,$2 as userid,
				--CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
				'Availablemodules' AS basketname, false as open
				FROM ac.COMPANYPACKS_PACKS_VIEW C
				LEFT JOIN ac.ROLE_USER_VIEW B ON C.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = $2
				WHERE C.COMPANYID = $1
				) SELECT Y.*,
					CASE WHEN (NULLIF(Y.allowedopsval, '{NULL}')) IS NULL AND (Y.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc
					FROM MYAA as Y;`
	*/

	qry = `WITH MY AS (
				SELECT DISTINCT RDPACKFUNCID FROM AC.ROLEDETAILS WHERE RDROLEMASTERID IN (SELECT ROLEMASTERID FROM AC.userrole WHERE userid = $1 AND COMPANYID = $2)
			), MYPF AS (
		   		SELECT z.*,'PUBLIC' AS BRANCHID,'' AS Roledetailid,'' AS ROLEMASTERID,
				CASE WHEN d.RDPACKFUNCID IS NULL THEN TRUE ELSE FALSE END as disablefunc,
				array_fill(FALSE, ARRAY[array_length(Z.allowedops,1)])  AS allowedopsval,
				--CASE WHEN (NULLIF(b.rdallowedopsval, '{NULL}')) IS NULL AND (z.TYPE IN ('function','module')) THEN TRUE ELSE FALSE END AS disablefunc,
				'Availablemodules' AS basketname, false as open
		   		FROM AC.COMPANYPACKS_PACKS_VIEW z
		   		--CROSS JOIN ac.rolemaster a
		   		--LEFT JOIN ac.roledetails b ON a.rolemasterid = b.rdrolemasterid AND z.packid = b.rdpackfuncid AND B.COMPANYID = Z.COMPANYID
		   		LEFT JOIN MY d ON z.PACKFUNCID = d.RDPACKFUNCID
				--JOIN MYP c on a.rolemasterid = c.rdrolemasterid AND z.packgroupid = c.gid
		   		WHERE  z.COMPANYID =  $2
	  		) SELECT * FROM MYPF;`

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &availmod, userinfo.UUID, rolereq.Companyid),
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
			SlugCode:   "ROLE-FETCH",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &models.RoleResp{}, errs
	}
	fmt.Println(availmod)
	createDataTree(&availmod)
	datosend.Availablemodules = availmod
	fmt.Println(availmod)
	fmt.Println(datosend)
	fmt.Println("-====++++datosend-====++++")

	//This will give all the rolewise details for the company if Roles are already created -- END

	//Fetch all Roles available for the company -- START

	//var selmod []models.RoleSelectModu
	//var selmod []models.TtblMytree
	//var selmod []string
	var selmod []models.TmpRoleSelectModu
	var fselmod []models.RoleSelectModu

	/* @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@*/
	qry = ` WITH MY AS (
				SELECT DISTINCT RDPACKFUNCID FROM AC.ROLEDETAILS WHERE RDROLEMASTERID IN (SELECT ROLEMASTERID FROM AC.userrole WHERE userid = $1 AND COMPANYID = $2)
		), MYP AS(
	   			
				   SELECT distinct rdrolemasterid,packgroupid gid FROM AC.ROLEDETAILS A
				   LEFT JOIN AC.PACKS B ON A.rdpackfuncid = B.PACKID			   
	   	), MYPF AS (
	   			SELECT z.*,a.*,b.roledetailid,b.rdallowedopsval as allowedopsval,
	   			CASE WHEN d.RDPACKFUNCID IS NULL THEN TRUE ELSE FALSE END as disablefunc,
				--CASE WHEN (NULLIF(b.rdallowedopsval, '{NULL}')) IS NULL AND (z.TYPE IN ('function','module')) THEN TRUE ELSE FALSE END AS disablefunc,
				'selectedmodules' AS basketname, false as open
	   			FROM AC.COMPANYPACKS_PACKS_VIEW z
	   			CROSS JOIN ac.rolemaster a
	   			LEFT JOIN ac.roledetails b ON a.rolemasterid = b.rdrolemasterid AND z.packid = b.rdpackfuncid AND B.COMPANYID = Z.COMPANYID
	   			LEFT JOIN MY d ON z.PACKFUNCID = d.RDPACKFUNCID
				JOIN MYP c on a.rolemasterid = c.rdrolemasterid AND z.packgroupid = c.gid 
	   			where  a.COMPANYID != 'PUBLIC' AND z.COMPANYID =  $2 AND a.rmstatus NOT IN ('D')
	   		) select SD.ROLEMASTERID,SD.RMNAME,SD.RMDISPLAYNAME,SD.RMDESCRIPTION,json_agg(SD) AS modules FROM MYPF sd GROUP BY  sd.rolemasterid,sd.rmname,sd.rmdisplayname ,sd.rmdescription;`
	/*@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@*/

	stmts = []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &selmod, userinfo.UUID, rolereq.Companyid),
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
			SlugCode:   "DOMAINREG-UPDATE",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return &models.RoleResp{}, errs
	}
	fmt.Println("before-----")
	fselmod = make([]models.RoleSelectModu, len(selmod))
	for i, s := range selmod {
		fselmod[i].Rolemasterid = s.Rolemasterid
		fselmod[i].Rolename = s.Rmname
		fselmod[i].Roledisplayname = s.Rmdisplayname
		fselmod[i].Roledescription = s.Rmdescription
		mapstructure.Decode(s.Modules, &fselmod[i].Modules)
		createDataTree(&fselmod[i].Modules)
	}

	datosend.Selectedmodules = fselmod
	return &datosend, err

}

// PackageFetch returns menus and access rights for it for a user.
// Parameters:
// packfuncid --> If you want only the packsfuncs sent. Send the PACKID form AC.PACKS as array.
//					This forcefully sent only those packs if it exists at company and user level(refer query)
// companyid  --> Always send company id here
//					Front don't send company id send userinfo.Companyid from calling side
//					else send whatever company id received from front end
// It returns PacksResp struct which is self explanatory and error.
func Roleupdate(app *application.Application, w http.ResponseWriter, r *http.Request, rolereq *models.RolesaveReq) error {
	fmt.Println("----------------- ROLE Update START -------------------")

	var data string
	var qry string
	var err error
	var rlmasid string
	var mtx pgx.Tx
	var dbtt dbtran.TranType

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "ROLE-MASTER-UPDATE")
	if errs != nil {
		return errs
	}

	fmt.Println(userinfo)
	//This will give all the modules available for role creation for the company -- START
	//This is always has same value irrespective of company has role or not.

	var rolemastresp []models.TblRolemaster
	var stmts []*dbtran.PipelineStmt
	var stmtcp *dbtran.PipelineStmt
	var myc dbtran.Resultset
	rlmas := rolereq.Rolemaster
	rldet := rolereq.Roledetails
	rladt := rolereq.Audit
	dbtt = dbtran.TranTypeFullSet
	mtx = nil
	valid := map[string]bool{"I": true, "U": true}
	validadt := map[string]bool{"I": true, "U": true, "M": true}
	fmt.Println(rlmas)
	if rlmas.Action != "D" {
		if valid[rlmas.Action] {
			switch rlmas.Action {
			case "I":
				fmt.Println("Inside delete")
				qry = `INSERT INTO ac.ROLEMASTER VALUES (DEFAULT,$1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) RETURNING *;`

				stmts = []*dbtran.PipelineStmt{
					dbtran.NewPipelineStmt("select", qry, &rolemastresp, rlmas.Rolename, rlmas.Roledisplayname, rlmas.Roledescription, rolereq.Companyid, rolereq.Branchid, "A"),
				}

			case "U":
				fmt.Println("Inside update")
				qry = `UPDATE ac.ROLEMASTER 
				SET rmname = $1, rmdisplayname = $2, rmdescription = $3, lmtime = CURRENT_TIMESTAMP
				WHERE rolemasterid = $4 AND companyid = $5 
		    		RETURNING *;`

				stmts = []*dbtran.PipelineStmt{
					dbtran.NewPipelineStmt("select", qry, &rolemastresp, rlmas.Rolename, rlmas.Roledisplayname, rlmas.Roledescription, rlmas.Rolemasterid, rolereq.Companyid),
				}
			}

			mtx, err = dbtran.WithTransaction(ctx, dbtt, app.DB.Client, mtx, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
				err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
				return err
			})

			fmt.Println("Just outside ")
			fmt.Println(rolemastresp)
			fmt.Println(err)

			if err != nil {
				goto mydberr
			}
			dbtt = dbtran.TranTypeLastSet
		}

		fmt.Println("----------------- ROLE master Update done -------------------")
		fmt.Println(rolemastresp)
		if rlmas.Rolemasterid == "NEW" {
			rlmasid = rolemastresp[0].Rolemasterid
		} else {
			rlmasid = rlmas.Rolemasterid
		}

		fmt.Println(len(rldet))
		stmts = []*dbtran.PipelineStmt{}

		if len(rldet) > 0 {
			fmt.Println("rldet count")
			for _, s := range rldet {
				s.Rolemasterid = rlmasid
				stmtcp = nil
				fmt.Println("record val", s)
				switch s.Action {
				case "I":
					qry = `INSERT INTO ac.ROLEDETAILS (rdrolemasterid, rdpackfuncid,  companyid, branchid,rdallowedopsval,octime,lmtime) 
			VALUES ($1,$2,$3,$4,$5,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) ;`
					stmts = append(stmts, dbtran.NewPipelineStmt("insert", qry, &myc, s.Rolemasterid, s.Packid, rolereq.Companyid, rolereq.Branchid, s.Allowedopsval))

				case "U":
					qry = `UPDATE ac.ROLEDETAILS 
					SET rdallowedopsval = $1, lmtime = CURRENT_TIMESTAMP
					WHERE rdrolemasterid = $2 AND roledetailid = $3 AND  companyid = $4 ;`

					stmts = append(stmts, dbtran.NewPipelineStmt("update", qry, &myc, s.Allowedopsval, s.Rolemasterid, s.Roledetailid, rolereq.Companyid))

				case "D":
					qry = `DELETE FROM ac.ROLEDETAILS 				
						WHERE rolemasterid = $1 AND rdroledetailid = $2 AND rdpackfuncid = $3 AND companyid = $4;`

					stmts = append(stmts, dbtran.NewPipelineStmt("delete", qry, &myc, s.Rolemasterid, s.Roledetailid, s.Packid, rolereq.Companyid))

				}
				fmt.Println("-------------------stmts both--------------")
				fmt.Println(stmtcp)
				fmt.Println(stmts)

			}

			fmt.Println("-------------------stmts--------------")
			fmt.Println(stmts)

			_, err = dbtran.WithTransaction(ctx, dbtt, app.DB.Client, mtx, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
				err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
				return err
			})

			if err != nil {
				goto mydberr
			}

		}
	} else {
		//Delete Rolemaster and its Roledetails
		fmt.Println("Delete the rolemaster")
		qry = `UPDATE ac.ROLEMASTER 
	SET rmstatus = 'D' , lmtime = CURRENT_TIMESTAMP
	WHERE rolemasterid = $1 AND companyid = $2 ;`
		fmt.Println(rlmas)
		fmt.Println(rlmas.Rolemasterid)
		stmts = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("update", qry, &myc, rlmas.Rolemasterid, rolereq.Companyid),
		}

		_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
			err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
			return err
		})

		if err != nil {
			goto mydberr
		}
	}

	type auditentryargs struct {
		Itemid string
		Action string
		Oldval interface{}
		Newval interface{}
		User   string
		Time   time.Time
	}

	if validadt[rladt.Action] {

		rladt.Itemkeys.Rolemasterid = rlmasid

		args, err1 := json.Marshal(auditentryargs{Itemid: rladt.Itemid, Action: rladt.Action,
			Oldval: rladt.Oldvalue, Newval: rladt.Newvalue,
			User: userinfo.UUID, Time: time.Now()})

		fmt.Println(err1)

		if err := app.Que.Enquejob(&gue.Job{Type: "Auditentry", Args: args}); err != nil {

			dd := httpresponse.SlugResponse{
				Err:        err,
				ErrType:    httpresponse.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Data:       map[string]interface{}{"message": data},
				SlugCode:   "ROLE-MASTER-UPDATE-ENQUE",
				LogMsg:     "update audit table -> Role audit update failed~rolemasterid: " + rlmasid,
			}
			//dd.HttpRespondWithError()
			dd.HttpRespond()
			return err
			//return err
		}
	}

mydberr:
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
			SlugCode:   "ROLE-MASTER-UPDATE",
			LogMsg:     pgErr.Error(),
		}
		dd.HttpRespond()
		return err
	}

	//fmt.Println(availmod)
	return err
}

/*

fmt.Println("**********************&&&&&&&&&&&&&&&&&&&&&&&&&&&&^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^%%%%%%%%%%%%%%%%%%%%%%%%%%%")
fmt.Println(rolereq)
fmt.Println("**********************&&&&&&&&&&&&&&&&&&&&&&&&&&&&^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^%%%%%%%%%%%%%%%%%%%%%%%%%%%")

					qry = `WITH MYVA AS (
								SELECT c.COMPANYID,A.BRANCHID,B.Roledetailid,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,
								c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,B.USERID,A.displayname as roledisplay,
								--CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
								'selectedmodules' AS basketname, false as open
								from ac.COMPANYPACKS_PACKS_VIEW C
								LEFT JOIN ac.rolemaster A ON A.COMPANYID = $1
								LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = $2 AND B.ROLEMASTERID = A.ROLEMASTERID
								WHERE C.COMPANYID = $1 AND A.ROLEMASTERID is NOT NULL
								ORDER BY A.ROLEMASTERID
							) , MYVAOS AS  (
								SELECT  COMPANYID,branchid,rolemasterid,packfuncid,ALLOWEDOPSVAL FROM  ac.ROLE_USER_VIEW WHERE COMPANYID = $1 AND rolemasterid in (select distinct rolemasterid from myva) GROUP BY companyid,branchid,rolemasterid,packfuncid,ALLOWEDOPSVAL
							), MYLAST AS (SELECT Y.*,PO.allowedopsval,CASE WHEN (NULLIF(PO.allowedopsval, '{NULL}')) IS NULL AND (Y.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc FROM MYVA y
								  JOIN MYVAOS as PO ON PO.COMPANYID = Y.COMPANYID AND PO.rolemasterid=Y.rolemasterid AND PO.Branchid = y.branchid AND PO.packfuncid = y.packid
							) SELECT sd.rolemasterid,sd.name,sd.roledisplay,sd.description,json_agg(SD) AS modules FROM mylast sd GROUP BY  sd.rolemasterid,sd.name,sd.roledisplay,sd.description;`


				qry = `WITH RECURSIVE MDATAR AS
				(
					SELECT * from ac.ROLE_USER_VIEW where companyid = $1
				),
				MyTree AS
				(
					SELECT C.COMPANYID,$2 As branchid,A.ROLEMASTERID,C.packid,c.name ,c.displayname ,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
					A.rmDISPLAYNAME,A.rmname,A.rmdescription,A.rmstatus,
					--CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
					CASE WHEN (NULLIF(A.allowedopsval, '{NULL}')) IS NULL AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
					'SELECTedmodules' AS basketname, false as open
					from ac.COMPANYPACKS_PACKS_VIEW C
					LEFT JOIN MDATAR A ON A.packfuncid = C.PACKID
					WHERE PACKGROUPID = (SELECT DISTINCT unnest(PACKGROUPID) FROM AC.PACKS where PACKID IN (SELECT DISTINCT PACKID FROM MDATAR))
						UNION
					SELECT C.COMPANYID,$2 As branchid,T.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
					t.rmDISPLAYNAME,t.rmname,t.rmdescription,t.rmstatus,
					--CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
					CASE WHEN (NULLIF(A.allowedopsval, '{NULL}')) IS NULL AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
					'SELECTedmodules' AS basketname, false as open
					from ac.COMPANYPACKS_PACKS_VIEW C
					LEFT JOIN MDATAR A ON  A.packfuncid = C.PACKID
					JOIN MyTree AS t ON C.packid = ANY(t.parent)

				),MYLAST AS (
				SELECT * FROM MyTree
				WHERE COMPANYID = $1
			)SELECT sd.rolemasterid,sd.rmname,sd.rmdisplayname ,sd.rmdescription,json_agg(SD) AS modules FROM mylast sd GROUP BY  sd.rolemasterid,sd.rmname,sd.rmdisplayname ,sd.rmdescription;`


		qry = `WITH RECURSIVE MDATAR AS
		(
			SELECT * from ac.ROLE_USER_VIEW where companyid = $1
		),
		MyTree AS
		(
			SELECT C.COMPANYID,$2 As branchid,A.ROLEMASTERID,C.packid,c.name ,c.displayname ,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
			A.rmDISPLAYNAME,A.rmname,A.rmdescription,A.rmstatus,
			--CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
			CASE WHEN (NULLIF(A.allowedopsval, '{NULL}')) IS NULL AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
			'SELECTedmodules' AS basketname, false as open
			from ac.COMPANYPACKS_PACKS_VIEW C
			LEFT JOIN MDATAR A ON A.packfuncid = C.PACKID
			WHERE PACKGROUPID && (SELECT DISTINCT unnest(PACKGROUPID) FROM AC.PACKS where PACKID IN (SELECT DISTINCT PACKID FROM MDATAR))
		),MYLAST AS (
		SELECT * FROM MyTree
		WHERE COMPANYID = $1
	)SELECT sd.rolemasterid,sd.rmname,sd.rmdisplayname ,sd.rmdescription,json_agg(SD) AS modules FROM mylast sd GROUP BY  sd.rolemasterid,sd.rmname,sd.rmdisplayname ,sd.rmdescription;`
*/
/*
func remove(slice []models.TtblMytree, s int) []models.TtblMytree {
	return append(slice[:s], slice[s+1:]...)
}


	if cmpy, errs = CompanyCheck(app, w, r, companyid); errs != nil {
		fmt.Println("TODO: Error handling")
		return &models.PacksResp{}, errs
	}



	qry =

		fmt.Println("---------------$$$end5")
		//fmt.Println(err.Error())
		if err != nil {
			//https://github.com/jackc/pgx/issues/474
			var pgErr *pgconn.PgError
			fmt.Println("---------------$$$end5a")
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
				SlugCode:   "DOMAINREG-UPDATE",
				LogMsg:     pgErr.Error(),
			}
			dd.HttpRespond()
			return &models.PacksResp{}, errs
		}

		fmt.Println("---------------$$$end6")
		dd, _ := json.Marshal(myc)
		fmt.Println(string(dd))
		fmt.Println("---------------$$$end6a")
		fmt.Printf("&myc is: %p\n", &myc)

		fmt.Println("---------------$$$end6b")
		dd1, _ := json.Marshal(myc)
		fmt.Printf("&myc is: %p\n", &myc)
		fmt.Println(string(dd1))
		fmt.Println("---------------$$$end7")
		fmt.Println("----------------- PACKAGE FETCH END -------------------")

		datosend.CpyLvlTreeforCpy = []models.ActiveEntityTree{{EntityType: "company",
			Entityid:   companyid,
			EntityTree: myc}}
		datosend.BrnLvlTreeforCpy = []models.ActiveEntityTree{}

		//return &datosend, nil

	} else {

		myca = make([]models.ActiveEntityTree, len(datosend.BranchLst))

		//https: //www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html
		//https://stackoverflow.com/questions/18805416/waiting-on-an-indeterminate-number-of-goroutines

		if havbrndetail {
			fmt.Println("--- start for loop")
			for i, s := range datosend.BranchLst {
				//var mycpp []models.TblMytree
				lo := i
				myca[lo].EntityType = "branch"
				myca[lo].Entityid = s.Branchid
				fmt.Println(lo, s.Branchid)
				wgbr.Add(1)
				go func() {
					defer wgbr.Done()
					//var mycpp []models.TblMytree
					var mycpp []models.TtblMytree


					stmts = []*dbtran.PipelineStmt{
						dbtran.NewPipelineStmt("select", qry, &mycpp, userinfo.UUID, companyid, s.Branchid),
					}

					_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
						err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
						return err
					})

					if err != nil {
						fmt.Println("TODO error handlin in go routine")
					}

					fmt.Println("---------------$$$end6aw")
					fmt.Println(mycpp)
					fmt.Printf("&myc is: %p\n", &mycpp)
					createDataTree(&mycpp)

					myca[lo].EntityTree = mycpp
					fmt.Println("---------------$$$end6bw")
					dd1, _ := json.Marshal(myca[lo])
					fmt.Printf("&myc is: %p\n", &myca[lo])
					fmt.Println(string(dd1))
					fmt.Println("---------------$$$end7w")
				}()
			}
			wgbr.Wait()
			fmt.Println("All routines completed")
			fmt.Printf("&myc is: %p\n", &myca)
			dd1, _ := json.Marshal(myca)
			fmt.Println(string(dd1))
			datosend.BrnLvlTreeforCpy = myca

		} else {
			/*
				TODO: navigate to ADDBRANCH
*/
/*
		datosend.BrnLvlTreeforCpy = []models.ActiveEntityTree{}
	}

	mycacp = make([]models.ActiveEntityTree, len(datosend.CompanyLst))

	for i, s := range datosend.CompanyLst {
		//var mycppp []models.TblMytree
		mycacp[i].EntityType = "company"
		mycacp[i].Entityid = s.Companyid
		fmt.Println(i, s.Companyid)
		wgcp.Add(1)
		go func() {
			defer wgcp.Done()
			//var mycppp []models.TblMytree
			var mycppp []models.TtblMytree
			qry = `WITH RECURSIVE MyTree AS
			(
				SELECT A.*,false as open,B.roledetailid,B.rolemasterid,B.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs A
				LEFT JOIN ac.roledetails B ON A.packid = B.PACKFUNCID
				WHERE A.packid IN
				(
					(	SELECT PACKFUNCID FROM ac.roledetails
						WHERE rolemasterid IN
							(SELECT DISTINCT rolemasterid FROM ac.userrole
								WHERE userid = $1
								AND status NOT IN ('D','I')
								AND companyid = $2
							)
						INTERSECT
						SELECT PACKFUNCID from ac.companypacks
							WHERE companyid = $2
							AND status NOT IN ('D','I')
							AND startdate <=  CURRENT_DATE
							AND expirydate >= CURRENT_DATE
					)
				)
				AND A.menulevel IN ('COMPANY')
				UNION
				SELECT M.*,false as open,N.roledetailid,N.rolemasterid,N.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs M
				LEFT JOIN ac.roledetails N ON M.packid = N.PACKFUNCID
				JOIN MyTree AS t ON M.packid = ANY(t.parent)
					/*SELECT m.*,false as open FROM ac.packs AS m JOIN MyTree AS t ON m.packid = ANY(t.parent)*/

/*
				)
				SELECT * FROM MyTree ORDER BY SORTORDER,TYPE,NAME;`

				stmts = []*dbtran.PipelineStmt{
					dbtran.NewPipelineStmt("select", qry, &mycppp, userinfo.UUID, companyid),
				}

				_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
					err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
					return err
				})

				if err != nil {
					fmt.Println("TODO error handlin in go routine")
				}

				fmt.Println("---------------$$$end6aw")
				fmt.Println(mycppp)
				fmt.Printf("&myc is: %p\n", &mycppp)
				createDataTree(&mycppp)
				mycacp[i].EntityTree = mycppp
				fmt.Println("---------------$$$end6bw")
				dd1, _ := json.Marshal(mycacp[i])
				fmt.Printf("&myc is: %p\n", &mycacp[i])
				fmt.Println(string(dd1))
				fmt.Println("---------------$$$end7w")
			}()
		}
		wgcp.Wait()

		fmt.Println("All routines completed")
		fmt.Printf("&myc is: %p\n", &mycacp)
		dd1, _ := json.Marshal(mycacp)
		fmt.Println(string(dd1))
		datosend.CpyLvlTreeforCpy = mycacp
	}

	return &datosend, err
}
*/

/*
	WITH MYVA AS (
		select A.COMPANYID,A.BRANCHID,B.Roledetailid,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,B.allowedopsval,B.USERID,
		CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
		'availablemodules' AS basketname, false as open
		from ac.COMPANYPACKS_PACKS_VIEW C
		LEFT JOIN ac.rolemaster A ON A.COMPANYID = 'CPYID132'
		LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = 'USER5555' AND B.ROLEMASTERID = A.ROLEMASTERID
		WHERE C.COMPANYID = 'CPYID132'
		ORDER BY A.ROLEMASTERID
		) SELECT * from myva



						qry = `WITH MyTree1 AS
						(
						WITH RECURSIVE MyTree AS
									(
										SELECT A.*,false as open,B.roledetailid,B.rolemasterid,B.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs A
										LEFT JOIN ac.roledetails B ON A.packid = B.PACKFUNCID
										WHERE A.packid IN
										(
											(	SELECT PACKFUNCID FROM ac.roledetails
												WHERE rolemasterid IN
													(SELECT DISTINCT rolemasterid FROM ac.userrole
														WHERE status NOT IN ('D','I')
														AND companyid = 'CPYID132'
													 	AND branchid && ('JDJDJD'::VARCHAR||ARRAY['ALL']::VARCHAR[])
													 	AND rolemasterid = 'ROLMA1'
													)
												INTERSECT
												SELECT PACKFUNCID from ac.companypacks
													WHERE companyid = 'CPYID132'
													AND status NOT IN ('D','I')
													AND startdate <=  CURRENT_DATE
													AND expirydate >= CURRENT_DATE
											)
										)
										AND A.menulevel = ANY( ARRAY['COMPANY'])
										UNION
										SELECT M.*,false as open,N.roledetailid,N.rolemasterid,N.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs M
										LEFT JOIN ac.roledetails N ON M.packid = N.PACKFUNCID
										JOIN MyTree AS t ON M.packid = ANY(t.parent)

									)
								SELECT * FROM MyTree AS MT WHERE MT.PARENT = ARRAY[NULL]  ORDER BY SORTORDER,TYPE,NAME

							), MyTree2 AS
							(
							WITH RECURSIVE MyTreeN AS
									(
										(
										SELECT A1.*,false as open,B1.roledetailid,B1.rolemasterid,B1.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs A1
											LEFT JOIN ac.roledetails B1 ON A1.packid = B1.PACKFUNCID AND rolemasterid = 'ROLMA1'
											WHERE A1.packid IN (SELECT packid FROM MYTREE1)
										)
											UNION
											SELECT M1.*,false as open,N1.roledetailid,N1.rolemasterid,N1.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs M1
											LEFT JOIN ac.roledetails N1 ON M1.packid = N1.PACKFUNCID AND rolemasterid = 'ROLMA1'

											JOIN MyTreeN AS t1 ON t1.packid = ANY(m1.parent)
										)
										(SELECT * FROM MyTreeN WHERE TYPE IN ('module','pack')
										UNION
										(SELECT M2.* FROM MyTreeN M2, ac.companypacks M3, ac.roledetails M4
											WHERE M2.packid = M3.PACKFUNCID AND M2.TYPE = 'function'
											)
										 )		ORDER BY SORTORDER,TYPE,NAME
							), MyTree3 AS
							(

										(
											SELECT *,true as witho FROM MyTree2 WHERE TYPE IN ('module','pack')
											UNION
											(SELECT *,false as witho  FROM MyTree2
												WHERE TYPE = 'function'
										 		AND packid not in(
										 		(SELECT PACKFUNCID FROM ac.roledetails
													WHERE rolemasterid IN
													(SELECT DISTINCT rolemasterid FROM ac.userrole
														WHERE status NOT IN ('D','I')
													 	AND userid = 'bQPmqQPcVVWS6paQYQ4eN6eReI83'
														AND companyid = 'CPYID132'
													 	AND branchid && ('JDJDJD'::VARCHAR||ARRAY['ALL']::VARCHAR[])
													)
										 		)
											)ORDER BY SORTORDER,TYPE,NAME
											 )
										 union
										 (SELECT *,true as witho  FROM MyTree2
											WHERE TYPE = 'function'
										 		AND packid in(
										 		(SELECT PACKFUNCID FROM ac.roledetails
													WHERE rolemasterid IN
													(SELECT DISTINCT rolemasterid FROM ac.userrole
														WHERE status NOT IN ('D','I')
													 	AND userid = 'bQPmqQPcVVWS6paQYQ4eN6eReI83'
														AND companyid = 'CPYID132'
													 	AND branchid && ('JDJDJD'::VARCHAR||ARRAY['ALL']::VARCHAR[])
													)
										 		)
											)

										 )		ORDER BY SORTORDER,TYPE,NAME
								)
							)
							SELECT * FROM MyTree3 ORDER BY SORTORDER,TYPE,NAME;
							++++++++++++++++++++
							company view:
								select * from ac.companypacks
								recursive from packs


							Role view:
								select companyid, roleid, userid, funcid,xxxx from ac.rolemaster
								left join ac.roldetails on <roleid>
								left join ac.userrole on <roledi>

								select A.companyid, A.branchid, A.rolemasterid, B.packfuncid,b.allowedopsval,c.userid from ac.rolemaster A
				left join ac.roledetails B on A.rolemasterid = B.rolemasterid
				left join ac.userrole C on A.rolemasterid = C.rolemasterid






								select c.* from companyview c
								Left join roleview r1 ON decode(r1.userid,null,'N',<LOGGEDINUSERID>,'Y','N') = 'Y' AND c.companyid = r1.companyid AND c.funcid = r1.funcid
								Left join roleview r2 ON decode(r2.roleid,null,'Y',<ROLEID>,'Y','N') = 'Y' AND c.companyid = r2.companyid AND c.funcid = r2.funcid
								where c.companyid = <companyid>

		select A.COMPANYID,A.BRANCHID,A.ROLEMASTERID,B.packfuncid,B.allowedopsval,B.USERID,C.* from ac.COMPANYPACKS_PACKS_VIEW C
		LEFT JOIN ac.rolemaster A ON A.COMPANYID = 'CPYID132' AND A.ROLEMASTERID = 'ROLMA3'
		LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND A.ROLEMASTERID = B.ROLEMASTERID AND B.packfuncid = C.PACKID AND B.USERID = 'USER5555' AND B.BRANCHID = 'PUBLIC' AND A.ROLEMASTERID = 'ROLMA3'
		WHERE C.COMPANYID = 'CPYID132'
		ORDER BY A.ROLEMASTERID

		--- Recent update on 13-Sep
		select C.*,A.COMPANYID,A.BRANCHID,A.ROLEMASTERID,B.packfuncid,B.allowedopsval,B.USERID from ac.COMPANYPACKS_PACKS_VIEW C
		LEFT JOIN ac.rolemaster A ON A.COMPANYID = 'CPYID132'
		LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND A.ROLEMASTERID = B.ROLEMASTERID AND B.packfuncid = C.PACKID AND B.USERID = 'USER5555' AND B.userbranchid && ARRAY['ALL'::VARCHAR,'PUBLIC'::VARCHAR]
		WHERE C.COMPANYID = 'CPYID132'
		ORDER BY A.ROLEMASTERID

		--- Recent update on 14-Sep
		select C.*,A.COMPANYID AS CPID,A.BRANCHID AS BRID,A.ROLEMASTERID AS RLMID,B.packfuncid,B.allowedopsval,B.USERID from ac.COMPANYPACKS_PACKS_VIEW C
		LEFT JOIN ac.rolemaster A ON A.COMPANYID = 'CPYID132'
		LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = 'USER5555'
		WHERE C.COMPANYID = 'CPYID132'
		ORDER BY A.ROLEMASTERID

		WITH MYVA AS (
		select A.COMPANYID,A.BRANCHID,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,B.allowedopsval,B.USERID,
		CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc
		from ac.COMPANYPACKS_PACKS_VIEW C
		LEFT JOIN ac.rolemaster A ON A.COMPANYID = 'CPYID132'
		LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = 'USER5555'
		WHERE C.COMPANYID = 'CPYID132'
		ORDER BY A.ROLEMASTERID
		) SELECT ROLEMASTERID,array_agg(myva) as ad from myva group by myva.rolemasterid


				WITH MYVA AS (
		select A.COMPANYID,A.BRANCHID,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,B.allowedopsval,B.USERID,
		CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
		'availablemodules' AS basketname, false as open
		from ac.COMPANYPACKS_PACKS_VIEW C
		LEFT JOIN ac.rolemaster A ON A.COMPANYID = 'CPYID132'
		LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = 'USER5555'
		WHERE C.COMPANYID = 'CPYID132'
		ORDER BY A.ROLEMASTERID
		) SELECT * from myva

							SELECT * FROM COMPANYPACKS_PACKS_VIEW;
							SELECT * FROM ROLEMASTER_PACKS;


						stmts = []*dbtran.PipelineStmt{
							dbtran.NewPipelineStmt("select", qry, &mycusr, userinfo.UUID, rolereq.Companyid, rolereq.Branchid),
						}

						_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
							err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
							return err
						})

						fmt.Println("---------------$$$end5")
						//fmt.Println(err.Error())
						if err != nil {
							//https://github.com/jackc/pgx/issues/474
							var pgErr *pgconn.PgError
							fmt.Println("---------------$$$end5a")
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
								SlugCode:   "DOMAINREG-UPDATE",
								LogMsg:     pgErr.Error(),
							}
							dd.HttpRespond()
							return &models.PacksResp{}, errs
						}

						createDataTree(&mycusr)

						qry1 = `WITH RECURSIVE MyTree AS
								(
									SELECT A.*,false as open,B.roledetailid,B.rolemasterid,B.allowedopsval,'Availablemodules' AS basketname  FROM ac.packs A
									LEFT JOIN ac.roledetails B ON A.packid = B.PACKFUNCID
									WHERE A.packid IN
									(
										(
											SELECT PACKFUNCID from ac.companypacks
												WHERE companyid = $1
												AND branchid ANY ARRAY['ALL'::VARCHAR, $2::VARCHAR]
												AND status NOT IN ('D','I')
												AND startdate <=  CURRENT_DATE
												AND expirydate >= CURRENT_DATE
										)
									)
									UNION
									SELECT M.*,false as open,N.roledetailid,N.rolemasterid,N.allowedopsval,'Availablemodules' AS basketname FROM ac.packs M
									LEFT JOIN ac.roledetails N ON M.packid = N.PACKFUNCID
									JOIN MyTree AS t ON M.packid = ANY(t.parent)
								)
								SELECT * FROM MyTree ORDER BY SORTORDER,TYPE,NAME;`

						stmts1 := []*dbtran.PipelineStmt{
							dbtran.NewPipelineStmt("select", qry1, &myc, userinfo.UUID, rolereq.Companyid, rolereq.Branchid),
						}

						_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
							err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
							return err
						})

						fmt.Println("---------------$$$end5")
						//fmt.Println(err.Error())
						if err != nil {
							//https://github.com/jackc/pgx/issues/474
							var pgErr *pgconn.PgError
							fmt.Println("---------------$$$end5a")
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
								SlugCode:   "DOMAINREG-UPDATE",
								LogMsg:     pgErr.Error(),
							}
							dd.HttpRespond()
							return &models.PacksResp{}, errs
						}

						createDataTree(&myc)

						rlmd := []string{}

						qry2 = `SELECT DISTINCT rolemasterid FROM AC.ROLEMASTER
														WHERE status NOT IN ('D','I')
														AND companyid = $1
														AND branchid ANY ARRAY['ALL'::VARCHAR, $2::VARCHAR]`

						stmts2 := []*dbtran.PipelineStmt{
							dbtran.NewPipelineStmt("select", qry2, &rlmd, rolereq.Companyid, rolereq.Branchid),
						}

						_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
							err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts2...)
							return err
						})

						fmt.Println("---------------$$$end5")
						//fmt.Println(err.Error())
						if err != nil {
							//https://github.com/jackc/pgx/issues/474
							var pgErr *pgconn.PgError
							fmt.Println("---------------$$$end5a")
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
								SlugCode:   "DOMAINREG-UPDATE",
								LogMsg:     pgErr.Error(),
							}
							dd.HttpRespond()
							return &models.PacksResp{}, errs
						}
						/*

							myca = make([]models.RoleSelectModu, len(rlmd))

							qry3 = `WITH RECURSIVE MyTree AS
															(
																SELECT A.*,false as open,B.roledetailid,B.rolemasterid,B.allowedopsval,'Selectedmodules' AS basketname  FROM ac.packs A
																LEFT JOIN ac.roledetails B ON A.packid = B.PACKFUNCID
																WHERE A.packid IN
																(
																	(	SELECT PACKFUNCID FROM ac.roledetails
																		WHERE status NOT IN ('D','I')
																				AND companyid = $2
																				AND branchid ANY ARRAY['ALL'::VARCHAR, $3::VARCHAR]
																				 AND rolemasterid = $4
																		INTERSECT
																		SELECT PACKFUNCID from ac.companypacks
																			WHERE companyid = $2
																			AND status NOT IN ('D','I')
																			AND startdate <=  CURRENT_DATE
																			AND expirydate >= CURRENT_DATE
																	)
																)WITH MYP AS(
SELECT DISTINCT RDPACKFUNCID FROM AC.ROLEDETAILS WHERE RDROLEMASTERID IN (SELECT ROLEMASTERID FROM AC.userrole WHERE userid = 'aeMhBaHZB0ShHXf8QmhLSZ4Ap9m2' AND COMPANYID = 'CPYID21')
	), MYPF AS (
			SELECT z.*,a.*,b.roledetailid,b.rdallowedopsval,
			CASE WHEN c.RDPACKFUNCID IS NULL THEN FALSE ELSE TRUE END as USERPACKSIDACCESS
			FROM AC.COMPANYPACKS_PACKS_VIEW z
			CROSS JOIN ac.rolemaster a
			LEFT JOIN ac.roledetails b ON a.rolemasterid = b.rdrolemasterid AND z.packid = b.rdpackfuncid AND B.COMPANYID = Z.COMPANYID
			LEFT JOIN MYP c ON z.PACKFUNCID = c.RDPACKFUNCID
		) select *
		from MYPF
																UNION
																SELECT M.*,false as open,N.roledetailid,N.rolemasterid,N.allowedopsval,'Selectedmodules' AS basketname FROM ac.packs M
																LEFT JOIN ac.roledetails N ON M.packid = N.PACKFUNCID
																JOIN MyTree AS t ON M.packid = ANY(t.parent)
														)
															SELECT * FROM MyTree ORDER BY SORTORDER,TYPE,NAME;`
								for i, s := range rlmd {
									//var mycpp []models.TblMytree
									lo := i
									myca[lo].Rolemasterid = s
									fmt.Println(lo, s)
									wgbr.Add(1)
									go func() {
										defer wgbr.Done()
										//var mycpp []models.TblMytree
										var mycpp []models.TtblMytree

										stmts3 := []*dbtran.PipelineStmt{
											dbtran.NewPipelineStmt("select", qry3, &mycpp, rolereq.Companyid, rolereq.Branchid, s),
										}

										_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
											err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts3...)
											return err
										})

										if err != nil {
											fmt.Println("TODO error handlin in go routine")
										}

										createDataTree(&mycpp)
										custmod := mycusr
										allmod := myc
										var ext bool
										//Take Master and impose roleid level values
										for j, t := range mycpp {
											ext = false
											newst := make([]models.TtblMytree, len(mycpp))
											for k, u := range allmod {
												if t.Packid == u.Packid {
													ext = true
													goto JUBR
												}
												JUBR:
												//copy the values from roleid level to master
												if(ext) {
													//Check if the user has access to this module
													for l, v := range custmod {
														if t.Packid == u.Packid {
															//Take the data from master
															newst[j] = t
														}



												}
											}

										}

										for j, t := range mycurcp {
											ext = false
											for k, l := range mycpp {
												if t.Packid == l.Packid {
													ext = true
													goto END
												}
											END:
												if ext {

													remove(mycpp, k)
												}

											}
										}

										myca[lo].Modules = mycpp
									}()
								}
								wgbr.Wait()
*/
