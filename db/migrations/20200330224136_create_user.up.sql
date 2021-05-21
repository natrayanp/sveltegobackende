CREATE SCHEMA ac;

---Login table
CREATE TABLE ac.userlogin (
    userid 		               varchar(100) NOT NULL,
	username    		       varchar(100),  --firebase. User.displayname
    useremail                  varchar(100),
    userpassword               varchar(1000),
    userstatus		           varchar(2) NOT NULL, --> (A- Active, B-Blocked, D-Deleted) 
    emailverified              boolean,
    siteid                     varchar(100) NOT NULL,  --> This tobe depricated
    userstatlstupdt	           timestamptz NOT NULL,    
    octime			           timestamptz NOT NULL,
    lmtime			           timestamptz NOT NULL,
    CONSTRAINT uid PRIMARY KEY (userid, siteid)  
    );

-- Creation of Admin user
--INSERT INTO ac.userlogin (userid,username,useremail,userpassword,userstatus,logintype,usertype,siteid,userstatlstupdt,octime,lmtime)
--VALUES ('fsvV7CG2yDZsBt0ZsNMgCnVZgl02','admin','nat@gmail.com','','A','T','I','ac',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

--INSERT INTO ac.userlogin (userid,username,useremail,userpassword,userstatus,logintype,usertype,userstatlstupdt,octime,lmtime)
--VALUES ('userid1','testuser@gmail.com','testuser@gmail.com','testpas1!','A','S','I','ac',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);


CREATE SEQUENCE companyid_seq START 1;

--- Domain map table which allow user to user their own URL

CREATE TABLE ac.domainmap (
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
    parent                text[],
    link                  varchar(1000),
    icon                  varchar(100),
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);

/*
insert into ac.packs values ('PKS1','POS','POS','POS has all the POS functionalities','pack','null','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS2','POS Function','POS Function','Functions related to POS','module',ARRAY[‘PKS1’],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS3','POS Reports','POS Reports','Reports related to POS','module',ARRAY[‘PKS1’],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS4','POS Settings','POS Settings','Setting for POS module','module',ARRAY['PKS1’,’PKS6'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS5','POS Generic Settings','Generic Settings','Generic settings for POS','function',ARRAY['PKS4'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS6','Settings','Settings','Settings','pack','null','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS7','Entity Settings','Entity Configuration','This module has all the entity level settings','module',ARRAY['PKS6'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS8','companysettigs','Company','This has the functions for company set up','function',ARRAY['PKS7'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS9','branchsettings','Branch','This has the functions for Branch set up','function',ARRAY['PKS7'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS10','User Settings','User Config','This module has all the user level settings','module',ARRAY['PKS6'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS11','user role','Roles','This has the functions for user role set up','function',ARRAY['PKS10'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
insert into ac.packs values ('PKS12','User Settings','Users','This has the functions for user set up','function',ARRAY['PKS10'],'A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);
*/


CREATE SEQUENCE companypacksid_seq START 1;

-- SITE Packages  (This attached to the company)
CREATE TABLE ac.companypacks (
    id                    text NOT NULL CONSTRAINT sitepackid PRIMARY KEY DEFAULT 'CPCKID'||nextval('companypacksid_seq'::regclass)::text,
    companyid             varchar(100) NOT NULL,
    packid                varchar(20) NOT NULL,  --> This can have only PACK type from pack table
    startdate             timestamptz NOT NULL,
    expirydate            timestamptz NOT NULL,
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

-- SITE Plan PACKs  (This is the plan card)
CREATE TABLE ac.planpacks (
    id                    varchar(20) NOT NULL CONSTRAINT planpackid PRIMARY KEY,    
    packid                varchar(20) NOT NULL,  --> This can have only PACK type from pack table
    planid                varchar(20) NOT NULL,
    userrolelimit         numeric(10),
    userlimit             numeric(10),  
    branchlimit           numeric(10),
    durationmonth         numeric(10),
    status                varchar(3),
    octime			      timestamptz NOT NULL,
    lmtime			      timestamptz NOT NULL
);



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
