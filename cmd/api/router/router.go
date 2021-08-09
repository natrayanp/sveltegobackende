package router

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/sveltegobackend/cmd/api/handlers/auth/login"
	"github.com/sveltegobackend/cmd/api/handlers/entity/branch"
	"github.com/sveltegobackend/cmd/api/handlers/entity/company"
	"github.com/sveltegobackend/cmd/api/handlers/refdata"

	"github.com/sveltegobackend/cmd/api/handlers/auth/signup"
	"github.com/sveltegobackend/cmd/api/handlers/createuser"
	"github.com/sveltegobackend/cmd/api/handlers/getuser"
	"github.com/sveltegobackend/pkg/application"
	"github.com/sveltegobackend/pkg/fireauth"
	logs "github.com/sveltegobackend/pkg/logger"

	"github.com/sveltegobackend/pkg/mymiddleware"

	//"github.com/julienschmidt/httprouter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Get(app *application.Application) *chi.Mux {

	r := chi.NewRouter()

	// Config
	setMiddlewares(app, r)

	/*
		// Protected routes
		r.Group(func(r chi.Router) {
			// Seek, verify and validate JWT tokens
			//r.Use(jwtauth.Verifier(tokenAuth))

			r.Use(fireauth.FirebaseClient{AuthClient: app.FireAuthclient.AuthClient}.FireMiddleware)

			// Handle valid / invalid tokens. In this example, we use
			// the provided authenticator middleware, but you can write your
			// own very easily, look at the Authenticator method in jwtauth.go
			// and tweak it, its not scary.
			//r.Use(jwtauth.Authenticator)
			r.Get("/users/:id", getuser.Do(app))
			r.Post("/users", createuser.Do(app))
			/*
				r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
					_, claims, _ := jwtauth.FromContext(r.Context())
					w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
				})
	*/
	/*})
	 */

	r.Mount("/auth", authorisedRouter(app))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {

			w.Write([]byte("welcome anonymous"))
		})
	})

	return r
}

func setMiddlewares(app *application.Application, router *chi.Mux) {

	//Config
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router)
	//addAuthMiddleware(router)
	//router.Use(mymiddleware.ParseHeadMiddleware(app))

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	fmt.Println(len(allowedOrigins))
	if len(allowedOrigins) == 0 {
		return
	}
	corsMiddleware := cors.New(cors.Options{
		//AllowedOrigins:   allowedOrigins,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}

// A completely separate router for administrator routes
func authorisedRouter(app *application.Application) chi.Router {
	r := chi.NewRouter()
	//setMiddlewares(r)
	//r.Use(AdminOnly)
	setMiddlewares(app, r)
	r.Use(fireauth.FirebaseClient{AuthClient: app.FireAuthclient.AuthClient}.FireMiddleware)
	r.Use(mymiddleware.ParseHeadMiddleware(app))
	r.Get("/test", signup.Do(app))
	r.Post("/signuptoken", signup.Do(app))
	r.Post("/logintoken", login.Do(app))
	r.Post("/regisdomain", login.DoRegisDomain(app))
	r.Get("/users/:id", getuser.Do(app))
	r.Post("/users", createuser.Do(app))
	r.Post("/regisplan", login.DoPacks(app))
	r.Get("/getcompany", company.DoFetch(app))
	r.Post("/savecompany", company.DoSave(app))
	r.Post("/getrefdata", refdata.Do(app))
	r.Get("/getbranch", branch.DoBrFetch(app))
	r.Post("/savebranch", branch.DoBrSave(app))

	/*

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("admin: index"))
		})
		r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("admin: list accounts.."))
		})
		r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
		})

	*/

	return r
}
