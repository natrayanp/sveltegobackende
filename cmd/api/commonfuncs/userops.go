package commonfuncs

import (
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
)

// PackageFetch returns menus and access rights for it for a user.
// Parameters:
// packfuncid --> If you want only the packsfuncs sent. Send the PACKID form AC.PACKS as array.
//					This forcefully sent only those packs if it exists at company and user level(refer query)
// companyid  --> Always send company id here
//					Front don't send company id send userinfo.Companyid from calling side
//					else send whatever company id received from front end
// It returns PacksResp struct which is self explanatory and error.
func UserFetch(app *application.Application, w http.ResponseWriter, r *http.Request, rolereq *models.RoleReq) error {

	qry := `with BR as(
						select branchId,branchName,branchStatus from ac.branch where companyid = 'CPYID23' and branchStatus NOT IN ('D')
					  ),BRF as (
									select * from ac.userprofile a
									left join ac.userlogin b ON a.usrprof_userid = b.userid and a.usrprof_companyid = b.companyid
									cross join BR c
									left join ac.userrole d ON a.usrprof_userid = d.usrrole_userid and a.usrprof_companyid = d.usrrole_companyid AND d.usrrole_status NOT IN ('D') and c.branchid = d.usrrole_userbranchacess 	
									where b.companyid = 'CPYID23' 
					  ), BRS as (
									select sd.usrprof_userid,sd.usrprof_firstname,sd.usrprof_lastname,sd.usrprof_department,sd.usrprof_mobile,sd.usrprof_email,sd.usrrole_status,sd.usrprof_companyid from BRF sd
					  ), BRT as (
									select sd.usrprof_userid, sd.usrprof_companyid, json_agg(sd) from BRS sd GROUP BY sd.usrprof_userid,sd.usrprof_companyid
					  ), BRFO as (
						  			select * from BRT a
									left join BRS b  on a.usrprof_userid = b.usrprof_userid AND a.usrprof_companyid = b.usrprof_companyid
					  ) select * from BRFO;`

	fmt.Println(qry)
	err := fmt.Errorf("nil")

	if err != nil {
		return err
	}
	return err
}
