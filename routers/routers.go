/*
Package routers defines the path to handler pairs for each endpoint.
*/
package routers

import (
	"net/http"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/paths"
)

// DefineRoutes returns a new servemux with all the required path and handler pairs attached.
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

	// Redirects
	mux.HandleFunc("/", handlers.RedirectToPathHandler(paths.HomePage))
	mux.HandleFunc("/favicon.ico", handlers.RedirectToPathHandler("/static/favicon.ico"))

	// pages
	mux.HandleFunc(paths.LoginOrSignupPage, handlers.HandleLoginOrSignupPageRequest)

	handleAuthenticated(mux, paths.HomePage, handlers.HandleHomePageRequest)
	handleAuthenticated(mux, paths.NotesPage, handlers.HandleNotesPageRequest)

	// ajax
	mux.HandleFunc(paths.User, handlers.HandleUserRequest)
	mux.HandleFunc(paths.Session, handlers.HandleSessionRequest)

	return mux
}

func handleAuthenticated(
	mux *http.ServeMux,
	pattern string,
	handlerFunc handlers.AuthentictedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrRedirectToLogin(handlerFunc))
}
