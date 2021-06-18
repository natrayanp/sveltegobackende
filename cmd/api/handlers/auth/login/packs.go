package login

/*
import (
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func getPacks(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n getPacks Start \n-------------------")
		defer r.Body.Close()

		var data string
		//var myc []models.TblMytree
		fmt.Println(r.Body)
		//isregis := false
		//havdom := false

		//Check user registered END
		userinfo, errs := commonfuncs.FetchUserinfoFromcontext(w, r, "PACKAGE-CHKCTX")
		if errs != nil {
			return
		}
		var myc *[]models.TblMytree
		myc, errs = commonfuncs.PackageFetch(app, w, r)
		if errs != nil {
			return
		}

		data = "Your domain registration successful. Login with your url - " + hostname
		ssd := map[string]interface{}{"message": data}
		//&nat{"nat1", "nat2"},
		fmt.Println("domain registration completed ss sent")
		ss := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     "SUCCESS",
			SlugCode:   "DOMAIN-REG",
			LogMsg:     "testing",
		}

		ss.HttpRespond()
		fmt.Println("-------------------\n getPacks Stop \n-------------------")
		return

	}
}

func DoPacks(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(getPacks(app))
	//return createUser(app)
}
*/
