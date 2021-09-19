package commonfuncs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func GetPacks(app *application.Application, w http.ResponseWriter, r *http.Request) (*[]models.TblCompanyPacks, error) {
	fmt.Println("----------------- PACKAGE CHECK START in GetPacks -------------------")

	var data string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "PACKAGE-CHKCTX")
	if errs != nil {
		return &[]models.TblCompanyPacks{}, errs
	}

	const qry = `SELECT * FROM ac.companypacks 
					WHERE companyid = $1
					AND planid = $2
					AND status in ('A')
					AND startdate <=  CURRENT_DATE
					AND expirydate >= CURRENT_DATE`

	var myc []models.TblCompanyPacks

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.Companyid, "PLANID1"),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
		return err
	})
	fmt.Println("+++++++++++++++++++++$$$end9")
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
		return &[]models.TblCompanyPacks{}, err
	}

	fmt.Println("----------------- PACKAGE CHECK END in GetPacks -------------------")

	return &myc, nil
}

// PackageFetch returns menus and access rights for a user.
// Parameters:
// packfuncid --> If you want only the packsfuncs sent. Send the PACKID form AC.PACKS as array.
//					This forcefully sent only those packs if it exists at company and user level(refer query)
// companyid  --> Always send company id here
//					if Front don't send company id send userinfo.Companyid from calling function
//					else send whatever company id received from front end
// It returns PacksResp struct which is self explanatory and error.
func PackageFetch(app *application.Application, w http.ResponseWriter, r *http.Request, packfuncids []string, companyid string) (*models.PacksResp, error) {
	fmt.Println("----------------- PACKAGE Fetch START -------------------")

	var data string
	var wgbr, wgcp sync.WaitGroup
	var qry string
	//var myc []models.TblMytree
	var myc []models.TtblMytree
	var myca []models.ActiveEntityTree
	var mycacp []models.ActiveEntityTree
	var cmpy *[]models.TblCompany
	var mybr *[]models.TblBranch
	var stmts []*dbtran.PipelineStmt
	var datosend models.PacksResp
	var err error
	havcpydetail := false
	havbrndetail := false
	datosend.Navstring = ""

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "PACKAGE-CHKCTX")
	if errs != nil {
		return &models.PacksResp{}, errs
	}

	datosend.EntityLst = userinfo.Entityid
	datosend.ActiveEntity = ""

	/* Company check start */
	if cmpy, errs = CompanyCheck(app, w, r, companyid); errs != nil {
		fmt.Println("TODO: Error handling")
		return &models.PacksResp{}, errs
	}

	cpyss, errse := getActiveCompany(cmpy, companyid)
	datosend.CompanyLst = *cmpy
	datosend.ActiveCompany = cpyss
	if errse != nil {
		havcpydetail = false
		// If no company exists no point proceeding
		// so return to calling function and handle error there
		// 		this is return error nil to handling the coniditions
		//		for new sign up where no company registered yet
		packfuncids = []string{"PKS8"}
		datosend.Navstring = "ADDCOMPANY"
		//return &datosend, nil
	} else {
		havcpydetail = true
	}
	fmt.Println(havcpydetail)
	/* Company check End */

	/* Branch Check start */
	if packfuncids[0] == "ALL" {
		if mybr, err = BranchCheck(app, w, r, datosend.ActiveCompany.Companyid, []string{"all"}); err != nil {
			fmt.Println("TODO: Error handling")
			return &datosend, errs
		}

		brss, errseb := getActiveBranch(mybr)
		datosend.BranchLst = *mybr
		datosend.ActiveBranch = brss
		if errseb != nil {
			havbrndetail = false
			// If no Branch exists no point proceeding
			// so return to calling function and handle error there
			// 		this is return error nil to handling the coniditions
			//		for new sign up where no branch registered yet
			packfuncids = []string{"PKS8", "PKS9"}
			datosend.Navstring = "ADDBRANCH"
			//	return &datosend, nil
		} else {
			havbrndetail = true
		}
	}
	/* Branch Check end */

	if packfuncids[0] != "ALL" {

		qry = `WITH RECURSIVE MDATA AS 
		(
			SELECT * from ac.ROLE_USER_VIEW where userid = $1 AND companyid = $2 AND packfuncid = ANY($4::varchar[])
				UNION
			SELECT * from ac.ROLE_USER_VIEW where userid = $1 AND companyid = 'PUBLIC' AND packfuncid = ANY($4::varchar[])
			AND (SELECT count(DISTINCT COMPANYID) from ac.ROLE_USER_VIEW where userid = $1 AND companyid = $2) = 0
		),  
		MyTree AS 
		(
			SELECT C.COMPANYID,$3 As branchid,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
			CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
			'SELECTedmodules' AS basketname, false as open
			from ac.COMPANYPACKS_PACKS_VIEW C
			JOIN MDATA A ON A.packfuncid = C.PACKID AND C.menulevel NOT IN ('COMPANY')
				UNION
			SELECT C.COMPANYID,$3 As branchid,T.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
			CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
			'SELECTedmodules' AS basketname, false as open
			from ac.COMPANYPACKS_PACKS_VIEW C
			LEFT JOIN MDATA A ON  A.packfuncid = C.PACKID AND C.menulevel NOT IN ('COMPANY')
			JOIN MyTree AS t ON C.packid = ANY(t.parent)	
		)
		SELECT * FROM MyTree 
		WHERE COMPANYID = $2	
		ORDER BY COMPANYID,SORTORDER,TYPE,NAME;`

		stmts = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("SELECT", qry, &myc, userinfo.UUID, companyid, packfuncids, "ALL"),
		}

		_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
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

		fmt.Println("---------------$$$end6")
		dd, _ := json.Marshal(myc)
		fmt.Println(string(dd))
		fmt.Println("---------------$$$end6a")
		fmt.Printf("&myc is: %p\n", &myc)
		createDataTree(&myc)
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
			fmt.Println("--- start for loop packsops")
			for i, s := range datosend.BranchLst {
				//var mycpp []models.TblMytree
				lo := i
				myca[lo].EntityType = "branch"
				myca[lo].Entityid = s.Branchid
				fmt.Println(lo, s.Branchid)
				wgbr.Add(1)
				go func(brid string) {
					defer wgbr.Done()
					//var mycpp []models.TblMytree
					var mycpp []models.TtblMytree

					qry = `WITH RECURSIVE MDATA AS 
					(
						SELECT * from ac.ROLE_USER_VIEW where userid = $1 AND companyid = $2 AND userbranchacess && ARRAY['ALL'::VARCHAR,$3::VARCHAR]
							UNION
						SELECT * from ac.ROLE_USER_VIEW where userid = $1 AND companyid = 'PUBLIC'
						AND (SELECT count(DISTINCT COMPANYID) from ac.ROLE_USER_VIEW where userid = $1 AND companyid = $2) = 0
					),  
					MyTree AS 
					(
						SELECT C.COMPANYID,$3 As branchid,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
						CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
						'SELECTedmodules' AS basketname, false as open
						from ac.COMPANYPACKS_PACKS_VIEW C
						JOIN MDATA A ON A.packfuncid = C.PACKID AND C.menulevel NOT IN ('COMPANY')
							UNION
						SELECT C.COMPANYID,$3 As branchid,T.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
						CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
						'SELECTedmodules' AS basketname, false as open
						from ac.COMPANYPACKS_PACKS_VIEW C
						LEFT JOIN MDATA A ON  A.packfuncid = C.PACKID AND C.menulevel NOT IN ('COMPANY')
						JOIN MyTree AS t ON C.packid = ANY(t.parent)	
				
					)
					SELECT * FROM MyTree 
					WHERE COMPANYID = $2	
					ORDER BY COMPANYID,SORTORDER,TYPE,NAME;`

					stmts = []*dbtran.PipelineStmt{
						dbtran.NewPipelineStmt("select", qry, &mycpp, userinfo.UUID, companyid, brid),
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
				}(s.Branchid)
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
			datosend.BrnLvlTreeforCpy = []models.ActiveEntityTree{}
		}

		mycacp = make([]models.ActiveEntityTree, len(datosend.CompanyLst))

		for i, sc := range datosend.CompanyLst {
			//var mycppp []models.TblMytree
			mycacp[i].EntityType = "company"
			mycacp[i].Entityid = sc.Companyid

			fmt.Println(i, sc.Companyid)
			fmt.Println("STARTING MY ")
			wgcp.Add(1)
			go func(cpid string) {
				defer wgcp.Done()
				//var mycppp []models.TblMytree
				var mycppp []models.TtblMytree

				qry = `WITH RECURSIVE MDATA AS 
				(
					SELECT * from ac.ROLE_USER_VIEW where userid = $1 AND companyid = $2
						UNION
					SELECT * from ac.ROLE_USER_VIEW where userid = $1 AND companyid = 'PUBLIC'
					AND (SELECT count(DISTINCT COMPANYID) from ac.ROLE_USER_VIEW where userid = $1 AND companyid = $2) = 0
				),  
				MyTree AS 
				(
					SELECT C.COMPANYID,$3 As branchid,A.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
					CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
					'SELECTedmodules' AS basketname, false as open
					from ac.COMPANYPACKS_PACKS_VIEW C
					JOIN MDATA A ON A.packfuncid = C.PACKID AND C.menulevel IN ('COMPANY')
						UNION
					SELECT C.COMPANYID,$3 As branchid,T.ROLEMASTERID,C.packid,c.name,c.displayname,c.description,c.type,c.parent,c.link,c.icon,c.startdate,c.expirydate,c.userrolelimit,c.userlimit,c.branchlimit,c.compstatus,c.sortorder,c.menulevel,c.allowedops,A.allowedopsval,A.USERID,
					CASE WHEN (TRUE = ANY(A.ALLOWEDOPSVAL) IS NULL) AND (C.TYPE = 'function') THEN TRUE ELSE FALSE END AS disablefunc,
					'SELECTedmodules' AS basketname, false as open
					from ac.COMPANYPACKS_PACKS_VIEW C
					LEFT JOIN MDATA A ON  A.packfuncid = C.PACKID AND C.menulevel IN ('COMPANY')
					JOIN MyTree AS t ON C.packid = ANY(t.parent)	
			
				)
				SELECT * FROM MyTree 
				WHERE COMPANYID = $2	
				ORDER BY COMPANYID,SORTORDER,TYPE,NAME;`

				stmts = []*dbtran.PipelineStmt{
					//dbtran.NewPipelineStmt("SELECT", qry, &mycppp, userinfo.UUID, cpid, defaultbranchid),
					dbtran.NewPipelineStmt("select", qry, &mycppp, userinfo.UUID, cpid, "ALL"),
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
			}(sc.Companyid)
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

//func createDataTree(mnodes *[]models.TblMytree) {
func createDataTree(mnodes *[]models.TtblMytree) {
	nodes := *mnodes
	fmt.Printf("&mnodes is: %p\n", mnodes)
	//var newnodes []models.TblMytree
	var newnodes []models.TtblMytree
	fmt.Println("---------------$$$end6a1")
	//m := make(map[pgtype.Varchar]*models.TblMytree)
	m := make(map[string]*models.TtblMytree)
	for i, _ := range nodes {
		//fmt.Printf("Setting m[%d] = <node with ID=%d>\n", n.ID, n.ID)
		m[nodes[i].Packid] = &nodes[i]
	}
	fmt.Println(m)
	fmt.Println("---------------$$$end6a2")
	for i, n := range nodes {
		//fmt.Printf("Setting <node with ID=%d>.Child to <node with ID=%d>\n", n.ID, m[n.ParentID].ID)
		fmt.Println(n)
		//fmt.Println(n.Parent.Dimensions[0].Length)
		fmt.Println(len(n.Parent))
		fmt.Println("---------------$$$end6a2a")
		//if n.Parent.Dimensions[0].Length > 0 {
		if len(n.Parent) > 0 {
			//for _, t := range n.Parent.Elements {
			for _, t := range n.Parent {
				//fmt.Println(t.Status)
				//fmt.Println(t.Status == pgtype.Null)
				fmt.Println("---------------$$$end6a2a1")
				fmt.Println(t)

				if t != nil {
					m[*t].Submenu = append(m[*t].Submenu, m[nodes[i].Packid])
					if len(m[*t].Submenu) > 1 {
						sort.Slice(m[*t].Submenu, func(i, j int) bool {
							return m[*t].Submenu[i].Sortorder < m[*t].Submenu[j].Sortorder
						})
					}
				}

				/*
					if t.Status != pgtype.Null {
						m[t].Submenu = append(m[t].Submenu, m[nodes[i].Id])
					}
				*/
			}
		}
	}
	fmt.Println("---------------$$$end6a3")
	for _, n := range m {
		fmt.Println(n)
		//fmt.Println(n.Parent.Elements[0].Status)
		//fmt.Println(n.Parent.Elements[0].Status == pgtype.Null)
		//if n.Parent.Elements[0].Status == pgtype.Null {
		fmt.Println(n.Parent[0])
		if n.Parent[0] == nil {
			fmt.Println(n)
			fmt.Println(newnodes)
			newnodes = append(newnodes, *n)
			fmt.Println(newnodes)
		}
	}

	sort.Slice(newnodes, func(i, j int) bool {
		return newnodes[i].Sortorder < newnodes[j].Sortorder
	})

	fmt.Println("---------------$$$end6a4")
	fmt.Printf("&mnodes is: %p\n", mnodes)
	fmt.Printf("&newnodes is: %p\n", &newnodes)
	*mnodes = newnodes
	fmt.Printf("&mnodes is: %p\n", mnodes)
	fmt.Println("---------------$$$end6a5")
}

func getActiveCompany(cmpy *[]models.TblCompany, companyid string) (models.TblCompany, error) {

	if len(*cmpy) == 1 {
		if companyid == "" || companyid == (*cmpy)[0].Companyid {
			return (*cmpy)[0], nil
		}
	} else if len(*cmpy) > 1 {

		for _, s := range *cmpy {
			if companyid == "" {
				if s.Isdefault == "Y" {
					return s, nil
				}
			} else {
				if companyid == s.Companyid {
					return (*cmpy)[0], nil
				}
			}
		}
		if companyid == "" {
			return (*cmpy)[0], nil
		}
	}
	return models.TblCompany{}, errors.New("Company/Requested Company setup doesnot exists")
}

func getActiveBranch(mybr *[]models.TblBranch) (models.TblBranch, error) {

	if len(*mybr) == 1 {
		return (*mybr)[0], nil
	} else if len(*mybr) > 1 {
		for _, s := range *mybr {
			if s.Isdefault == "Y" {
				return s, nil
			}
		}
		return (*mybr)[0], nil
	}
	return models.TblBranch{}, errors.New("Branch setup doesnot exists")
}
