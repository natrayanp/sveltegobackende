package models

type TblUserlogin struct {
	Userid          string `json:"userid"`
	Username        string `json:"username"`
	Useremail       string `json:"useremail"`
	Userpassword    string `json:"userpassword"`
	Userstatus      string `json:"userstatus"`
	Emailverified   bool   `json:"emailverified"`
	Siteid          string `json:"siteid"`
	Domainmapid     string `json:"domainmapid"`
	Userstatlstupdt string `json:"userstatuslastupdate"`
	Octime          string `json:"creattime"`
	Lmtime          string `json:"lasmodifytime"`
}
