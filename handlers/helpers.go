package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"snippetbox/models"
	"time"
)

func (a *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (a *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (a *Application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

//if there is a *models.User struct in the request with contextKeyUser then we know
//that the request is coming from an authenticated and valid user
func (a *Application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func (a *Application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = a.Sessions.PopString(r, "flash")
	td.AuthenticatedUser = a.authenticatedUser(r)
	return td

}

func (a *Application) render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := a.TemplateCache[name]
	if !ok {
		a.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, a.addDefaultData(td, r))
	if err != nil {

		a.serverError(w, err)
		return
	}
	buf.WriteTo(w)

}
