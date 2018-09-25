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
}

func (mux *routeHandler) handleAuthenticatedPage(
	pattern string,
	handlerFunc handlers.AuthenticatedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrRedirect(handlerFunc, paths.LoginOrSignupPage))
}

func (mux *routeHandler) handleAuthenticatedApi(
	pattern string,
	handlerFunc handlers.AuthenticatedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrReturnUnauthorized(handlerFunc))
}

// DefineRoutes returns a new servemux with all the required path and handler pairs attached.
func DefineRoutes() http.Handler {
	mux := &routeHandler{http.NewServeMux()}
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
	mux.HandleFunc(paths.LoginOrSignupPage, handlers.HandleLoginOrSignupPageRequest)

	mux.handleAuthenticatedPage(paths.HomePage, handlers.HandleHomePageRequest)
	mux.handleAuthenticatedPage(paths.NotesPage, handlers.HandleNotesPageRequest)

	// api

	mux.HandleFunc(paths.UserApi, handlers.HandleUserApiRequest)
	mux.HandleFunc(paths.SessionApi, handlers.HandleSessionApiRequest)

	mux.handleAuthenticatedApi(paths.NoteApi, handlers.HandleNoteApiRequest)
	mux.handleAuthenticatedApi(paths.CategoryApi, handlers.HandleNoteCateogryApiRequest)

	return mux
}
