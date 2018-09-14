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

func (mux *routeHandler) handleAuthenticated(
	pattern string,
	handlerFunc handlers.AuthentictedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrRedirect(handlerFunc, paths.LoginOrSignupPage))
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
	mux.HandleFunc("/favicon.ico", handlers.RedirectToPathHandler("/static/favicon.ico"))

	// pages
	mux.HandleFunc(paths.LoginOrSignupPage, handlers.HandleLoginOrSignupPageRequest)

	mux.handleAuthenticated(paths.HomePage, handlers.HandleHomePageRequest)
	mux.handleAuthenticated(paths.NotesPage, handlers.HandleNotesPageRequest)

	// api

	mux.HandleFunc(paths.UserApi, handlers.HandleUserApiRequest)
	mux.HandleFunc(paths.SessionApi, handlers.HandleSessionApiRequest)

	mux.handleAuthenticated(paths.NoteApi, handlers.HandleNoteApiRequest)

	return mux
}
