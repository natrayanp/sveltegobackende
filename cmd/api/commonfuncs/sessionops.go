package commonfuncs

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func SessionOps(app *application.Application, w http.ResponseWriter, r *http.Request, userinfo *fireauth.User) error {
	//Check user registered Start

	ctx := r.Context()
	fmt.Println("--------- inside session ops ------------")
	fmt.Println((*userinfo).Session)

	if userinfo.Session == "" {

		qry := `UPDATE ac.loginh SET logoutime = CURRENT_TIMESTAMP
			WHERE userid = $1
			AND companyid = $2
			AND logoutime IS NULL`

		qry1 := `INSERT INTO ac.loginh (userid, ipaddress, sessionid, companyid, logintime) 
			VALUES ($1, $2, $3,$4 ,CURRENT_TIMESTAMP) RETURNING *;`

		currentTime := time.Now().String()
		mysess := getHash(userinfo.UUID+currentTime, "")
		fmt.Println("----------------printing session------------------")
		fmt.Println(mysess)

		var myc dbtran.Resultset
		var myc1 dbtran.Resultset

		stmts := []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("update", qry, &myc, userinfo.UUID, userinfo.Companyid),
			dbtran.NewPipelineStmt("insert", qry1, &myc1, userinfo.UUID, "", mysess, userinfo.Companyid),
		}

		_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
			err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
			return err
		})

		if err != nil {
			//		dd := errors.SlugError{
			dd := httpresponse.SlugResponse{
				Err:        err,
				ErrType:    httpresponse.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Data:       map[string]interface{}{"message": "Database error"},
				SlugCode:   "SESSION-CREAT",
				LogMsg:     "Database error",
			}
			//dd.HttpRespondWithError()
			dd.HttpRespond()
			userinfo.Session = ""
			return err
		} else {
			userinfo.Session = mysess
		}
	} else {

		//TODO check if session exists is it valid if not return error
		fmt.Println("else loop")

		qry3 := `SELECT sessionid FROM ac.loginh 
		WHERE userid = $1
		AND companyid = $2
		AND logoutime IS NULL
		AND sessionid = $3`

		//var myc3 []models.ResultCount

		type SessionResult struct {
			Sessionid *string
		}

		var myc3 []SessionResult

		stmts := []*dbtran.PipelineStmt{
			dbtran.NewPipelineStmt("select", qry3, &myc3, userinfo.UUID, userinfo.Companyid, userinfo.Session),
		}

		_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
			err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts...)
			return err
		})
		res := false
		if err != nil {
			res = true
		}
		if len(myc3) == 1 {
			res = true

			if *(myc3[0].Sessionid) == userinfo.Session {
				res = false
			}
		} else {
			res = true
		}

		if res {
			fmt.Println("going inside error")
			var ds map[string]interface{}
			var lms string
			var sc string

			if err != nil {
				fmt.Println("going inside s")
				ds = map[string]interface{}{"message": "Database error from"}
				lms = "Database error from"
				sc = "SESSION-CHECK"
			} else {
				fmt.Println("going inside e")
				err = fmt.Errorf("user session fetch error")
				ds = map[string]interface{}{"message": "Invalid session"}
				lms = "Invalid session"
				sc = "SESSION-INVALID"
			}

			dd := httpresponse.SlugResponse{
				Err:        err,
				ErrType:    httpresponse.ErrorTypeDatabase,
				RespWriter: w,
				Request:    r,
				Data:       ds, //map[string]interface{}{"message": "Database error"},
				SlugCode:   sc,
				LogMsg:     lms, //"Database error",
			}
			//dd.HttpRespondWithError()
			dd.HttpRespond()
			return err
		}
		fmt.Println("going insiderrrrrr error")

	}

	/*
		ctx1 := context.WithValue(ctx, fireauth.GetUserCtxKey(), userinfo)
		r = r.WithContext(ctx1)
		r.WithContext(context.WithValue(r.Context(), fireauth.GetUserCtxKey(), userinfo))
	*/
	fmt.Println("++++++++++++ddsds++++++++++++++")

	return nil
}

func getHash(data string, secret string) string {
	//sum := sha256.Sum256([]byte("hello world\n"))
	if secret == "" {
		secret = "sesstkn"
	}

	//data := "data"
	fmt.Printf("Secret: %s Data: %s\n", secret, data)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	//sha := hex.EncodeToString(h.Sum(nil))

	//sum := sha256.Sum256([]byte(val))
	//var bv []byte = sum[:]

	bv := h.Sum(nil)

	hasher := sha1.New()
	hasher.Write(bv)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
