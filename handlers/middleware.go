package handlers

import (
	"context"
	"fmt"
	"net/http"
	"snippetbox/models"
)

//If program panics this middleware shuts it gracefully
func (a *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil { //value returned by recover is interface{}
				w.Header().Set("Connection", "close")
				a.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//Prints out user's rmeoveAddr,proto,method and user agent on every request
func (a *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

//restrict the unAuthorized user access these urls
func (a Application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

//Checks if the userID value exists in the session
func (a *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := a.Sessions.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		//fetch the details of the current user from the database.
		//if no matching record is found,m remove the (invalid) userID
		//from their session and call the next handler in the chain
		user, err := a.Users.Get(a.Sessions.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			a.Sessions.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			a.serverError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

//Sets secure headers
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}
