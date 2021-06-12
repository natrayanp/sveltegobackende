package models

import "github.com/jackc/pgtype"

type TblUserlogin struct {
	Userid          pgtype.Varchar     `json:"userid"`
	Username        pgtype.Varchar     `json:"username"`
	Useremail       pgtype.Varchar     `json:"useremail"`
	Userpassword    pgtype.Varchar     `json:"userpassword"`
	Userstatus      pgtype.Varchar     `json:"userstatus"`
	Emailverified   pgtype.Bool        `json:"emailverified"`
	Siteid          pgtype.Varchar     `json:"siteid"`
	Hostname		pgtype.Text	     `json:"hostname"`
	Companyid     pgtype.Varchar     `json:"companyid"`
	Userstatlstupdt pgtype.Timestamptz `json:"userstatuslastupdate"`
	Octime          pgtype.Timestamptz `json:"creattime"`
	Lmtime          pgtype.Timestamptz `json:"lasmodifytime"`
}

type TblMytree struct {
	Id          pgtype.Varchar     `json:"id"`
	Name        pgtype.Varchar     `json:"name"`
	DisplayName pgtype.Varchar     `json:"displayname"`
	Description pgtype.Varchar     `json:"description"`
	Type        pgtype.Varchar     `json:"type"`
	Paren       pgtype.Text        `json:"parent"`
	Link        pgtype.Varchar     `json:"link"`
	Icon        pgtype.Varchar     `json:"icon"`
	Status      pgtype.Varchar     `json:"status"`
	Octime      pgtype.Timestamptz `json:"creattime"`
	Lmtime      pgtype.Timestamptz `json:"lasmodifytime"`
	Open        pgtype.Bool        `json:"open"`
}



