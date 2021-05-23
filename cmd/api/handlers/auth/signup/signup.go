package signup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/errors"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/mymiddleware"
)

func userSignup(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------\n signuptoken Start \n-------------------")
		defer r.Body.Close()
		ctx := r.Context()
		userinfo, err := fireauth.UserFromCtxs(ctx)

		//Check user registered Start

		qry := `SELECT * FROM ac.userlogins
		WHERE userid = $1
		AND  siteid = $2;`

		var myc []models.TblUserlogin

		stmts := []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Siteid),
		}

		_, err = dbtran.WithTransaction(ctx, dbtran.TranTypeNoTran, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
			err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
			return err
		})

		if err != nil {
			dd := errors.SlugError{
				Err:        err,
				ErrType:    errors.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Slug:       "Database error",
				SlugCode:   "INT",
				LogMsg:     "Database error",
			}
			dd.HttpRespondWithError()
			return
		}

		fmt.Println(myc)

		if len(myc) == 0 || len(myc) > 1 {
			dd := errors.SlugError{
				Err:        err,
				ErrType:    errors.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Slug:       "Invalid Company Profile Setup Exists.  Contact Support",
				SlugCode:   "NOMULCPY",
				LogMsg:     "Company Details Not set or Have multiple Company; sql:" + qry,
			}
			dd.HttpRespondWithError()
			return
		}

		//Check user registered end

		user := &models.User{}
		json.NewDecoder(r.Body).Decode(user)

		if err := user.Create(r.Context(), app); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Oops")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(user)
		fmt.Println("-------------------\n signuptoken Stop \n-------------------")
		w.Write(response)
	}
}

func Do(app *application.Application) http.HandlerFunc {
	return mymiddleware.Chain(userSignup(app))
	//return createUser(app)
}
