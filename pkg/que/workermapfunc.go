package que

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vgarvardt/gue/v2"
)

var wm = gue.WorkMap{
	"PrintName":  printName,
	"AssignRole": assign_role_after_domain_regis,
}

type printNameArgs struct {
	Name string
}

var printName = func(j *gue.Job) error {
	var args printNameArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}
	fmt.Printf("Hello %s!\n", args.Name)
	return nil
}

type assingrole struct {
	UUID string
}

var assign_role_after_domain_regis = func(j *gue.Job) error {
	var args assingrole
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}

	//Apply "SignupAdmin" = 'ROLMA1' role to the user after domain registration
	const qry = `INSERT INTO ac.userrole VALUES ($1,'ROLMA1','PUBLIC','PUBLIC','A','Y',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`

	_, err := j.Tx().Exec(context.Background(), qry, args.UUID)

	fmt.Println(err)
	return nil

}
