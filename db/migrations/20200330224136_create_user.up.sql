CREATE TABLE IF NOT EXISTS gue_jobs
(
    job_id      bigserial   NOT NULL PRIMARY KEY,
    priority    smallint    NOT NULL,
    run_at      timestamptz NOT NULL,
    job_type    text        NOT NULL,
    args        json        NOT NULL,
    error_count integer     NOT NULL DEFAULT 0,
    last_error  text,
    queue       text        NOT NULL,
    created_at  timestamptz NOT NULL,
    updated_at  timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS "idx_gue_jobs_selector" ON "gue_jobs" ("queue", "run_at", "priority");

COMMENT ON TABLE gue_jobs IS '1';






CREATE SCHEMA ac;

---Login table
CREATE TABLE ac.userlogin (
    userid 		               varchar(100) NOT NULL,
	username    		       varchar(100),  --firebase. User.displayname
    useremail                  varchar(100),
    userpassword               varchar(1000),
    userstatus		           varchar(2) NOT NULL, --> (A- Active, B-Blocked, D-Deleted) 
    emailverified              boolean,
    siteid                     varchar(100) NOT NULL,
    hostname    text NOT NULL,          
    selecthostname    text UNIQUE,          
    companyid   varchar(100),
    userstatlstupdt	           timestamptz NOT NULL,    
    octime			           timestamptz NOT NULL,
    lmtime			           timestamptz NOT NULL,
    CONSTRAINT uid PRIMARY KEY (userid, hostname,siteid)  
    );

-- Creation of Admin user
--INSERT INTO ac.userlogin (userid,username,useremail,userpassword,userstatus,logintype,usertype,siteid,userstatlstupdt,octime,lmtime)
--VALUES ('fsvV7CG2yDZsBt0ZsNMgCnVZgl02','admin','nat@gmail.com','','A','T','I','ac',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

--INSERT INTO ac.userlogin (userid,username,useremail,userpassword,userstatus,logintype,usertype,userstatlstupdt,octime,lmtime)
--VALUES ('userid1','testuser@gmail.com','testuser@gmail.com','testpas1!','A','S','I','ac',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);


CREATE SEQUENCE companyid_seq START 1;

--- Domain map table which allow user to user their own URL

CREATE TABLE ac.domainmap (
    domainmapid varchar(100) NOT NULL DEFAULT 'DMAPID'||nextval('companyid_seq'::regclass)::varchar(100),
    hostname    text NOT NULL UNIQUE,    
    siteid      text NOT NULL,
    companyid   varchar(100) NOT NULL DEFAULT 'CPYID'||nextval('companyid_seq'::regclass)::varchar(100),
    status      varchar(3) NOT NULL,
    octime      timestamptz NOT NULL,
    lmtime      timestamptz NOT NULL
);



---Secret key detail table

CREATE TABLE ac.secrettkn (
	secretcode 		            varchar(100) NOT NULL CONSTRAINT secretcode PRIMARY KEY, 
    seccdid 		            varchar(30) NOT NULL,  --> DDMMYYYYHHMMSS
    secoctime			        timestamp NOT NULL
    );
INSERT INTO ac.secrettkn VALUES ('secret01','31082019193003',CURRENT_TIMESTAMP);



-- Login history
CREATE TABLE ac.loginh (
    userid 		               varchar(100) NOT NULL,
    ipaddress                  varchar(25),
    sessionid                  varchar(100),
    companyid                  varchar(100),
    logintime                  timestamptz NOT NULL, 
    logoutime                  timestamptz 
);


-- Packages
CREATE TABLE ac.packs (
    id                    varchar(20) NOT NULL CONSTRAINT packid PRIMARY KEY,
    name                  varchar(100) NOT NULL,
    displayname           varchar(50) NOT NULL,
    description           varchar NOT NULL,
    type                  varchar(30) NOT NULL,
    parent                varchar(20)[],
    sortorder             smallserial NOT NULL,
    link                  varchar(1000),
    icon                  varchar(100),
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);

/*
insert into ac.packs values ('PKS1','POS','POS','POS has all the POS functionalities','pack',ARRAY[NULL],'','radio_button_checked','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS2','POS Function','POS Function','Functions related to POS','module',ARRAY['PKS1'],'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS3','POS Reports','POS Reports','Reports related to POS','module',ARRAY['PKS1'],'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS4','POS Settings','POS Settings','Setting for POS module','module',ARRAY['PKS1’,’PKS6'],'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS5','POS Generic Settings','Generic Settings','Generic settings for POS','function',ARRAY['PKS4'],'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS6','Settings','Settings','Settings','pack',ARRAY[NULL],'/landing/settings','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS7','Entity Settings','Entity Configuration','This module has all the entity level settings','module',ARRAY['PKS6'],'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS8','companysettigs','Company','This has the functions for company set up','function',ARRAY['PKS7'],'/landing/settings/companysettings','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS9','branchsettings','Branch','This has the functions for Branch set up','function',ARRAY['PKS7'],'/landing/settings/branchsettings','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS10','User Settings','User Config','This module has all the user level settings','module',ARRAY['PKS6'],'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS11','user role','Roles','This has the functions for user role set up','function',ARRAY['PKS10'],'/landing/settings/roles','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS12','User Settings','Users','This has the functions for user set up','function',ARRAY['PKS10'],'/landing/settings/users','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

insert into ac.packs values ('PKS13','Pricing','Pricing','Pricing plans avaialble','pack',ARRAY[NULL],'/landing/pricing','fa-hand-holding-heart','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS14','Pricing','Pricing','Pricing plans avaialble','module',ARRAY['PKS13'],'','','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS15','Pricing','Pricing','Pricing plans avaialble','function',ARRAY['PKS14'],'','','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/




CREATE SEQUENCE companypacksid_seq START 1;

-- SITE Packages  (This attached to the company)
CREATE TABLE ac.companypacks (
    id                    text NOT NULL CONSTRAINT sitepackid PRIMARY KEY DEFAULT 'CPCKID'||nextval('companypacksid_seq'::regclass)::text,
    companyid             varchar(100) NOT NULL,
    planid                    varchar(20) NOT NULL ,
    packid                varchar(20) NOT NULL,  --> This can have only PACK type from pack table
    startdate             date NOT NULL,
    expirydate            date NOT NULL,
    userrolelimit         numeric(10),
    userlimit             numeric(10),  
    branchlimit           numeric(10),
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);



-- SITE Plan (This is the plan card)
CREATE TABLE ac.plan (
    id                    varchar(20) NOT NULL CONSTRAINT planid PRIMARY KEY,
    name                  varchar(100) NOT NULL,
    displayname           varchar(50) NOT NULL,
    description           varchar NOT NULL,
    currency              varchar(3),
    amount                numeric(15,2),
    discountedamt         numeric(15,2),
    startdate             timestamptz NOT NULL,
    enddate               timestamptz NOT NULL,
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);

--insert into ac.plan values ('PLANID1','Free','Free','Free plan available for all','SGD',0,0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);


-- SITE Plan PACKs  (This is the plan card)
CREATE TABLE ac.planpacks (
    id                    varchar(20) NOT NULL CONSTRAINT planpackid PRIMARY KEY,    
    packid                varchar(20) NOT NULL,  --> This can have only PACK type from pack table
    planid                varchar(20) NOT NULL,
    userrolelimit         numeric(10),
    userlimit             numeric(10),  
    branchlimit           numeric(10),
    durationdays         integer,
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);


/*
insert into ac.planpacks values ('PACKID1','PKS8','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.planpacks values ('PACKID2','PKS9','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.planpacks values ('PACKID3','PKS14','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/


--RoleMaster

CREATE TABLE ac.rolemaster (    
    id                    varchar(20) NOT NULL CONSTRAINT rolemasterid PRIMARY KEY, 
    name                  varchar(100) NOT NULL,
    displayname           varchar(50) NOT NULL,
    description           varchar NOT NULL,
    companyid             varchar(30) NOT NULL,
    branchid              varchar(30) NOT NULL,     
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);

/*
insert into ac.ROLEMASTER values ('ROLMA1','SignupAdmin','SignupAdmin','This is the role given to users when they sign up','PUBLIC','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.ROLEMASTER values ('ROLMA2','defaultadmin','defaultadmin','This is the role given to users when they completed creation of their first Company and branch','PUBLIC','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/


--defaultRoleMaster

    CREATE TABLE ac.defaultrolemaster (    
        id                    varchar(20) NOT NULL CONSTRAINT defaultrolemasterid PRIMARY KEY, 
        name                  varchar(100) NOT NULL,
        displayname           varchar(50) NOT NULL,
        description           varchar NOT NULL,
        companyid             varchar(30) NOT NULL,
        branchid              varchar(30) NOT NULL,        
        status                varchar(3) NOT NULL,
        octime			      timestamptz NOT NULL,
        lmtime			      timestamptz NOT NULL
    );

/*
insert into ac.defaultrolemaster values ('DROLMA1','SignupAdmin','SignupAdmin','This is the role given to users when they sign up','PUBLIC','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultrolemaster values ('DROLMA2','defaultadmin','defaultadmin','This is the role given to users when they completed creation of their first Company and branch','PUBLIC','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/

--RoleDetails

CREATE TABLE ac.roledetails (    
    id                    varchar(20) NOT NULL CONSTRAINT roledetailsid PRIMARY KEY, 
    rolemasterid          varchar(100) NOT NULL,    
    packfuncid            varchar(20) NOT NULL,  --> This can have only function type from pack table
    planid               varchar(30) NOT NULL,
    companyid             varchar(30) NOT NULL,
    branchid              varchar(30) NOT NULL,                    
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);

/*
insert into ac.Roledetails values ('ROLDET1','ROLMA1','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET2','ROLMA1','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET3','ROLMA1','PKS14','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

insert into ac.Roledetails values ('ROLDET1','ROLMA1','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET2','ROLMA1','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET3','ROLMA2','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET4','ROLMA2','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET5','ROLMA2','PKS11','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET6','ROLMA2','PKS12','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/

--defaultRoleDetails

CREATE TABLE ac.defaultroledetails (    
    id                    varchar(20) NOT NULL CONSTRAINT droledetailsid PRIMARY KEY, 
    rolemasterid          varchar(100) NOT NULL,    
    packid                varchar(20) NOT NULL,  --> This can have only function type from pack table
    companyid             varchar(30) NOT NULL,
    branchid              varchar(30) NOT NULL,       
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);
/*
insert into ac.defaultroledetails values ('DROLDET1','ROLMA1','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET2','ROLMA1','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET3','ROLMA2','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET4','ROLMA2','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET5','ROLMA2','PKS11','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET6','ROLMA2','PKS12','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/

--userrole

CREATE TABLE ac.userrole ( 
    userid 		          varchar(100) NOT NULL,
    rolemasterid          varchar(20) NOT NULL,       
    companyid             varchar(30) NOT NULL,
    branchid              varchar(30) NOT NULL,      
    status                varchar(3) NOT NULL,  --> D -Delete / A - Active
    isdefault             varchar(3) NOT NULL,  --> Y/N
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);


CREATE SEQUENCE refid_seq START 1;


CREATE TABLE ac.refdata (
    id                     text NOT NULL CONSTRAINT refid PRIMARY KEY DEFAULT 'REFID'||nextval('refid_seq'::regclass)::text,
    refcode               varchar(100) NOT NULL,
    refvalcat              varchar(100) NOT NULL,
    refvalue               varchar(100) NOT NULL,
    description           varchar NOT NULL,
    parent                varchar(20)[],
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);
/*
INSERT INTO ac.refdata values (DEFAULT,'industype','industype','Hotel','industry type of the company',ARRAY[NULL],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'compcat','compcat','Food','Company type category',ARRAY[NULL],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'compcat','compcat','FMCG','Company type FMCG',ARRAY[NULL],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);*/
---
CREATE TABLE ac.company (
        companyid                      varchar(100) NOT NULL UNIQUE,        
        companyName             text NOT NULL,
        companyShortName        varchar(100) NOT NULL,
        companyCategory         varchar(100) NOT NULL,
        companyStatus           varchar(3) NOT NULL,  --> D -Delete / A - Active
        companyDescription      text NOT NULL,
        companyImageUrl         text,
        companyLogo             text,
        companyIndustry         varchar(100) NOT NULL,
        companyTaxID            varchar(100) NOT NULL,
        companyAddLine1         varchar(100) NOT NULL,
        companyAddLine2         varchar(100) NOT NULL,
        companyCity             varchar(100) NOT NULL,
        companyState            varchar(100) NOT NULL,
        companyCountry          varchar(100) NOT NULL,
        companyPinCode          numeric	 NOT NULL,
        companyPhone            text,
        companyFax              text,
        companyMobile           text,
        companyWebsite          text,
        companyEmail            text,
        companyStartDate        date NOT NULL,
        companyFiscalYear       numeric	 NOT NULL,
        companyTimeZone         text,
        companyBaseCurency      varchar(3) NOT NULL,
        companysParent          text,   
        isdefault               varchar(3) NOT NULL,  --> Y/N
        lmuserid                varchar(100) NOT NULL,                  
        octime			        timestamptz NOT NULL,
        lmtime			        timestamptz NOT NULL
);























---Entity details
create table entity
(
    entityid          varchar(100)  default (('EN'::text || to_char(CURRENT_TIMESTAMP, 'DDMMYYYY'::text)) ||
                                            nextval('unihot.entityid_sequence'::regclass)) not null
                                             constraint enity_pkey primary key,
    entityname        varchar(100)  not null,
    entityshortname   varchar(100),
    entitycategory    varchar(100),
    entitystatus      varchar(1),    
    entityimageurl    varchar(100),
    entitylogo        varchar(100),
    entityindustry    varchar(100),
    entitytaxid       varchar(100),
    entityaddline1    varchar(100),
    entityaddline2    varchar(100),
    entitycity        varchar(100),
    entitystate       varchar(100),
    entitycountry     varchar(100),
    entitypincode     varchar(100),
    entityphone       varchar(100),
    entityfax         varchar(100),
    entitymobile      varchar(100),
    entitywebsite     varchar(100),
    entityemail       varchar(100),
    entitystartdate   varchar(100),
    entityfiscalyear  varchar(100),
    entitytimezone    varchar(100),
    octime            timestamp with time zone                                             not null,
    lmtime            timestamp with time zone                                             not null
);



--INSERT for public entity

---Entity details
create table entitybranch
(
    entitybranchid          varchar(100) default (('BR'::text || to_char(CURRENT_TIMESTAMP, 'DDMMYYYY'::text)) ||
                                                  nextval('unihot.entitybranchid_sequence'::regclass)) not null
        constraint enitybranch_pkey
            primary key,
    entityid                varchar(100)                                                               not null
        constraint entitybranch_entityid_fkey
            references entity,
    entitybranchname        varchar(100),
    entitybranchshortname   varchar(100),
    entitybranchcategory    varchar(100),
    entitybranchstatus      varchar(100),
    entitybranchdescription varchar(100),
    entitybranchimageurl    varchar(100),
    entitybranchaddline1    varchar(100),
    entitybranchaddline2    varchar(100),
    entitybranchcity        varchar(100),
    entitybranchstate       varchar(100),
    entitybranchcountry     varchar(100),
    entitybranchpincode     varchar(100),
    entitybranchphone       varchar(100),
    entitybranchfax         varchar(100),
    entitybranchmobile      varchar(100),
    entitybranchwebsite     varchar(100),
    entitybranchemail       varchar(100),
    entitybranchstartdate   varchar(100),
    octime                  timestamp with time zone                                                   not null,
    lmtime                  timestamp with time zone                                                   not null
);
--INSERT for public entity

---user access permission table
CREATE TABLE unihot.useraccess (
    userid 		               varchar(100) NOT NULL,
    logintype                  varchar(2) NOT NULL,  --> based on admin user (T - Thirdparty, S - Standalone)
    usertype                   varchar(2) NOT NULL, --> based on admin user (C - Thirdparty COMPANY, I - Thirdparty individual, S - Standalone Company)
    entity                     varchar(20) NOT NULL ,
    entitybranch               varchar(10) NOT NULL ,
    defaultindicator           varchar(1) NOT NULL,
    roleid                     varchar(100),  --> from role setup table ADMIN,READONLY,WRITE,NODELETE
    site                       varchar(100),  --> nc - Nawalcube, dv - developer, au - auth
    accessstatus	           varchar(2) NOT NULL, --> (A- Active, B-Blocked) for the site
    octime			           timestamptz NOT NULL,
    lmtime			           timestamptz NOT NULL,
    CONSTRAINT usac PRIMARY KEY (userid, logintype, usertype, entity, entitybranch, site)  
    );

--INSERT for public entity
INSERT INTO unihot.useraccess VALUES ('01','PUBLIC','01','nc','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

---role id table
CREATE TABLE unihot.roledetails (
    roleid 		               varchar(100) NOT NULL,
    rolename                   varchar(100) NOT NULL,
    entity                     varchar(20) NOT NULL REFERENCES unihot.enity(entityid),
    site                       varchar(100),  --> nc - Nawalcube, dv - developer, au - auth
    roleidstatus	           varchar(2) NOT NULL, --> (A- Active, B-Blocked) for the site
    octime			           timestamptz NOT NULL,
    lmtime			           timestamptz NOT NULL,
    CONSTRAINT us PRIMARY KEY (roleid, entity, site)  
    );

--INSERT for public entity
INSERT INTO unihot.roledetails VALUES ('01','PUBLIC','01','nc','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);