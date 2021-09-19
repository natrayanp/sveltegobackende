package models

import (
	"time"
)

/*
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
*/

type TblUserlogin struct {
	Userid          string    `json:"userid"`
	Username        *string   `json:"username"`
	Useremail       *string   `json:"useremail"`
	Userpassword    *string   `json:"userpassword"`
	Userstatus      string    `json:"userstatus"`
	Emailverified   *bool     `json:"emailverified"`
	Siteid          string    `json:"siteid"`
	Hostname        string    `json:"hostname"`
	Selecthostname  *string   `json:"selecthostname"`
	Companyid       *string   `json:"companyid"`
	Companyowner    string    `json:"companyowner"`
	Entityid        []*string `json:"entityid"`
	Userstatlstupdt time.Time `json:"userstatuslastupdate"`
	Octime          time.Time `json:"creattime"`
	Lmtime          time.Time `json:"lasmodifytime"`
}

/*
type TblMytree struct {
	Id          pgtype.Varchar      `json:"id"`
	Name        pgtype.Varchar      `json:"name"`
	Displayname pgtype.Varchar      `json:"displayname"`
	Description pgtype.Varchar      `json:"description"`
	Type        pgtype.Varchar      `json:"type"`
	Menulevel   pgtype.Varchar      `json:"menulevel"`
	Allowedops  pgtype.BoolArray    `json:"allowedops"`
	Parent      pgtype.VarcharArray `json:"parent"`
	Sortorder   pgtype.Int2         `json:"sortorder"`
	Link        pgtype.Varchar      `json:"link"`
	Icon        pgtype.Varchar      `json:"icon"`
	Status      pgtype.Varchar      `json:"status"`
	Octime      pgtype.Timestamptz
	Lmtime      pgtype.Timestamptz
	Open        pgtype.Bool  `json:"open"`
	Submenu     []*TblMytree `json:"submenu"`
}
*/
type TtblMytree struct {
	Companyid     string        `json:"companyid"`
	Branchid      string        `json:"branchid"`
	Rolemasterid  string        `json:"rolemasterid"`
	Packid        string        `json:"packid"`
	Packfuncid    string        `json:"packfuncid"`
	Status        string        `json:"Status"`
	Name          string        `json:"name"`
	Displayname   string        `json:"displayname"`
	Description   string        `json:"description"`
	Type          string        `json:"type"`
	Parent        []*string     `json:"parent"`
	Link          *string       `json:"link"`
	Icon          *string       `json:"icon"`
	Startdate     time.Time     `json:"startdate"`
	Expirydate    time.Time     `json:"expirydate"`
	Userrolelimit *int          `json:"userrolelimit"`
	Userlimit     *int          `json:"userlimit"`
	Branchlimit   *int          `json:"branchlimit"`
	Compstatus    string        `json:"Compstatus"`
	Sortorder     int           `json:"sortorder"`
	Menulevel     string        `json:"menulevel"`
	Allowedops    []bool        `json:"allowedops"`
	Allowedopsval []bool        `json:"allowedopsval"`
	Userid        *string       `json:"userid"`
	Disablefunc   bool          `json:"disablefunc"`
	Basketname    string        `json:"basketname"`
	Open          bool          `json:"open"`
	Submenu       []*TtblMytree `json:"submenu"`
	Roledetailid  *string       `json:"roledetailid"`
	Octime        time.Time     `json:"creattime"`
	Lmtime        time.Time     `json:"lasmodifytime"`
}

/*
type TblCompanyPacks struct {
	Id            pgtype.Text        `json:"id"`
	Companyid     pgtype.Varchar     `json:"companyid"`
	Planid        pgtype.Varchar     `json:"planid"`
	Packfuncid    pgtype.Varchar     `json:"packfuncid"`
	Startdate     pgtype.Date        `json:"startdate"`
	Expirydate    pgtype.Date        `json:"expirydate"`
	Userrolelimit pgtype.Numeric     `json:"userrolelimit"`
	Userlimit     pgtype.Numeric     `json:"userlimit"`
	Branchlimit   pgtype.Numeric     `json:"branchlimit"`
	Status        pgtype.Varchar     `json:"status"`
	Octime        pgtype.Timestamptz `json:"octime"`
	Lmtime        pgtype.Timestamptz `json:"lmtime"`
}
*/

type TblCompanyPacks struct {
	Cpypacksid    string    `json:"cpypacksid"`
	Companyid     string    `json:"companyid"`
	Planid        string    `json:"planid"`
	Packfuncid    string    `json:"packfuncid"`
	Startdate     time.Time `json:"startdate"`
	Expirydate    time.Time `json:"expirydate"`
	Userrolelimit *int      `json:"userrolelimit"`
	Userlimit     *int      `json:"userlimit"`
	Branchlimit   *int      `json:"branchlimit"`
	Status        string    `json:"status"`
	Octime        time.Time `json:"octime"`
	Lmtime        time.Time `json:"lmtime"`
}

/*
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
*/
type TblCompany struct {
	Companyid          string    `json:"companyId"`
	Companyname        string    `json:"companyName"`
	Companyshortname   string    `json:"companyShortName"`
	Companycategory    string    `json:"companyCategory"`
	Companystatus      string    `json:"companyStatus"`
	Companydescription string    `json:"companyDescription"`
	Companyimageurl    *string   `json:"companyImageUrl"`
	Companylogo        *string   `json:"companyLogo"`
	Companyindustry    string    `json:"companyIndustry"`
	Companytaxid       string    `json:"companyTaxID"`
	Companyaddline1    string    `json:"companyAddLine1"`
	Companyaddline2    string    `json:"companyAddLine2"`
	Companycity        string    `json:"companyCity"`
	Companystate       string    `json:"companyState"`
	Companycountry     string    `json:"companyCountry"`
	Companypincode     string    `json:"companyPinCode"`
	Companyphone       *string   `json:"companyPhone"`
	Companyfax         *string   `json:"companyFax"`
	Companymobile      *string   `json:"companyMobile"`
	Companywebsite     *string   `json:"companyWebsite"`
	Companyemail       *string   `json:"companyEmail"`
	Companystartdate   time.Time `json:"companyStartDate"`
	Companyfiscalyear  string    `json:"companyFiscalYear"`
	Companytimezone    *string   `json:"companyTimeZone"`
	Companybasecurency string    `json:"companyBaseCurency"`
	Companysparent     *string   `json:"companysParent"`
	Isdefault          string    `json:"isdefault"`
	Lmuserid           string    `json:"lmuserid"`
	Octime             time.Time `json:"octime"`
	Lmtime             time.Time `json:"lmtime"`
}

/*
type TblRefdata struct {
	Id        pgtype.Text      `json:"id"`
	Refvalcat pgtype.Varchar   `json:"refvalcat"`
	Refvalue  pgtype.Varchar   `json:"refvalue"`
	Sortorder pgtype.Int2      `json:"sortorder"`
	Parent    pgtype.TextArray `json:"parent"`
	Submenu   []*TblRefdata    `json:"submenu"`
}
*/

type TblRefdata struct {
	Refid     string        `json:"Refid"`
	Refvalcat string        `json:"refvalcat"`
	Refvalue  string        `json:"refvalue"`
	Sortorder int           `json:"sortorder"`
	Parent    []*string     `json:"parent"`
	Submenu   []*TblRefdata `json:"submenu"`
}

/*
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
*/
type TblBranch struct {
	Companyid         string    `json:"companyId"`
	Companyname       string    `json:"companyName"`
	Branchid          string    `json:"branchId"`
	Branchname        string    `json:"branchName"`
	Branchshortname   string    `json:"branchShortName"`
	Branchcategory    string    `json:"branchCategory"`
	Branchstatus      string    `json:"branchStatus"`
	Branchdescription string    `json:"branchDescriptio"`
	Branchimageurl    *string   `json:"branchImageUrl"`
	Branchaddline1    string    `json:"branchAddLine1"`
	Branchaddline2    string    `json:"branchAddLine2"`
	Branchcountry     string    `json:"branchCountry"`
	Branchstate       string    `json:"branchState"`
	Branchcity        string    `json:"branchCity"`
	Branchpincode     string    `json:"branchPinCode"`
	Branchphone       *string   `json:"branchPhone"`
	Branchfax         *string   `json:"branchFax"`
	Branchmobile      *string   `json:"branchMobile"`
	Branchwebsite     *string   `json:"branchWebsite"`
	Branchemail       *string   `json:"branchEmail"`
	Branchstartdate   time.Time `json:"branchStartDate"`
	Isdefault         string    `json:"isdefault"`
	Lmuserid          string    `json:"lmuserid"`
	Octime            time.Time `json:"octime"`
	Lmtime            time.Time `json:"lmtime"`
}
