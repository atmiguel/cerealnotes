package routers

import (
	"net/http"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/paths"
)

func DefineRoutes() http.Handler {
	mux := http.NewServeMux()
	// static files
	{
		staticDirectoryName := "static"
		staticDirectoryPaddedWithSlashes := "/" + staticDirectoryName + "/"

		fileServer := http.FileServer(http.Dir(staticDirectoryName))

		mux.Handle(
			staticDirectoryPaddedWithSlashes,
			http.StripPrefix(staticDirectoryPaddedWithSlashes, fileServer))
	}

	mux.HandleFunc("/", handlers.RedirectRequestToHome)

	// templates
	mux.HandleFunc(paths.LoginOrSignupPath, handlers.HandleLoginOrSignupRequest)

	// forms
	mux.HandleFunc("/user", handlers.HandleUserRequest)
	mux.HandleFunc("/session", handlers.HandleSessionRequest)

	// requires Authentication
	handleAuthenticated(mux, paths.HomePath, handlers.HandleHomeRequest)

	return mux
}

func handleAuthenticated(
	mux *http.ServeMux,
	pattern string,
	handlerFunc handlers.AuthentictedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrRedirectToLogin(handlerFunc))
}
