package main

import (
	"github.com/joho/godotenv"
	"github.com/sveltegobackend/cmd/api/router"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/exithandler"
	"github.com/sveltegobackend/pkg/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		//	logs.Info.Println("failed to load env vars")
	}

	app, err := application.Get()
	if err != nil {
		//logs.Error.Fatal(err.Error())
	}

	srv := server.
		Get().
		WithAddr(app.Cfg.GetAPIPort()).
		WithRouter(router.Get(app))
	//.		WithErrLogger(logs.Error)

	go func() {
		//logs.Info.Printf("starting server at %s", app.Cfg.GetAPIPort())
		if err := srv.Start(); err != nil {
			//logs.Error.Fatal(err.Error())
		}
	}()

	exithandler.Init(func() {
		if err := srv.Close(); err != nil {
			//logs.Error.Println(err.Error())
		}

		if err := app.DB.Close(); err != nil {
			//logs.Error.Println(err.Error())
		}
	})
}
