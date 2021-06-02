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

func SessionOps(app *application.Application, w http.ResponseWriter, r *http.Request) error {
	//Check user registered Start

	ctx := r.Context()
	userinfo, ok := ctx.Value(fireauth.UserContextKey).(fireauth.User)

	if !ok {
		err := fmt.Errorf("Empty context")
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": "Technical Error.  Please contact support"},
			SlugCode:   "SESSION-CHKCTX",
			LogMsg:     "Context fetch error",
		}
		dd.HttpRespond()
		return err
	}

	qry := `UPDATE ac.loginh SET logoutime = CURRENT_TIMESTAMP
			WHERE userid = $1
			AND companyid = $2
			AND logoutime IS NULL`

	qry1 := `INSERT INTO ac.loginh (userid, ipaddress, sessionid, companyid, logintime) 
			VALUES ($1, $2, $3,$4 ,CURRENT_TIMESTAMP) RETURNING *;`

	currentTime := time.Now().String()
	mysess := getHash(userinfo.UUID+currentTime, "")
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
		return err
	}

	userinfo.Session = mysess
	ctx1 := context.WithValue(ctx, fireauth.GetUserCtxKey(), userinfo)
	r = r.WithContext(ctx1)

	fmt.Println("ddsds")
	fmt.Println(myc)

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
