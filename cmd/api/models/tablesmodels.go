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
	Domainmapid     pgtype.Varchar     `json:"domainmapid"`
	Userstatlstupdt pgtype.Timestamptz `json:"userstatuslastupdate"`
	Octime          pgtype.Timestamptz `json:"creattime"`
	Lmtime          pgtype.Timestamptz `json:"lasmodifytime"`
}
