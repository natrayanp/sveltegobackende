package login

import (
	"fmt"
	"net/http"
    "encoding/json"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func registerdomain(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n registerdomain Start \n-------------------")
		defer r.Body.Close()
		ctx := r.Context()

		var data string
		var myc []models.TblMytree
		var p models.DomainRegis
		fmt.Println(r.Body)
		isregis:=false
		havdom:=false
		
		err := json.NewDecoder(r.Body).Decode(&p)
		fmt.Println(p)
		

		if err != nil {
			return
		}

		//Check user registered END
		userinfo, ctxfetchok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

		if !ctxfetchok {
			dd := httpresponse.SlugResponse{
				Err:        fmt.Errorf("User Context Fetch error"),
				ErrType:    httpresponse.ErrorTypeContexFetchFail,
				RespWriter: w,
				Request:    r,
				Data:       map[string]interface{}{"message": "Technical issue. Please contact support"},
				SlugCode:   "AUTH-USRREGCHKFAIL",
				LogMsg:     "Logic Failed",
			}
			dd.HttpRespond()
			//dd.HttpRespondWithError()
			return
		}
		var hostname string
		fmt.Println(p.Registype)
		fmt.Println(p.Registype == "subdomain")
		if (p.Registype == "subdomain") {
			// Check for domain registration
			hostname =p.Siteid+"."+userinfo.Hostname
			fmt.Println("calling regis: ",hostname)
		} else {
			//User registration Start			
			hostname =p.Siteid
			fmt.Println("calling regis else: ",hostname)
		}

		commonfuncs.DomRegis(app,w,r,hostname)





		ssd := map[string]interface{}{"message": data, "isregistered": isregis, "havedomain": havdom, "menu": &myc}
		//&nat{"nat1", "nat2"},
		fmt.Println("registration completed ss sent")
		ss := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     "SUCCESS",
			SlugCode:   "AUTH-RES",
			LogMsg:     "testing",
		}

		ss.HttpRespond()
		fmt.Println("-------------------\n registerdomain Stop \n-------------------")
		return

	}
}

func DoRegisDomain(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(registerdomain(app))
	//return createUser(app)
}
