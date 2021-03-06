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
    userid 		                varchar(100) NOT NULL,
	username    		        varchar(100),  --firebase. User.displayname
    useremail                   varchar(100),
    userpassword                varchar(1000),
    userstatus		            varchar(2) NOT NULL, --> (A- Active, B-Blocked, D-Deleted) 
    emailverified               boolean,
    siteid                      varchar(100) NOT NULL,
    hostname                    text NOT NULL,          
    selecthostname              text UNIQUE,          
    companyid                   varchar(100),
    companyowner                varchar(3) NOT NULL, --> Y, N, DK
    entityid                    varchar(100)[],  --> This holds the ID enity id will have companyids tagged in another table which is TODO
    userstatlstupdt	            timestamptz NOT NULL,    
    octime			            timestamptz NOT NULL,
    lmtime			            timestamptz NOT NULL,
    CONSTRAINT uid PRIMARY KEY (userid, hostname,siteid)  
    );

-- Creation of Admin user
--INSERT INTO ac.userlogin (userid,username,useremail,userpassword,userstatus,logintype,usertype,siteid,userstatlstupdt,octime,lmtime)
--VALUES ('fsvV7CG2yDZsBt0ZsNMgCnVZgl02','admin','nat@gmail.com','','A','T','I','ac',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

--INSERT INTO ac.userlogin (userid,username,useremail,userpassword,userstatus,logintype,usertype,userstatlstupdt,octime,lmtime)
--VALUES ('userid1','testuser@gmail.com','testuser@gmail.com','testpas1!','A','S','I','ac',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);


CREATE SEQUENCE ac.companyid_seq START 1;

--- Domain map table which allow user to use their own URL

CREATE TABLE ac.domainmap (
    domainmapid varchar(100) NOT NULL DEFAULT 'DMAPID'||nextval('ac.companyid_seq'::regclass)::varchar(100),
    hostname    text NOT NULL UNIQUE,    
    siteid      text NOT NULL,
    companyid   varchar(100) NOT NULL DEFAULT 'CPYID'||nextval('ac.companyid_seq'::regclass)::varchar(100),
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



CREATE TABLE ac.userprofile (
    usrprof_userid 		        varchar(100) NOT NULL,  -- User id from ac.userlogin
    usrprof_firstname           varchar(100) NOT NULL,
    usrprof_lastname            varchar(100) NOT NULL,
    usrprof_designation         varchar(100) NOT NULL,
    usrprof_department          varchar(100) NOT NULL,
    usrprof_gender              varchar(100) NOT NULL,
    usrprof_dob                 date NOT NULL,
    usrprof_AddLine1            varchar(100) NOT NULL,
    usrprof_AddLine2            varchar(100) NOT NULL,
    usrprof_city                varchar(100) NOT NULL,
    usrprof_state               varchar(100) NOT NULL,
    usrprof_country             varchar(100) NOT NULL,
    usrprof_pinCode             varchar(50)	 NOT NULL,
    usrprof_mobile              text,
    usrprof_email               text,
    usrprof_joiningdate         date NOT NULL,
    usrprof_lastdate            date NOT NULL,
    usrprof_taxid               date NOT NULL,
    usrprof_companyid           varchar(100),
    usrprof_edurefid            varchar(50),  -- future implementation
    usrprof_exprefid            varchar(50),  -- future implementation
    usrprof_imagelink           text,
    usrprof_octime			    timestamptz NOT NULL,
    usrprof_lmtime			    timestamptz NOT NULL
);





-- Packages
CREATE TABLE ac.packs (
    packid                varchar(20) NOT NULL CONSTRAINT packid PRIMARY KEY,
    packgroupid           varchar(20)[] NOT NULL,  
    name                  varchar(100) NOT NULL,
    displayname           varchar(50) NOT NULL,
    description           varchar NOT NULL,
    type                  varchar(30) NOT NULL,
    menulevel             varchar(30) NOT NULL,
    allowedops            boolean[] NOT NULL,  -- True-show checkbox; False-hide checkbox. link-roledetails.allowedopsval;ac.refdata.refcode = 'allowedops'
                                                -- This should hold value for all values in refcode index should be equal to (sort order-1).
                                                -- This value should rollup meaning...all function may not have export but the parent should have true for an ops if atleast one of it child has that ops
    parent                varchar(20)[],
    sortorder             smallserial NOT NULL,
    link                  varchar(1000),
    icon                  varchar(100),
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);


insert into ac.packs values ('PKS1',ARRAY['PKSGP1'],'POS','POS','POS has all the POS functionalities','pack','COMPANY', ARRAY[true,true,true,true,false],ARRAY[NULL],1,'','radio_button_checked','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS2',ARRAY['PKSGP1'],'POS Function','POS Function','Functions related to POS','module','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS1'],1,'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS3',ARRAY['PKSGP1'],'POS Reports','POS Reports','Reports related to POS','module','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS1'],2,'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS4',ARRAY['PKSGP1','PKSGP2'],'POS Settings','POS Settings','Setting for POS module','module','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS1???,???PKS6'],3,'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS5',ARRAY['PKSGP1'],'POS Generic Settings','Generic Settings','Generic settings for POS','function','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS4'],1,'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS6',ARRAY['PKSGP2'],'Settings','Settings','Settings','pack','COMPANY', ARRAY[true,true,true,true,false],ARRAY[NULL],2,'/landing/settings','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS7',ARRAY['PKSGP2'],'Entity Settings','Entity Configuration','This module has all the entity level settings','module','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS6'],1,'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS8',ARRAY['PKSGP2'],'companysettigs','Company','This has the functions for company set up','function','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS7'],1,'/landing/settings/companysettings','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS9',ARRAY['PKSGP2'],'branchsettings','Branch','This has the functions for Branch set up','function','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS7'],2,'/landing/settings/branchsettings','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS10',ARRAY['PKSGP2'],'User Settings','User Config','This module has all the user level settings','module','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS6'],3,'','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS11',ARRAY['PKSGP2'],'user role','Roles','This has the functions for user role set up','function','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS10'],1,'/landing/settings/roles','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS12',ARRAY['PKSGP2'],'User Settings','Users','This has the functions for user set up','function','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS10'],1,'/landing/settings/users','fa-cog','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS13',ARRAY['PKSGP3'],'Pricing','Pricing','Pricing plans avaialble','pack','COMPANY', ARRAY[true,true,true,true,false],ARRAY[NULL],4,'/landing/pricing','fa-hand-holding-heart','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS14',ARRAY['PKSGP3'],'Pricing','Pricing','Pricing plans avaialble','module','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS13'],1,'','','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS15',ARRAY['PKSGP3'],'Pricing','Pricing','Pricing plans avaialble','function','COMPANY', ARRAY[true,true,true,true,false],ARRAY['PKS14'],1,'','','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);





CREATE SEQUENCE ac.companypacksid_seq START 1;

-- SITE Packages  (This attached to the company)
CREATE TABLE ac.companypacks (
    cpypacksid                    text NOT NULL CONSTRAINT sitepackid PRIMARY KEY DEFAULT 'CPCKID'||nextval('ac.companypacksid_seq'::regclass)::text,
    companyid             varchar(100) NOT NULL,
    planid                    varchar(20) NOT NULL ,
    packfuncid             varchar(20) NOT NULL,  --> This can have only PACK type from pack table
    startdate             date NOT NULL,
    expirydate            date NOT NULL,
    userrolelimit         numeric(10),
    userlimit             numeric(10),  
    branchlimit           numeric(10),
    status                varchar(3) NOT NULL,
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);



-- SITE Plan (This is the plan card)
CREATE TABLE ac.plan (
    planid                    varchar(20) NOT NULL CONSTRAINT planid PRIMARY KEY,
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

insert into ac.plan values ('PLANID1','Free','Free','Free plan available for all','SGD',0,0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);


-- SITE Plan PACKs  (This is the plan card)
CREATE TABLE ac.planpacks (
    planpackid                    varchar(20) NOT NULL CONSTRAINT planpackid PRIMARY KEY,    
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


/* This is to be created for each plan
    currently iam adding a free plan
*/
insert into ac.planpacks values ('PACKID1','PKS8','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.planpacks values ('PACKID2','PKS9','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.planpacks values ('PACKID3','PKS14','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.planpacks values ('PACKID4','PKS11','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.planpacks values ('PACKID5','PKS12','PLANID1',10,10,10,90,'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);



--RoleMaster

CREATE SEQUENCE ac.rolemasterid_seq START 1;

CREATE TABLE ac.rolemaster (    
    rolemasterid         varchar(100) DEFAULT (('ROLMA'::text || to_char(CURRENT_TIMESTAMP, 'DDMMYYYY'::text)) ||
                                                   nextval('ac.rolemasterid_seq'::regclass)) NOT NULL
                                                   constraint rolemaster_pkey primary key,
    rmname                  varchar(100) NOT NULL,
    rmdisplayname           varchar(50) NOT NULL,
    rmdescription           varchar NOT NULL,
    rmcompanyid             varchar(30) NOT NULL,
    rmbranchid              varchar(30) NOT NULL,     
    rmstatus                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);


insert into ac.ROLEMASTER values ('ROLMA1','SignupAdmin','SignupAdmin','This is the role given to users when they sign up','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
/*
insert into ac.ROLEMASTER values ('ROLMA2','defaultadmin','defaultadmin','This is the role given to users when they completed creation of their first Company and branch','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/

/*
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


insert into ac.defaultrolemaster values ('ROLMA1','SignupAdmin','SignupAdmin','This is the role given to users when they sign up','PUBLIC','PUBLIC','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

*/
--RoleDetails
/* PACKFUNCID will have the leaf node value*/

CREATE SEQUENCE ac.roledetailid_seq START 1;

CREATE TABLE ac.roledetails (    
    roledetailid    varchar(100) DEFAULT (('ROLDET'::text || to_char(CURRENT_TIMESTAMP, 'DDMMYYYY'::text)) ||
                                                                    nextval('ac.roledetailid_seq'::regclass)) NOT NULL
                                                                    constraint roledetail_pkey primary key,
    rdrolemasterid          varchar(100) NOT NULL  REFERENCES ac.rolemaster (rolemasterid),    
    rdpackfuncid            varchar(20) NOT NULL,  --> This can have only function type from pack table
    rdplanid               varchar(30),
    rdcompanyid             varchar(30) NOT NULL,
    rdbranchid              varchar(30) NOT NULL,
    rdallowedopsval         boolean[] NOT NULL,    -- True/False - represent checkbox Value. link-packs.allowedops;ac.refdata.refcode = 'allowedops'                                                                 
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);

/*  All the Functions and Modules that are to be part of sign up admin should have should be included here   
    Eventhough we have entry in this table for sign up role.  User will get the modules/functions that are common between this and company packs
*/
insert into ac.Roledetails values ('ROLDET1','ROLMA1','PKS8','PUBLIC','PUBLIC','PUBLIC', ARRAY[true,true,true,true,false],CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET2','ROLMA1','PKS9','PUBLIC','PUBLIC','PUBLIC', ARRAY[true,true,true,true,false],CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET3','ROLMA1','PKS14','PUBLIC','PUBLIC','PUBLIC',ARRAY[true,true,true,true,false],CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET4','ROLMA1','PKS11','PUBLIC','PUBLIC','PUBLIC',ARRAY[true,true,true,true,false],CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.Roledetails values ('ROLDET5','ROLMA1','PKS12','PUBLIC','PUBLIC','PUBLIC',ARRAY[true,true,true,true,false],CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

/*
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

insert into ac.defaultroledetails values ('DROLDET1','ROLMA1','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET2','ROLMA1','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET3','ROLMA2','PKS8','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET4','ROLMA2','PKS9','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET5','ROLMA2','PKS11','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.defaultroledetails values ('DROLDET6','ROLMA2','PKS12','PUBLIC','PUBLIC','PUBLIC',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/

--userrole
/*
CREATE TABLE ac.userrole ( 
    userid 		          varchar(100) NOT NULL,
    rolemasterid          varchar(20) NOT NULL,       
    companyid             varchar(30) NOT NULL,
    userbranchacess    varchar(30)[] NOT NULL,      
    status                varchar(3) NOT NULL,  --> D -Delete / A - Active
    isdefault             varchar(3) NOT NULL,  --> Y/N
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);
*/

CREATE TABLE ac.userrole ( 
    usrrole_userid 		          varchar(100) NOT NULL,
    usrrole_rolemasterid          varchar(20)[] NOT NULL,       
    usrrole_companyid             varchar(30) NOT NULL,
    usrrole_branchidaccess       varchar(30) NOT NULL,      
    usrrole_status                varchar(3) NOT NULL,  --> D -Delete / A - Active
    usrrole_isdefault             varchar(3) NOT NULL,  --> Y/N
    usrrole_octime			      timestamptz NOT NULL,
    usrrole_lmtime			      timestamptz NOT NULL
);



CREATE SEQUENCE ac.refid_seq START 1;



CREATE TABLE ac.refdata (
    refid                     text NOT NULL CONSTRAINT refid PRIMARY KEY DEFAULT 'REFID'||nextval('ac.refid_seq'::regclass)::text,
    refcode               varchar(100) NOT NULL,
    refvalcat              varchar(100) NOT NULL,
    refvalue               varchar(100) NOT NULL,
    description           varchar NOT NULL,
    parent                varchar(20)[],
    status                varchar(3),
    sortorder             smallserial NOT NULL,
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);


INSERT INTO ac.refdata values (DEFAULT,'allowedops','allowedops','READ','allowed operations master used in packs.allowedops and roledetails.allowedopsval',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'allowedops','allowedops','ADD','allowed operations master used in packs.allowedops and roledetails.allowedopsval',ARRAY[NULL],'A',2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'allowedops','allowedops','EDIT','allowed operations master used in packs.allowedops and roledetails.allowedopsval',ARRAY[NULL],'A',3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'allowedops','allowedops','DELETE','allowed operations master used in packs.allowedops and roledetails.allowedopsval',ARRAY[NULL],'A',4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'allowedops','allowedops','EXPORT','allowed operations master used in packs.allowedops and roledetails.allowedopsval',ARRAY[NULL],'A',5,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'country','country','India','Country',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'state','state','Tamilnadu','Country',ARRAY['REFID6'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'state','state','Karnataka','Country',ARRAY['REFID6'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'state','state','Kerala','Country',ARRAY['REFID6'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'state','state','Andhrapradesh','Country',ARRAY['REFID6'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Chennai','cities in a state',ARRAY['REFID7'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Coimbatore','cities in a state',ARRAY['REFID7'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Salem','cities in a state',ARRAY['REFID7'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Namakkal','cities in a state',ARRAY['REFID7'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Bangalore','cities in a state',ARRAY['REFID8'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Mangalore','cities in a state',ARRAY['REFID8'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Cochin','cities in a state',ARRAY['REFID9'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Tiruvandapuram','cities in a state',ARRAY['REFID9'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Guruvayur','cities in a state',ARRAY['REFID9'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Hyrderabad','cities in a state',ARRAY['REFID10'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Secundrabad','cities in a state',ARRAY['REFID10'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'city','city','Vijayawada','cities in a state',ARRAY['REFID10'],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'industype','industype','Hotel','industry type of the company',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'compcat','compcat','Food','Company type category',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'compcat','compcat','FMCG','Company type FMCG',ARRAY[NULL],'A',2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'timezone','timezone','IST','Country timezone',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'timezone','timezone','SGT','Country timezone',ARRAY[NULL],'A',2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'timezone','timezone','GMT','Country timezone',ARRAY[NULL],'A',3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'currency','currency','INR','iso currency code',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'currency','currency','SGD','iso currency code',ARRAY[NULL],'A',2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'currency','currency','USD','iso currency code',ARRAY[NULL],'A',3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'currency','currency','EUR','iso currency code',ARRAY[NULL],'A',4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'finyear','Financial year','JAN-DEC','financial year',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'finyear','Financial year','APR-MAR','financial year',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'gender','gender','Male','Gender value',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'gender','gender','Female','Gender value',ARRAY[NULL],'A',2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'gender','gender','Other','Gender value',ARRAY[NULL],'A',3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'dept','department','Marketing','Department value',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'dept','department','HR','Department value',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'dept','department','Delivery','Department value',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'dept','department','Store','Department value',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'dept','department','Frontdesk','Department value',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

INSERT INTO ac.refdata values (DEFAULT,'designa','designation','Junior officer','Company designation values',ARRAY[NULL],'A',1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'designa','designation','Senior officer','Company designation values',ARRAY[NULL],'A',2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'designa','designation','Assitant Vice President - Band 1','Company designation values',ARRAY[NULL],'A',3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'designa','designation','Assitant Vice President - Band 2','Company designation values',ARRAY[NULL],'A',4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'designa','designation','Senior Vice President - Band 5','Company designation values',ARRAY[NULL],'A',5,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
INSERT INTO ac.refdata values (DEFAULT,'designa','designation','Senior Vice President - Band 4','Company designation values',ARRAY[NULL],'A',6,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

---
CREATE TABLE ac.company (
        companyid               varchar(100) NOT NULL UNIQUE,        
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
        companyPinCode          varchar(50)	 NOT NULL,
        companyPhone            text,
        companyFax              text,
        companyMobile           text,
        companyWebsite          text,
        companyEmail            text,
        companyStartDate        date NOT NULL,
        companyFiscalYear       varchar(10)	 NOT NULL,
        companyTimeZone         text,
        companyBaseCurency      varchar(3) NOT NULL,
        companysParent          text,   
        isdefault               varchar(3) NOT NULL,  --> Y/N
        lmuserid                varchar(100) NOT NULL,                  
        octime			        timestamptz NOT NULL,
        lmtime			        timestamptz NOT NULL
);



CREATE SEQUENCE ac.branchid_sequence START 1;



CREATE TABLE ac.branch (
            companyid           varchar(100) NOT NULL REFERENCES ac.company(companyid),
            branchId            varchar(100)  default (('BR'::text || to_char(CURRENT_TIMESTAMP, 'DDMMYYYY'::text)) ||
                                        nextval('ac.branchid_sequence'::regclass)) NOT NULL
                                        constraint branch_pkey primary key,
            branchName          text NOT NULL,
            branchShortName     varchar(100) NOT NULL,
            branchStatus        varchar(3) NOT NULL,  --> D -Delete / A - Active
            branchImageUrl      text,
            branchAddLine1      varchar(100) NOT NULL,
            branchAddLine2      varchar(100) NOT NULL,
            branchCity          varchar(100) NOT NULL,
            branchState         varchar(100) NOT NULL,
            branchCountry       varchar(100) NOT NULL,
            branchPinCode       varchar(50) NOT NULL,
            branchPhone         text,
            branchFax           text,
            branchMobile        text,
            branchWebsite       text,
            branchEmail         text,
            branchStartDate     date NOT NULL,
            isdefault           varchar(3) NOT NULL,  --> Y/N
            lmuserid            varchar(100) NOT NULL,                  
            octime			    timestamptz NOT NULL,
            lmtime			    timestamptz NOT NULL
);



CREATE SEQUENCE ac.auditid_sequence START 1;

CREATE TABLE ac.audit (
    audit_id        varchar(100)  default (('AD'::text || to_char(CURRENT_TIMESTAMP, 'DDMMYYYY'::text)) ||
                                                   nextval('ac.auditid_sequence'::regclass)) NOT NULL
                                                   constraint audit_pkey primary key,
    itemid         varchar(100) NOT NULL,-- Role/packs
    itemkeys        JSONB NOT NULL,
    action          varchar(3) NOT NULL,-- A-ADD/E-EDIT/D-DELETE
    oldvalue        JSONB[] NOT NULL,-- old value of the fields
    newvalue        JSONB[] NOT NULL,-- New values of the fields
    companyid       varchar(100),
    actionuser      varchar(100) NOT NULL,-- modified user
    octime          timestamptz NOT NULL -- server time of the update
);




/* CREATE VIEWS */
CREATE VIEW AC.COMPANYPACKS_PACKS_VIEW AS(
WITH recursive COMPANYPACKSVIEW AS(
	SELECT B.COMPANYID,b.packfuncid,b.startdate,b.expirydate,b.userrolelimit,B.USERLIMIT,B.BRANCHLIMIT,B.STATUS AS COMPSTATUS,A.* FROM AC.packs A
	JOIN AC.COMPANYPACKS B ON B.PACKFUNCID = A.PACKID 	
	UNION
	SELECT CASE WHEN B.COMPANYID IS NULL THEN t.companyid END::VARCHAR(100),
		   CASE WHEN B.packfuncid IS NULL THEN A.PACKID END::VARCHAR(20),	
		   CASE WHEN B.startdate IS NULL THEN t.startdate END,
		 	CASE WHEN B.expirydate IS NULL THEN t.expirydate END,
			CASE WHEN B.userrolelimit IS NULL THEN t.userrolelimit END::numeric(10,0),			
			CASE WHEN B.USERLIMIT IS NULL THEN t.USERLIMIT END::numeric(10,0),			
			CASE WHEN B.BRANCHLIMIT IS NULL THEN t.BRANCHLIMIT END::numeric(10,0),
			CASE WHEN B.STATUS IS NULL THEN t.STATUS END::varchar(3) AS COMPSTATUS ,
			A.* FROM AC.packs A
	LEFT JOIN AC.COMPANYPACKS B ON B.PACKFUNCID = A.PACKID
	JOIN COMPANYPACKSVIEW AS t ON A.packid = ANY(t.parent)	
) SELECT * FROM COMPANYPACKSVIEW
	);


/*
CREATE VIEW COMPANYPACKS_PACKS_VIEW AS(
WITH recursive COMPANYPACKSVIEW AS(
	SELECT B.COMPANYID,b.packfuncid,b.startdate,b.expirydate,b.userrolelimit,B.USERLIMIT,B.BRANCHLIMIT,B.STATUS AS COMPSTATUS,A.* FROM AC.packs A
	JOIN AC.COMPANYPACKS B ON B.PACKFUNCID = A.PACKID 	
	UNION
	SELECT B.COMPANYID,b.packfuncid,b.startdate,b.expirydate,b.userrolelimit,B.USERLIMIT,B.BRANCHLIMIT,B.STATUS AS COMPSTATUS,A.* FROM AC.packs A
	LEFT JOIN AC.COMPANYPACKS B ON B.PACKFUNCID = A.PACKID
	JOIN COMPANYPACKSVIEW AS t ON A.packid = ANY(t.parent)	
) SELECT * FROM COMPANYPACKSVIEW
	);
*/

CREATE VIEW ac.ROLE_USER_VIEW AS(
 SELECT a.rmcompanyid,
    a.rmbranchid,
    a.rolemasterid,
    a.rmname,
    a.rmdisplayname,
    a.rmdescription,
    a.rmstatus,
    b.roledetailid,
    b.rdpackfuncid,
    b.rdallowedopsval,
    c.usrrole_userid as userid,
    c.usrrole_branchidaccess as userbranchacess
   FROM ac.rolemaster a
     LEFT JOIN ac.roledetails b ON a.rolemasterid = b.rdrolemasterid
     LEFT JOIN ac.userrole c ON a.rolemasterid = ANY(c.usrrole_rolemasterid)
);


/*To get back the value as JSON format
select json_agg(art)
from (
	select *,
    	(select json_agg(b)
		from (
			select * 
            from pfstklist where pfportfolioid = a.pfportfolioid ) as b
         ) as pfstklist,
        	(select json_agg(c)
		from (
			select * 
            from pfmflist where pfportfolioid = a.pfportfolioid ) as c
         ) as pfmflist
	from pfmaindetail as a) art
*/