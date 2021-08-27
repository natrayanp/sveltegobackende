package models

type DomainRegis struct {
	Siteid    string
	Registype string
}

type DrResult struct {
	Domainmapid string
}

type QueConfig struct {
	WorkerEnabled bool
	WorkerCount   int64
	QueName       string
}

type PackSelect struct {
	Planid string
}

type RefDatReq struct {
	Reftype string // valid values are group, single
	Refname string // menu name (eg: Company regist) or individual item name (eg: country)
}

type RefDatReqFinal struct {
	Refs      []RefDatReq
	RefResult map[string]interface{}
}

/*
type TtblRefdata struct {
	Id        string         `json:"id"`
	Refvalcat string         `json:"refvalcat"`
	Refvalue  string         `json:"refvalue"`
	Sortorder int            `json:"sortorder"`
	Parent    []*string      `json:"parent"`
	Submenu   []*TtblRefdata `json:"submenu"`
}
*/
type Cpy struct {
	CompanyId          string
	CompanyName        string
	CompanyShortName   string
	CompanyCategory    string
	CompanyStatus      string
	CompanyLogoUrl     string
	CompanyLogo        string
	CompanyIndustry    string
	CompanyTaxID       string
	CompanyStartDate   string
	CompanyAddLine1    string
	CompanyAddLine2    string
	CompanyCountry     string
	CompanyCity        string
	CompanyState       string
	CompanyPinCode     string
	CompanyPhone       string
	CompanyFax         string
	CompanyMobile      string
	CompanyEmail       string
	CompanyWebsite     string
	CompanyFiscalYear  string
	CompanyTimeZone    string
	CompanyBaseCurency string
	CompanysParent     string
	Entityid           string
}

type Cpyops struct {
	Optype string
}

type ResultCount struct {
	Count int
}

type Brn struct {
	CompanyId         string
	Companyname       string
	BranchId          string
	BranchName        string
	BranchShortName   string
	BranchCategory    string
	BranchStatus      string
	BranchDescription string
	BranchImageUrl    string
	BranchAddLine1    string
	BranchAddLine2    string
	BranchCity        string
	BranchState       string
	BranchCountry     string
	BranchPinCode     string
	BranchPhone       string
	BranchFax         string
	BranchMobile      string
	BranchWebsite     string
	BranchEmail       string
	BranchStartDate   string
	Isdefault         string
}

type BrnResp struct {
	Optype     string
	Branchdata TblBranch
}

type RegisChk struct {
	Isregis      bool
	Companyowner string
}

type ActiveEntityTree struct {
	EntityType string
	Entityid   string
	EntityTree []TtblMytree
}

type ReqEntityTree struct {
	EntityType string
	Entityid   []string
	EntityTree []TtblMytree
}

type PacksResp struct {
	Navstring        string
	EntityLst        []*string
	ActiveEntity     string
	CompanyLst       []TblCompany
	ActiveCompany    TblCompany
	BranchLst        []TblBranch
	ActiveBranch     TblBranch
	BrnLvlTreeforCpy []ActiveEntityTree
	CpyLvlTreeforCpy []ActiveEntityTree
}
