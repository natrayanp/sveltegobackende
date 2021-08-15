package models

import "github.com/jackc/pgtype"

type TblUserlogin struct {
	Userid          pgtype.Varchar      `json:"userid"`
	Username        pgtype.Varchar      `json:"username"`
	Useremail       pgtype.Varchar      `json:"useremail"`
	Userpassword    pgtype.Varchar      `json:"userpassword"`
	Userstatus      pgtype.Varchar      `json:"userstatus"`
	Emailverified   pgtype.Bool         `json:"emailverified"`
	Siteid          pgtype.Varchar      `json:"siteid"`
	Hostname        pgtype.Text         `json:"hostname"`
	Selecthostname  pgtype.Text         `json:"selecthostname"`
	Companyid       pgtype.Varchar      `json:"companyid"`
	Companyowner    pgtype.Varchar      `json:"companyowner"`
	Entityid        pgtype.VarcharArray `json:"entityid"`
	Userstatlstupdt pgtype.Timestamptz  `json:"userstatuslastupdate"`
	Octime          pgtype.Timestamptz  `json:"creattime"`
	Lmtime          pgtype.Timestamptz  `json:"lasmodifytime"`
}

type TblMytree struct {
	Id          pgtype.Varchar      `json:"id"`
	Name        pgtype.Varchar      `json:"name"`
	Displayname pgtype.Varchar      `json:"displayname"`
	Description pgtype.Varchar      `json:"description"`
	Type        pgtype.Varchar      `json:"type"`
	Parent      pgtype.VarcharArray `json:"parent"`
	Sortorder   pgtype.Int2         `json:"sortorder"`
	Link        pgtype.Varchar      `json:"link"`
	Icon        pgtype.Varchar      `json:"icon"`
	Status      pgtype.Varchar      `json:"status"`
	Octime      pgtype.Timestamptz  `json:"creattime"`
	Lmtime      pgtype.Timestamptz  `json:"lasmodifytime"`
	Open        pgtype.Bool         `json:"open"`
	Submenu     []*TblMytree        `json:"submenu"`
}

type TblCompanyPacks struct {
	Id            pgtype.Text        `json:"id"`
	Companyid     pgtype.Varchar     `json:"companyid"`
	Planid        pgtype.Varchar     `json:"planid"`
	Packid        pgtype.Varchar     `json:"packid"`
	Startdate     pgtype.Date        `json:"startdate"`
	Expirydate    pgtype.Date        `json:"expirydate"`
	Userrolelimit pgtype.Numeric     `json:"userrolelimit"`
	Userlimit     pgtype.Numeric     `json:"userlimit"`
	Branchlimit   pgtype.Numeric     `json:"branchlimit"`
	Status        pgtype.Varchar     `json:"status"`
	Octime        pgtype.Timestamptz `json:"octime"`
	Lmtime        pgtype.Timestamptz `json:"lmtime"`
}

type TblCompany struct {
	Companyid          pgtype.Varchar     `json:"companyId"`
	Companyname        pgtype.Text        `json:"companyName"`
	Companyshortname   pgtype.Varchar     `json:"companyShortName"`
	Companycategory    pgtype.Varchar     `json:"companyCategory"`
	Companystatus      pgtype.Varchar     `json:"companyStatus"`
	Companydescription pgtype.Text        `json:"companyDescription"`
	Companyimageurl    pgtype.Text        `json:"companyImageUrl"`
	Companylogo        pgtype.Text        `json:"companyLogo"`
	Companyindustry    pgtype.Varchar     `json:"companyIndustry"`
	Companytaxid       pgtype.Varchar     `json:"companyTaxID"`
	Companyaddline1    pgtype.Varchar     `json:"companyAddLine1"`
	Companyaddline2    pgtype.Varchar     `json:"companyAddLine2"`
	Companycity        pgtype.Varchar     `json:"companyCity"`
	Companystate       pgtype.Varchar     `json:"companyState"`
	Companycountry     pgtype.Varchar     `json:"companyCountry"`
	Companypincode     pgtype.Varchar     `json:"companyPinCode"`
	Companyphone       pgtype.Text        `json:"companyPhone"`
	Companyfax         pgtype.Text        `json:"companyFax"`
	Companymobile      pgtype.Text        `json:"companyMobile"`
	Companywebsite     pgtype.Text        `json:"companyWebsite"`
	Companyemail       pgtype.Text        `json:"companyEmail"`
	Companystartdate   pgtype.Date        `json:"companyStartDate"`
	Companyfiscalyear  pgtype.Varchar     `json:"companyFiscalYear"`
	Companytimezone    pgtype.Text        `json:"companyTimeZone"`
	Companybasecurency pgtype.Varchar     `json:"companyBaseCurency"`
	Companysparent     pgtype.Text        `json:"companysParent"`
	Isdefault          pgtype.Varchar     `json:"isdefault"`
	Lmuserid           pgtype.Varchar     `json:"lmuserid"`
	Octime             pgtype.Timestamptz `json:"octime"`
	Lmtime             pgtype.Timestamptz `json:"lmtime"`
}

type TblRefdata struct {
	Id        pgtype.Text      `json:"id"`
	Refvalcat pgtype.Varchar   `json:"refvalcat"`
	Refvalue  pgtype.Varchar   `json:"refvalue"`
	Parent    pgtype.TextArray `json:"parent"`
	Submenu   []*TblRefdata    `json:"submenu"`
}

type TblBranch struct {
	Companyid         pgtype.Varchar     `json:"companyId"`
	Companyname       pgtype.Text        `json:"companyName"`
	Branchid          pgtype.Varchar     `json:"branchId"`
	Branchname        pgtype.Text        `json:"branchName"`
	Branchshortname   pgtype.Varchar     `json:"branchShortName"`
	Branchcategory    pgtype.Varchar     `json:"branchCategory"`
	Branchstatus      pgtype.Varchar     `json:"branchStatus"`
	Branchdescription pgtype.Text        `json:"branchDescriptio"`
	Branchimageurl    pgtype.Text        `json:"branchImageUrl"`
	Branchaddline1    pgtype.Varchar     `json:"branchAddLine1"`
	Branchaddline2    pgtype.Varchar     `json:"branchAddLine2"`
	Branchcountry     pgtype.Varchar     `json:"branchCountry"`
	Branchstate       pgtype.Varchar     `json:"branchState"`
	Branchcity        pgtype.Varchar     `json:"branchCity"`
	Branchpincode     pgtype.Varchar     `json:"branchPinCode"`
	Branchphone       pgtype.Text        `json:"branchPhone"`
	Branchfax         pgtype.Text        `json:"branchFax"`
	Branchmobile      pgtype.Text        `json:"branchMobile"`
	Branchwebsite     pgtype.Text        `json:"branchWebsite"`
	Branchemail       pgtype.Text        `json:"branchEmail"`
	Branchstartdate   pgtype.Date        `json:"branchStartDate"`
	Isdefault         pgtype.Varchar     `json:"isdefault"`
	Lmuserid          pgtype.Varchar     `json:"lmuserid"`
	Octime            pgtype.Timestamptz `json:"octime"`
	Lmtime            pgtype.Timestamptz `json:"lmtime"`
}
