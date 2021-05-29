package signup

import (
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func userSignup(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n signuptoken Start \n-------------------")
		defer r.Body.Close()
		ctx := r.Context()

		//Check user registered Start
		var ctxfetchok bool
		var userinfo fireauth.User
		var data string

		isregis, errs := commonfuncs.CheckUserRegistered(app, w, r)

		if errs != nil {
			return
		}

		//Check user registered END
		userinfo, ctxfetchok = ctx.Value(fireauth.UserContextKey).(fireauth.User)

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

		fmt.Println(isregis)
		if isregis {
			data = userinfo.Email + "Already a registered user"
			if !userinfo.EmailVerified {
				data = data + ". Verify your email before login."
			}
		} else {
			//User registration Start
			fmt.Println("calling regis")

			registersuccess, err := commonfuncs.RegisterUser(app, w, r)

			if err != nil {
				return
			}

			fmt.Println("registration completed")

			if !registersuccess {
				dd := httpresponse.SlugResponse{
					ErrType:    httpresponse.ErrorTypeDatabase,
					RespWriter: w,
					Request:    r,
					Userinfo:   userinfo,
					Data:       map[string]interface{}{"message": "User Registration Failed. Please contact support"},
					SlugCode:   "AUTH-USRREGFAIL",
					LogMsg:     "Logic Failed",
				}
				dd.HttpRespond()
				return
			}
			data = "Registration successful for " + userinfo.Email + ". Verify your email before login."
			//User registration End
		}

		ssd := map[string]interface{}{"message": data}
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
		return

		/*
			dds, stat := ss.RespData()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(stat)
			response, _ := json.Marshal(dds)
			fmt.Println(response)
			fmt.Println("-------------------\n signuptoken Stop \n-------------------")
			w.Write(response)
		*/
	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(userSignup(app))
	//return createUser(app)
}
