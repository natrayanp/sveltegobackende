package application

import (
	"fmt"

	"github.com/sveltegobackend/pkg/config"
	"github.com/sveltegobackend/pkg/db"
	"github.com/sveltegobackend/pkg/fireauth"
)

// Application holds commonly used app wide data, for ease of DI
type Application struct {
	DB             *db.DB
	Cfg            *config.Config
	FireAuthclient *fireauth.FirebaseClient
}

// Get captures env vars, establishes DB connection and keeps/returns
// reference to both
func Get() (*Application, error) {
	cfg := config.Get()

	fmt.Println(cfg.GetDBConnStr())

	db, err := db.Get(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	client, err := fireauth.Get(cfg.GetFireaccoutn())

	return &Application{
		DB:             db,
		Cfg:            cfg,
		FireAuthclient: client,
	}, nil
}
