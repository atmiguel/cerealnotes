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
	env *handlers.Environment,
	pattern string,
	handlerFunc handlers.AuthenticatedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrRedirect(env, handlerFunc, paths.LoginOrSignupPage))
}

func (mux *routeHandler) handleAuthenticatedApi(
	env *handlers.Environment,
	pattern string,
	handlerFunc handlers.AuthenticatedRequestHandlerType,
) {
	mux.HandleFunc(pattern, handlers.AuthenticateOrReturnUnauthorized(env, handlerFunc))
}

func (mux *routeHandler) handleUnAutheticedRequest(
	env *handlers.Environment,
	pattern string,
	handlerFunc handlers.UnauthenticatedEndpointHandlerType,
) {
	mux.HandleFunc(pattern, handlers.WrapUnauthenticatedEndpoint(env, handlerFunc))
}

// DefineRoutes returns a new servemux with all the required path and handler pairs attached.
func DefineRoutes(env *handlers.Environment) http.Handler {
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
	mux.handleUnAutheticedRequest(env, paths.LoginOrSignupPage, handlers.HandleLoginOrSignupPageRequest)

	mux.handleAuthenticatedPage(env, paths.HomePage, handlers.HandleHomePageRequest)
	mux.handleAuthenticatedPage(env, paths.NotesPage, handlers.HandleNotesPageRequest)

	// api

	mux.handleUnAutheticedRequest(env, paths.UserApi, handlers.HandleUserApiRequest)
	mux.handleUnAutheticedRequest(env, paths.SessionApi, handlers.HandleSessionApiRequest)

	mux.handleAuthenticatedApi(env, paths.NoteApi, handlers.HandleNoteApiRequest)
	mux.handleAuthenticatedApi(env, paths.NoteCategoryApi, handlers.HandleNoteCateogryApiRequest)

	return mux
}
