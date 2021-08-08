package login

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func userLogin(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n logintoken Start \n-------------------")
		defer r.Body.Close()

		//Check user registered Start
		//var ctxfetchok bool
		//var userinfo fireauth.User
		//var ctx context.Context
		var data string
		var myc *[]models.TblMytree
		var cppks *[]models.TblCompanyPacks
		var cmpy *[]models.TblCompany
		var brnc *[]models.TblBranch
		var userinfo fireauth.User
		var ctxfetchok bool
		//TODO Add branch
		nxtaction := "DOMAINREGIS"
		havdom := false
		havpacks := false
		havcpydetail := false
		//havbrndetail := false

		//gotolanding := true

		ctx := r.Context()
		userinfo, ctxfetchok = ctx.Value(fireauth.UserContextKey).(fireauth.User)

		fmt.Println(ctxfetchok)

		isregis, errs := commonfuncs.CheckUserRegistered(app, w, r)

		if errs != nil {
			return
		}

		//Check user registered END

		fmt.Println(isregis)
		if isregis {

			/*
				if errs := commonfuncs.SessionOps(app, w, r); errs != nil {
					return
				}
			*/

			/*
				sess, errs := commonfuncs.SessionOps(app, w, r)
				if errs != nil {
					return
				}
				fmt.Println(sess)
				userinfo.Session = sess
				fmt.Println(userinfo)
				r = r.WithContext(context.WithValue(ctx, fireauth.UserContextKey, userinfo))
				userinfo, ctxfetchok = ctx.Value(fireauth.UserContextKey).(fireauth.User)
			*/
			/*
				sess, errs := commonfuncs.SessionOps(app, w, r, &ctx)
				if errs != nil {
					return
				}
				fmt.Println(sess)

					userinfo.Session = sess
					ctx1 := context.WithValue(ctx, fireauth.UserContextKey, userinfo)
					r = r.WithContext(ctx1)
			*/

			fmt.Println("0909090909")
			fmt.Println(userinfo)
			fmt.Println(userinfo.Session)
			fmt.Println("0909090909")

			if !ctxfetchok {
				dd := httpresponse.SlugResponse{
					Err:        errors.New("User Context Fetch error"),
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

			// Check for domain registration
			if userinfo.Companyid != "DEFAULT" {
				havdom = true
				nxtaction = "LANDING"

				//TODO Check for COMPANY PACKS... if none NAV to pricing page

				if cppks, errs = commonfuncs.GetPacks(app, w, r); errs != nil {
					return
				}

				if len(*cppks) > 0 {
					havpacks = true
				} else {
					nxtaction = "ADDPACKS"
					goto NAVCHKEND
				}

				//TODO Check for COMPANY DETAILS CAPTURED... if none NAV to comapny details page

				if cmpy, errs = commonfuncs.CompanyCheck(app, w, r); errs != nil {
					return
				}

				if len(*cmpy) == 1 {
					havcpydetail = true
				} else {
					nxtaction = "ADDCOMPANY"
					goto NAVCHKEND
				}

				//TODO Check for BRANCH DETAILS CAPTURED... if none NAV to branch details page
				if brnc, errs = commonfuncs.BranchCheck(app, w, r); errs != nil {
					return
				}

				if len(*brnc) == 1 {
					//havbrndetail = true
				} else {
					//havbrndetail = false
					nxtaction = "ADDBRANCH"
					goto NAVCHKEND
				}

				//TODO if all the above check satisfied, nav to landing page

			NAVCHKEND:
				//nxtaction = "LANDING"
				switch nxtaction {
				case "LANDING":
					if myc, errs = commonfuncs.PackageFetch(app, w, r, []string{"ALL"}); errs != nil {
						return
					}
				case "ADDCOMPANY":
					if myc, errs = commonfuncs.PackageFetch(app, w, r, []string{"PKS8"}); errs != nil {
						return
					}
					cmpy = &[]models.TblCompany{}
					brnc = &[]models.TblBranch{}

				case "ADDBRANCH":
					if myc, errs = commonfuncs.PackageFetch(app, w, r, []string{"PKS8", "PKS9"}); errs != nil {
						return
					}

				default:
					myc = &[]models.TblMytree{}
					cmpy = &[]models.TblCompany{}
					brnc = &[]models.TblBranch{}
				}

				/*
					if (nxtaction != "ADDPACKS") || (nxtaction != "ADDCOMPANY") {
						//TODO fecth menu tree
						if myc, errs = commonfuncs.PackageFetch(app, w, r, "ALL"); errs != nil {
							return
						}
						//myc = []models.TblMytree{}
					} else if nxtaction == "ADDCOMPANY" {
						if myc, errs = commonfuncs.PackageFetch(app, w, r, "PKS7"); errs != nil {
							return
						}
					} else {
						myc = &[]models.TblMytree{}
					}*/

			} else {
				nxtaction = "DOMAINREGIS"
				fmt.Println("else loop in login tblmytre")
				data = "Subdomain not registered"
				myc = &[]models.TblMytree{}
				cmpy = &[]models.TblCompany{}
				brnc = &[]models.TblBranch{} //Add to send empty branch
			}

			// Return menu

		} else {
			//User registration Start
			fmt.Println("calling regis")
			//gotolanding = false
			data = "Not a Registered user. Register to continue."
			nxtaction = "NOTREGISTERED"
			//EMPTY session id sent
			myc = &[]models.TblMytree{}
			cmpy = &[]models.TblCompany{}
			brnc = &[]models.TblBranch{} //Add to send empty branch
			//User registration End

			cc := httpresponse.SlugResponse{
				Err:        fmt.Errorf("Not a Registered user"),
				ErrType:    httpresponse.ErrorTypeIncorrectInput,
				RespWriter: w,
				Request:    r,
				Data:       map[string]interface{}{"message": data, "nextaction": nxtaction, "menu": myc, "company": cmpy, "branch": brnc, "sessionid": userinfo.Session},
				Status:     "ERROR",
				SlugCode:   "AUTH-USRNOTREG",
				LogMsg:     "User trying to login with non registered user",
				Userinfo:   userinfo,
			}
			cc.HttpRespond()
			//dd.HttpRespondWithError()
			return

		}

		fmt.Println(isregis, havdom, havpacks, havcpydetail)
		fmt.Println(userinfo)
		fmt.Println(userinfo.Session)
		data = "User registration successful."
		//ssd := map[string]interface{}{"message": data, "isregistered": isregis, "havedomain": havdom, "havepacks": havpacks, "havecompany": havcpydetail, "menu": &myc}
		ssd := map[string]interface{}{"message": data, "nextaction": nxtaction, "menu": myc, "company": cmpy, "branch": brnc}
		//&nat{"nat1", "nat2"},
		fmt.Println("registration completed ss sent")
		ss := httpresponse.SlugResponse{
			RespWriter: w,
			Request:    r,
			Data:       ssd,
			Status:     "SUCCESS",
			SlugCode:   "AUTH-RES",
			LogMsg:     "testing",
			Userinfo:   userinfo,
		}
		fmt.Println(userinfo.Session)

		fmt.Println("-------------------\n logintoken Stop \n-------------------")

		ss.HttpRespond()
		//return

		/*
			dds, stat := ss.RespData()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(stat)
			response, _ := json.Marshal(dds)
			fmt.Println(response)

			w.Write(response)
		*/
	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(userLogin(app))
	//return createUser(app)
}
