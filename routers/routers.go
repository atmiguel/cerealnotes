/*
Package routers defines the path to handler pairs for each endpoint.
*/
package routers

import (
	"net/http"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/paths"
)

type routeHandler struct {
	*http.ServeMux
	Env *handlers.Environment
}

func (mux *routeHandler) handleAuthenticatedPage(
	pattern string,
	handlerFunc handlers.AuthenticatedRequestHandlerType,
) {
	mux.HandleFunc(pattern, mux.Env.AuthenticateOrRedirect(handlerFunc, paths.LoginOrSignupPage))
}

func (mux *routeHandler) handleAuthenticatedApi(
	pattern string,
	handlerFunc handlers.AuthenticatedRequestHandlerType,
) {
	mux.HandleFunc(pattern, mux.Env.AuthenticateOrReturnUnauthorized(handlerFunc))
}

// DefineRoutes returns a new servemux with all the required path and handler pairs attached.
func DefineRoutes(env *handlers.Environment) http.Handler {
	mux := &routeHandler{http.NewServeMux(), env}
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
	mux.HandleFunc("/api/", http.NotFound)
	mux.HandleFunc("/favicon.ico", handlers.RedirectToPathHandler("/static/favicon.ico"))

	// pages
	mux.HandleFunc(paths.LoginOrSignupPage, env.HandleLoginOrSignupPageRequest)

	mux.handleAuthenticatedPage(paths.HomePage, env.HandleHomePageRequest)
	mux.handleAuthenticatedPage(paths.NotesPage, env.HandleNotesPageRequest)

	// api

	mux.HandleFunc(paths.UserApi, env.HandleUserApiRequest)
	mux.HandleFunc(paths.SessionApi, env.HandleSessionApiRequest)

	mux.handleAuthenticatedApi(paths.NoteApi, env.HandleNoteApiRequest)
	mux.handleAuthenticatedApi(paths.CategoryApi, env.HandleCategoryApiRequest)

	return mux
}
