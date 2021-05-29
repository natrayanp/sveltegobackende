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
	"github.com/sveltegobackend/cmd/api/models"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/db/dbtran"
	"github.com/sveltegobackend/pkg/fireauth"
	"github.com/sveltegobackend/pkg/httpresponse"
)

func SessionOperation(app *application.Application, w http.ResponseWriter, r *http.Request) (bool, error) {
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
		return false, err
	}

	qry := `SELECT * FROM ac.userlogin
			WHERE userid = $1
			AND  siteid = $2;`

	const qry = `UPDATE ac.loginh SET logoutime = CURRENT_TIMESTAMP
			WHERE userid = $1
			AND siteid = $2
			AND logoutime IS NULL`

	currentTime := time.Now().String()
	mysess := getHash(userinfo.UUID+currentTime, "")
	fmt.Println(mysess)

	var myc []models.TblUserlogin

	stmts := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("select", qry, &myc, userinfo.UUID, userinfo.Siteid),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeNoTran, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
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
			SlugCode:   "AUTH-INT",
			LogMsg:     "Database error",
		}
		//dd.HttpRespondWithError()
		dd.HttpRespond()
		return false, err
	}
	fmt.Println("ddsds")
	fmt.Println(myc)

	if len(myc) > 1 {
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": "Invalid Company Profile Setup Exists.  Contact Support"},
			SlugCode:   "AUTH-NOMULCPY",
			LogMsg:     "Company Details Not set or Have multiple Company; sql:" + qry,
		}

		dd.HttpRespond()
		return false, fmt.Errorf("Invalid Company Profile Setup Exists.  Contact Support")
	} else if len(myc) == 0 {
		fmt.Println("no record db success")
		return false, nil
	}

	//Check user registered end
	return true, nil
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
