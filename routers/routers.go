package routers

import (
	"github.com/atmiguel/cerealnotes/handlers"
	"net/http"
)

func SetRoutes() {
	// static files
	{
		staticDirectoryName := "static"
		staticDirectoryPaddedWithSlashes := "/" + staticDirectoryName + "/"

		fileServer := http.FileServer(http.Dir(staticDirectoryName))

		http.Handle(
			staticDirectoryPaddedWithSlashes,
			http.StripPrefix(staticDirectoryPaddedWithSlashes, fileServer))
	}

	// templates
	http.HandleFunc("/login-or-signup", handlers.HandleLoginOrSignupRequest)

	// forms
	http.HandleFunc("/user", handlers.HandleUserRequest)
	http.HandleFunc("/session", handlers.HandleSessionRequest)

	// requires Authentication
	http.HandleFunc("/", handlers.AuthenticateOrRedirectToLogin(handlers.HandleRootRequest))

}
