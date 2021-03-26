package handlers

import (
	"net/http"
)

func (a *Application) Route() http.Handler {
	mux := http.NewServeMux()
	//order do matter
	mux.Handle("/", a.Sessions.Enable(a.authenticate(http.HandlerFunc(a.Home))))
	//mux.Handle("/snippet/create", http.HandlerFunc(a.CreateSnippetForm))
	mux.Handle("/snippet/create", a.Sessions.Enable(a.authenticate(a.requireAuthenticatedUser(http.HandlerFunc(a.CreateSnippet)))))
	mux.Handle("/snippet/", a.Sessions.Enable(a.authenticate(http.HandlerFunc(a.ShowSnippet))))
	//user handlers
	mux.Handle("/user/signup", a.Sessions.Enable(a.authenticate(http.HandlerFunc(a.SignupUser))))
	mux.Handle("/user/login", a.Sessions.Enable(a.authenticate(http.HandlerFunc(a.LoginUser))))
	mux.Handle("/user/logout", a.Sessions.Enable(a.authenticate(a.requireAuthenticatedUser(http.HandlerFunc(a.LogoutUser)))))
	return a.recoverPanic(a.logRequest(secureHeaders(mux)))

}
