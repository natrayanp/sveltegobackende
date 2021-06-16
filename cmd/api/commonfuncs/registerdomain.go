package commonfuncs

//"github.com/sveltegobackend/cmd/api/models"

/*
func DomRegis1(app *application.Application, w http.ResponseWriter, r *http.Request, dom string) error {
	fmt.Println("----------------- DOM REGIS START -------------------")
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
			SlugCode:   "SESSION-CHKCTX",
			LogMsg:     "Context fetch error",
		}
		dd.HttpRespond()
		return err
	}
	fmt.Println(dom)
	const qry1 = `UPDATE ac.userlogin SET
						companyid = 'CPYID'||nextval('companyid_seq'::regclass),
						lmtime = CURRENT_TIMESTAMP,
						selecthostname = $1,
						hostname = $1
					WHERE userid = $2
					AND siteid = $3
					AND selecthostname isnull`

	var myc1 dbtran.Resultset

	stmts1 := []*dbtran.PipelineStmt{
		dbtran.NewPipelineStmt("update", qry1, &myc1, dom, userinfo.UUID, userinfo.Siteid),
	}

	_, err := dbtran.WithTransaction(ctx, dbtran.TranTypeFullSet, app.DB.Client, nil, func(ctx context.Context, typ dbtran.TranType, db *pgxpool.Pool, ttx dbtran.Transaction) error {
		err := dbtran.RunPipeline(ctx, typ, db, ttx, stmts1...)
		return err
	})

	if err != nil {
		//https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Error())
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				data = "Subdomain or Domain already in use.  Please select a new value"
			}
		}
		fmt.Println(data)
		//		dd := errors.SlugError{
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "DOMAINREG-UPDATE",
			LogMsg:     pgErr.Error(),
		}
		//dd.HttpRespondWithError()
		dd.HttpRespond()
		return err
	}

	type assingrole struct {
		UUID string
		Cpid string
	}
	args, err := json.Marshal(assingrole{UUID: userinfo.UUID, Cpid: userinfo.Companyid})

	if err != nil {
		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "DOMAINREG-ARGSEXTRACT",
			LogMsg:     "json Marsha operation error",
		}
		//dd.HttpRespondWithError()
		dd.HttpRespond()
		return err
	}

	if err := app.Que.Enquejob(&gue.Job{Type: "AssignRole", Args: args}); err != nil {

		dd := httpresponse.SlugResponse{
			Err:        err,
			ErrType:    httpresponse.ErrorTypeDatabase,
			RespWriter: w,
			Request:    r,
			Data:       map[string]interface{}{"message": data},
			SlugCode:   "DOMAINREG-ENQUE",
			LogMsg:     "update error table -> signup failure~default roleassgingment failed~userid" + userinfo.UUID,
		}
		//dd.HttpRespondWithError()
		dd.HttpRespond()
		return err
	}

	fmt.Println("----------------- DOM REGIS END -------------------")
	return nil

}
*/
