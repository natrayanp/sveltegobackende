package commonfuncs

import (
	"fmt"
	"net/http"

	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func fetchUserinfoFromcontext(w http.ResponseWriter, r *http.Request, slugcode string) (*fireauth.User, error) {

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	data := "Technical Error.  Please contact support"

	if !ok {
		err := fmt.Errorf("Empty context")
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   slugcode,
			LogMsg:     "Context fetch error",
		}
		dd.HttpRespond()
		return &fireauth.User{}, err
	}
	return &userinfo, nil

}
