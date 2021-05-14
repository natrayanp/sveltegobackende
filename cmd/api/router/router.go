package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/sveltegobackend/cmd/api/handlers/createuser"
	"github.com/sveltegobackend/cmd/api/handlers/getuser"
	"github.com/sveltegobackend/pkg/application"

	//"github.com/julienschmidt/httprouter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Get(app *application.Application) *chi.Mux {

	r := chi.NewRouter()

	// Config
	setMiddlewares(r)

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		//r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(app.FireAuth.)

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
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})
	})

	return r
}

func setMiddlewares(router *chi.Mux) {

	//Config
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router)
	//addAuthMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}
