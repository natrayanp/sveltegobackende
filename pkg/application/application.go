package application

import (
	"github.com/sveltegobackende/pkg/config"
	"github.com/sveltegobackende/pkg/db"
	"github.com/sveltegobackende/pkg/fireauth"
)

// Application holds commonly used app wide data, for ease of DI
type Application struct {
	DB       *db.DB
	Cfg      *config.Config
	FireAuth *fireauth.FirebaseClient
}

// Get captures env vars, establishes DB connection and keeps/returns
// reference to both
func Get() (*Application, error) {
	cfg := config.Get()

	db, err := db.Get(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	client, err := fireauth.Get(cfg.GetFireaccoutn())

	return &Application{
		DB:       db,
		Cfg:      cfg,
		FireAuth: client,
	}, nil
}
