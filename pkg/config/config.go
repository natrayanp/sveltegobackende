package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type QueConfig struct {
	WorkerEnabled bool
	WorkerCount   int64
	QueName       string
}

type Config struct {
	dbUser           string
	dbPswd           string
	dbHost           string
	dbPort           string
	dbName           string
	testDBHost       string
	testDBName       string
	apiPort          string
	migrate          string
	queWorkerEnabled string
	queName          string
	queWorkerCount   string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("POSTGRES_USER"), "DB user name")
	flag.StringVar(&conf.dbPswd, "dbpswd", os.Getenv("POSTGRES_PASSWORD"), "DB pass")
	flag.StringVar(&conf.dbPort, "dbport", os.Getenv("POSTGRES_PORT"), "DB port")
	flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("POSTGRES_HOST"), "DB host")
	flag.StringVar(&conf.dbName, "dbname", os.Getenv("POSTGRES_DB"), "DB name")
	flag.StringVar(&conf.testDBHost, "testdbhost", os.Getenv("TEST_DB_HOST"), "test database host")
	flag.StringVar(&conf.testDBName, "testdbname", os.Getenv("TEST_DB_NAME"), "test database name")
	flag.StringVar(&conf.apiPort, "apiPort", os.Getenv("API_PORT"), "API Port")
	flag.StringVar(&conf.migrate, "migrate", "up", "specify if we should be migrating DB 'up' or 'down'")
	flag.StringVar(&conf.queWorkerEnabled, "workerenabled", os.Getenv("QUEUE_WORKER_ENABLED"), "Worker is part of this application")
	flag.StringVar(&conf.queName, "quename", os.Getenv("QUEUE_WORKER_QUENAME"), "Worker que name for this app")
	flag.StringVar(&conf.queWorkerCount, "queworkercount", os.Getenv("QUEUE_WORKERCOUNT"), "Worker que name for this app")

	flag.Parse()

	return conf
}

func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

func (c *Config) GetTestDBConnStr() string {
	return c.getDBConnStr(c.testDBHost, c.testDBName)
}

func (c *Config) getDBConnStr(dbhost, dbname string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.dbUser,
		c.dbPswd,
		dbhost,
		c.dbPort,
		dbname,
	)
}

func (c *Config) GetAPIPort() string {
	return ":" + c.apiPort
}

func (c *Config) GetMigration() string {
	return c.migrate
}

func (c *Config) GetFireaccoutn() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path + "/pkg/config/firebaseServiceAccount.json"
}

func (c *Config) GetQueconf() *QueConfig {
	i, err := strconv.ParseBool(c.queWorkerEnabled)
	if nil != err {
		return &QueConfig{}
	}

	j, err := strconv.ParseInt(c.queWorkerCount, 10, 0)
	if nil != err {
		return &QueConfig{}
	}

	return &QueConfig{
		WorkerEnabled: i,
		WorkerCount:   j,
		QueName:       c.queName,
	}
}
