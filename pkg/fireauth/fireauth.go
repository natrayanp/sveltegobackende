package fireauth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/sirupsen/logrus"
	commonerrors "github.com/sveltegobackend/pkg/errors"
	"github.com/sveltegobackend/pkg/errors/httperr"

	"google.golang.org/api/option"
)

type FirebaseClient struct {
	AuthClient *auth.Client
}

func Get(clientName string) (*FirebaseClient, error) {
	authClient, err := get(clientName)
	if err != nil {
		return nil, err
	}

	return &FirebaseClient{
		AuthClient: authClient,
	}, nil
}

func get(acjson string) (*auth.Client, error) {
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
		//router.Use(auth.HttpMockMiddleware)
		return nil, fmt.Errorf("MOCK AUTH not allowed")
	}

	//var opts []option.ClientOption
	/*
		if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
			opts = append(opts, option.WithCredentialsFile(file))
		}
	*/
	opt := option.WithCredentialsFile(acjson)
	config := &firebase.Config{ProjectID: "my-project-id"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("Unable to create firebase Auth client")
		return nil, err
	}

	return client, nil

}

func (a FirebaseClient) fireMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bearerToken := a.tokenFromHeader(r)
		if bearerToken == "" {
			httperr.Unauthorised("empty-bearer-token", nil, w, r)
			return
		}

		token, err := a.AuthClient.VerifyIDToken(ctx, bearerToken)
		if err != nil {
			httperr.Unauthorised("unable-to-verify-jwt", err, w, r)
			return
		}

		// it's always a good idea to use custom type as context value (in this case ctxKey)
		// because nobody from the outside of the package will be able to override/read this value
		ctx = context.WithValue(ctx, userContextKey, User{
			UUID:        token.UID,
			Email:       token.Claims["email"].(string),
			Role:        token.Claims["role"].(string),
			DisplayName: token.Claims["name"].(string),
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a FirebaseClient) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}

type User struct {
	UUID  string
	Email string
	Role  string

	DisplayName string
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

var (
	// if we expect that the user of the function may be interested with concrete error,
	// it's a good idea to provide variable with this error
	NoUserInContextError = commonerrors.NewAuthorizationError("no user in context", "no-user-found")
)

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}

	return User{}, NoUserInContextError
}
