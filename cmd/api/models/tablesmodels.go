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
	Hostname        pgtype.Text        `json:"hostname"`
	Selecthostname  pgtype.Text        `json:"selecthostname"`
	Companyid       pgtype.Varchar     `json:"companyid"`
	Userstatlstupdt pgtype.Timestamptz `json:"userstatuslastupdate"`
	Octime          pgtype.Timestamptz `json:"creattime"`
	Lmtime          pgtype.Timestamptz `json:"lasmodifytime"`
}

type TblMytree struct {
	Id          pgtype.Varchar     `json:"id"`
	Name        pgtype.Varchar     `json:"name"`
	Displayname pgtype.Varchar     `json:"displayname"`
	Description pgtype.Varchar     `json:"description"`
	Type        pgtype.Varchar     `json:"type"`
	Parent      pgtype.Text        `json:"parent"`
	Link        pgtype.Varchar     `json:"link"`
	Icon        pgtype.Varchar     `json:"icon"`
	Status      pgtype.Varchar     `json:"status"`
	Octime      pgtype.Timestamptz `json:"creattime"`
	Lmtime      pgtype.Timestamptz `json:"lasmodifytime"`
	Open        pgtype.Bool        `json:"open"`
}

type TblCompanyPacks struct {
	Id            pgtype.Text        `json:"id"`
	Companyid     pgtype.Varchar     `json:"companyid"`
	Packid        pgtype.Varchar     `json:"packid"`
	Startdate     pgtype.Timestamptz `json:"startdate"`
	Expirydate    pgtype.Timestamptz `json:"expirydate"`
	Userrolelimit pgtype.Numeric     `json:"userrolelimit"`
	Userlimit     pgtype.Numeric     `json:"userlimit"`
	Branchlimit   pgtype.Numeric     `json:"branchlimit"`
	Status        pgtype.Varchar     `json:"status"`
	Octime        pgtype.Timestamptz `json:"octime"`
	Lmtime        pgtype.Timestamptz `json:"lmtime"`
}

type TblCompany struct {
	CompanyId          pgtype.Varchar     `json:"companyId"`
	CompanyName        pgtype.Text        `json:"companyName"`
	CompanyShortName   pgtype.Varchar     `json:"companyShortName"`
	CompanyCategory    pgtype.Varchar     `json:"companyCategory"`
	CompanyStatus      pgtype.Varchar     `json:"companyStatus"`
	CompanyDescription pgtype.Text        `json:"companyDescription"`
	CompanyImageUrl    pgtype.Text        `json:"companyImageUrl"`
	CompanyLogo        pgtype.Text        `json:"companyLogo"`
	CompanyIndustry    pgtype.Varchar     `json:"companyIndustry"`
	CompanyTaxID       pgtype.Varchar     `json:"companyTaxID"`
	CompanyAddLine1    pgtype.Varchar     `json:"companyAddLine1"`
	CompanyAddLine2    pgtype.Varchar     `json:"companyAddLine2"`
	CompanyCity        pgtype.Varchar     `json:"companyCity"`
	CompanyState       pgtype.Varchar     `json:"companyState"`
	CompanyCountry     pgtype.Varchar     `json:"companyCountry"`
	CompanyPinCode     pgtype.Numeric     `json:"companyPinCode"`
	CompanyPhone       pgtype.Text        `json:"companyPhone"`
	CompanyFax         pgtype.Text        `json:"companyFax"`
	CompanyMobile      pgtype.Text        `json:"companyMobile"`
	CompanyWebsite     pgtype.Text        `json:"companyWebsite"`
	CompanyEmail       pgtype.Text        `json:"companyEmail"`
	CompanyStartDate   pgtype.Date        `json:"companyStartDate"`
	CompanyFiscalYear  pgtype.Numeric     `json:"companyFiscalYear"`
	CompanyTimeZone    pgtype.Text        `json:"companyTimeZone"`
	CompanyBaseCurency pgtype.Varchar     `json:"companyBaseCurency"`
	CompanysParent     pgtype.Text        `json:"companysParent"`
	Isdefault          pgtype.Varchar     `json:"isdefault"`
	Lmuserid           pgtype.Varchar     `json:"lmuserid"`
	Octime             pgtype.Timestamptz `json:"octime"`
	Lmtime             pgtype.Timestamptz `json:"lmtime"`
}

type TblPacks struct {
	Id          pgtype.Varchar     `json:"id"`
	Name        pgtype.Varchar     `json:"name"`
	Displayname pgtype.Varchar     `json:"displayname"`
	Description pgtype.Varchar     `json:"description"`
	Type        pgtype.Varchar     `json:"type"`
	Parent      pgtype.TextArray   `json:"parent"`
	Link        pgtype.Varchar     `json:"link"`
	Icon        pgtype.Varchar     `json:"icon"`
	Status      pgtype.Varchar     `json:"status"`
	Octime      pgtype.Timestamptz `json:"octime"`
	Lmtime      pgtype.Timestamptz `json:"lmtime"`
	Open        pgtype.Bool        `json:"open"`
}
