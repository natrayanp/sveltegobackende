package commonfuncs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mitchellh/mapstructure"
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

	qry = `SELECT c.COMPANYID,'PUBLIC' AS BRANCHID,'' AS Roledetailid,'' AS ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,
			c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,
			array_fill(FALSE, ARRAY[array_length(c.allowedops,1)])  AS allowedopsval,$2 as userid,
			CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
			'Availablemodules' AS basketname, false as open
			FROM ac.COMPANYPACKS_PACKS_VIEW C
			LEFT JOIN ac.ROLE_USER_VIEW B ON C.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = $2
			WHERE C.COMPANYID = $1;`

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &availmod, rolereq.Companyid, userinfo.UUID),
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

	//This will give all the rolewise details for the company if Roles are already created -- END

	//Fetch all Roles available for the company -- START

	//var selmod []models.RoleSelectModu
	//var selmod []models.TtblMytree
	//var selmod []string
	var selmod []models.TmpRoleSelectModu
	var fselmod []models.RoleSelectModu

	qry = `WITH MYVA AS (
				SELECT c.COMPANYID,A.BRANCHID,B.Roledetailid,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,
				c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,B.allowedopsval,B.USERID,
				CASE WHEN (TRUE = ANY(B.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
				'selectedmodules' AS basketname, false as open
				from ac.COMPANYPACKS_PACKS_VIEW C
				LEFT JOIN ac.rolemaster A ON A.COMPANYID = $1
				LEFT JOIN ac.ROLE_USER_VIEW B ON A.COMPANYID = B.COMPANYID AND  B.packfuncid = C.PACKID AND B.USERID = $2 AND B.ROLEMASTERID = A.ROLEMASTERID
				WHERE C.COMPANYID = $1
				ORDER BY A.ROLEMASTERID
			) , MYVAGROUP AS 
			( 
				SELECT rolemasterid,json_agg(MYVA) AS modules FROM myva GROUP BY MYVA.rolemasterid
			 ) SELECT X.*,Y.displayname FROM myvagroup X
			   LEFT JOIN ac.rolemaster Y ON X.rolemasterid = Y.rolemasterid;`

	stmts = []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &selmod, rolereq.Companyid, userinfo.UUID),
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
		fselmod[i].Displayname = s.Displayname
		mapstructure.Decode(s.Modules, &fselmod[i].Modules)
		createDataTree(&fselmod[i].Modules)
	}

	datosend.Selectedmodules = fselmod
	return &datosend, err

}

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
																)
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
