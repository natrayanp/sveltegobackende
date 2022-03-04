package que

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vgarvardt/gue/v2"
)

var wm = gue.WorkMap{
	"PrintName":  printName,
	"AssignRole": assign_role_after_domain_regis,
	"Auditentry": audit_entry,
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
	Cpid string
}

var assign_role_after_domain_regis = func(j *gue.Job) error {
	var args assingrole
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}

	//Apply "SignupAdmin" = 'ROLMA1' role to the user after domain registration
	//const qry = `INSERT INTO ac.userrole VALUES ($1,'ROLMA1',$2,ARRAY['ALL'],'A','Y',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`
	const qry = `INSERT INTO ac.userrole VALUES ($1,ARRAY['ROLMA1'],$2,'ALL','A','Y',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`

	_, err := j.Tx().Exec(context.Background(), qry, args.UUID, args.Cpid)

	//TODO: write logic to populate the company packs for default modules
	//like fetch all packs from PACKS1 and loop through it and insert one by one
	/*
		const qry1 = `INSERT INTO ac.companypacks VALUES ($1,'PKS8',CURRENT_DATE,CURRENT_DATE,10,10,10,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`
		const qry2 = `INSERT INTO ac.companypacks VALUES ($1,'PKS9',CURRENT_DATE,CURRENT_DATE,10,10,10,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`
		const qry3 = `INSERT INTO ac.companypacks VALUES ($1,'PKS14',CURRENT_DATE,CURRENT_DATE,10,10,10,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`

		_, err = j.Tx().Exec(context.Background(), qry1, args.Cpid)
		_, err = j.Tx().Exec(context.Background(), qry2, args.Cpid)
		_, err = j.Tx().Exec(context.Background(), qry3, args.Cpid)

		fmt.Println(err)

	*/
	fmt.Println(err)
	return nil

}

type auditentryargs struct {
	Itemid    string
	Itemkeys  interface{}
	Action    string
	Oldval    interface{}
	Newval    interface{}
	Companyid string
	User      string
	Time      time.Time
}

var audit_entry = func(j *gue.Job) error {
	fmt.Println("---------------- AUDIT STARTS ------------------")
	fmt.Println("i am printing inside audit entry")
	var args auditentryargs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		fmt.Println(")))))))))))))")
		fmt.Println(err)
		return err
	}

	fmt.Println(args)
	const qry = `INSERT INTO ac.audit VALUES (DEFAULT,$1,$2,$3,$4,$5,$6,$7,$8)`
	ol := fmt.Sprintf("%v", args.Oldval)
	ne := fmt.Sprintf("%v", args.Newval)
	fmt.Println(ol)
	fmt.Println(ne)
	//_, err := j.Tx().Exec(context.Background(), qry, args.Itemid, args.Action, args.Oldval, args.Newval, args.User, args.Time)
	_, err := j.Tx().Exec(context.Background(), qry, args.Itemid, args.Itemkeys, args.Action, args.Oldval, args.Newval, args.Companyid, args.User, args.Time)
	fmt.Println(")))))))))))))+++++++++++++++")
	fmt.Println(err)
	fmt.Println("---------------- AUDIT ENDS ------------------")
	return nil
}
