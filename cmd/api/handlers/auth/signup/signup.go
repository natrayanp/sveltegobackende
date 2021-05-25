package signup

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sveltegobackend/cmd/api/commonfuncs"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/errors"
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

		isregis := commonfuncs.CheckUserRegistered(app, w, r)
		//Check user registered END

		userinfo, ctxfetchok = ctx.Value(fireauth.UserContextKey).(fireauth.User)

		if !ctxfetchok {
			dd := errors.SlugError{
				ErrType:    errors.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Slug:       "Technical issue. Please contact support",
				SlugCode:   "AUTH-USRREGCHKFAIL",
				LogMsg:     "Logic Failed",
			}
			dd.HttpRespondWithError()
		}
		fmt.Println(isregis)
		if isregis {
			data = userinfo.Email + "Already a registered user"
			if !userinfo.EmailVerified {
				data = data + ". Verify your email before login."
			}
		} else {
			//User registration Start
			registersuccess := commonfuncs.RegisterUser(app, w, r)
			if !registersuccess {
				dd := errors.SlugError{
					ErrType:    errors.ErrorTypeDatabase,
					RespWriter: w,
					Request:    r,
					Slug:       "User Registration Failed. Please contact support",
					SlugCode:   "AUTH-USRREGFAIL",
					LogMsg:     "Logic Failed",
				}
				dd.HttpRespondWithError()
			}
			data = "Registration successful for " + userinfo.Email + ". Verify your email before login."
			//User registration End
		}

		myd := httpresponse.SucessResponse{Data: data, Respcode: ""}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(myd)
		fmt.Println(response)
		fmt.Println("-------------------\n signuptoken Stop \n-------------------")
		w.Write(response)
	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(userSignup(app))
	//return createUser(app)
}
