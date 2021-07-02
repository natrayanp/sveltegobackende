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
