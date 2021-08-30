package commonfuncs

/*
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

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
func RoleFetch(app *application.Application, w http.ResponseWriter, r *http.Request, packfuncids []string, companyid string) (*models.PacksResp, error) {
	fmt.Println("----------------- ROLE Fetch START -------------------")

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

/*

	if cmpy, errs = CompanyCheck(app, w, r, companyid); errs != nil {
		fmt.Println("TODO: Error handling")
		return &models.PacksResp{}, errs
	}



	qry = `WITH RECURSIVE MyTree AS
	(
		SELECT A.*,false as open,B.roledetailid,B.rolemasterid,B.allowedopsval,'Availablemodules' AS basketname  FROM ac.packs A
		LEFT JOIN ac.roledetails B ON A.packid = B.PACKFUNCID
		WHERE A.packid IN
		(
			(	SELECT PACKFUNCID FROM ac.roledetails
				WHERE rolemasterid IN
					(SELECT DISTINCT rolemasterid FROM ac.userrole
						WHERE userid = $1
						AND status NOT IN ('D','I')
						AND companyid = $2
						AND branchid && ARRAY['ALL'::VARCHAR,$3::VARCHAR]
					)
				INTERSECT
				SELECT PACKFUNCID from ac.companypacks
					WHERE companyid = $2
					AND status NOT IN ('D','I')
					AND startdate <=  CURRENT_DATE
					AND expirydate >= CURRENT_DATE
			)
		)
		AND A.menulevel NOT IN ('COMPANY')
		UNION
		SELECT M.*,false as open,N.roledetailid,N.rolemasterid,N.allowedopsval,'Availablemodules' AS basketname FROM ac.packs M
		LEFT JOIN ac.roledetails N ON M.packid = N.PACKFUNCID
		JOIN MyTree AS t ON M.packid = ANY(t.parent)
			/*SELECT m.*,false as open FROM ac.packs AS m JOIN MyTree AS t ON m.packid = ANY(t.parent)*/
/*
	)
	SELECT * FROM MyTree ORDER BY SORTORDER,TYPE,NAME;`

		stmts = []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, companyid, packfuncids),
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
