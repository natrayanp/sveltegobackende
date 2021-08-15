package commonfuncs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func GetPacks(app *application.Application, w http.ResponseWriter, r *http.Request) (*[]models.TblCompanyPacks, error) {
	fmt.Println("----------------- PACKAGE CHECK START -------------------")

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

	fmt.Println("----------------- PACKAGE CHECK END -------------------")

	return &myc, nil
}

func PackageFetch(app *application.Application, w http.ResponseWriter, r *http.Request, packfuncid []string, companyid string) (*[]models.TblMytree, error) {
	fmt.Println("----------------- PACKAGE Fetch START -------------------")

	var data string

	ctx := r.Context()
	//userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	userinfo, errs := FetchUserinfoFromcontext(w, r, "PACKAGE-CHKCTX")
	if errs != nil {
		return &[]models.TblMytree{}, errs
	}

	var qry string
	var myc []models.TblMytree
	var myca [][]models.TblMytree
	var mybranch *[]models.TblBranch
	var stmts []*dbtran.PipelineStmt
	var datosend models.PacksResp
	var err error
	var havbrndetail bool

	datosend.EntityLst = userinfo.Entityid
	datosend.ActiveEntity = ""
	datosend.CompanyLst = []string{companyid}
	datosend.ActiveCompany = datosend.CompanyLst[0]

	mybr, err := BranchCheck(app, w, r, []string{"All"})
	if err != nil {
		fmt.Println("TODO: Error handling")
	}

	if mybr, errs := commonfuncs.BranchCheck(app, w, r, []string{"all"}); errs != nil {
		return
	}

	if len(*mybr) > 0 {
		havbrndetail = true
		datosend.BranchLst = *mybr
		datosend.ActiveBranch = datosend.BranchLst[0]
	} else {
		havbrndetail = false
		datosend.BranchLst = *mybr
		datosend.ActiveBranch = ""
	}

	/*
		if packfuncid[0] == "ALL" {


		qry = `WITH RECURSIVE MyTree AS
					(
						SELECT *,false as open FROM ac.packs WHERE id IN
						(
							(	SELECT PACKFpackfuncidemasterid IN
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
						UNION
						SELECT m.*,false as open FROM ac.packs AS m JOIN MyTree AS t ON m.id = ANY(t.parent)
					)
					SELECT * FROM MyTree ORDER BY TYPE, SORTORDER,NAME;`

			stmts = []*dbtran.PipelineStmt{
				dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Companyid),
			}


		} else {
	*/
	if packfuncid[0] != "ALL" {
		qry = `WITH RECURSIVE MyTree AS 
				(
					SELECT *,false as open FROM ac.packs WHERE id IN
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
							INTERSECT	
							SELECT * FROM UNNEST($3) AS PACKFUNCID
						)
					)
					UNION
					SELECT m.*,false as open FROM ac.packs AS m JOIN MyTree AS t ON m.id = ANY(t.parent)
				)
				SELECT * FROM MyTree ORDER BY TYPE, SORTORDER,NAME;`

		stmts = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Companyid, packfuncid),
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
			return &[]models.TblMytree{}, err
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

		return &myc, nil

	} else {

		if havbrndetail {

			for i, s := range datosend.BranchLst {
				fmt.Println(i, s.Branchid.String)

				qry = `WITH RECURSIVE MyTree AS 
			(
				SELECT $3 as branchid, *,false as open FROM ac.packs WHERE id IN
				(
					(	SELECT PACKFpackfuncidemasterid IN 
							(SELECT DISTINCT rolemasterid FROM ac.userrole 
								WHERE userid = $1
								AND status NOT IN ('D','I') 
								AND companyid = $2
								AND branchid = $3
							)
						INTERSECT	
						SELECT PACKFUNCID from ac.companypacks 
							WHERE companyid = $2
							AND status NOT IN ('D','I')
							AND startdate <=  CURRENT_DATE
							AND expirydate >= CURRENT_DATE
					)
				)
				UNION
				SELECT $3 as branchid,m.*,false as open FROM ac.packs AS m JOIN MyTree AS t ON m.id = ANY(t.parent)
			)
			SELECT * FROM MyTree ORDER BY TYPE, SORTORDER,NAME;`

				append(stmts, []*dbtran.PipelineStmt{
					dbtran.NewPipelineStmt("select", qry, &myca[i], userinfo.UUID, userinfo.Companyid, s.Branchid.String),
				})

			}

		}

	}
}

func createDataTree(mnodes *[]models.TblMytree) {
	nodes := *mnodes
	fmt.Printf("&mnodes is: %p\n", mnodes)
	var newnodes []models.TblMytree
	fmt.Println("---------------$$$end6a1")
	m := make(map[pgtype.Varchar]*models.TblMytree)
	for i, _ := range nodes {
		//fmt.Printf("Setting m[%d] = <node with ID=%d>\n", n.ID, n.ID)
		m[nodes[i].Id] = &nodes[i]
	}
	fmt.Println(m)
	fmt.Println("---------------$$$end6a2")
	for i, n := range nodes {
		//fmt.Printf("Setting <node with ID=%d>.Child to <node with ID=%d>\n", n.ID, m[n.ParentID].ID)
		fmt.Println(n)
		fmt.Println(n.Parent.Dimensions[0].Length)
		fmt.Println("---------------$$$end6a2a")
		if n.Parent.Dimensions[0].Length > 0 {
			fmt.Println(n.Parent.Status)

			for _, t := range n.Parent.Elements {
				fmt.Println(t.Status)
				fmt.Println(t.Status == pgtype.Null)
				fmt.Println("---------------$$$end6a2a1")
				if t.Status != pgtype.Null {
					m[t].Submenu = append(m[t].Submenu, m[nodes[i].Id])
				}
			}
		}
	}
	fmt.Println("---------------$$$end6a3")
	for _, n := range m {
		fmt.Println(n)
		fmt.Println(n.Parent.Elements[0].Status)
		fmt.Println(n.Parent.Elements[0].Status == pgtype.Null)
		if n.Parent.Elements[0].Status == pgtype.Null {
			fmt.Println(n)
			fmt.Println(newnodes)
			newnodes = append(newnodes, *n)
			fmt.Println(newnodes)
		}
	}
	fmt.Println("---------------$$$end6a4")
	fmt.Printf("&mnodes is: %p\n", mnodes)
	fmt.Printf("&newnodes is: %p\n", &newnodes)
	*mnodes = newnodes
	fmt.Printf("&mnodes is: %p\n", mnodes)
	fmt.Println("---------------$$$end6a5")
}
