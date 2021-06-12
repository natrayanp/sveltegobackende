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
