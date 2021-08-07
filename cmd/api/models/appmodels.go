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
