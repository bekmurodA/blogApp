package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"snippetbox/forms"
	"snippetbox/models"
	"snippetbox/mysql"
	"strconv"

	"github.com/golangcollege/sessions"
)

type contextKey string
var contextKeyUser = contextKey("user")


type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Snippets      *mysql.SnippetModel
	Sessions      *sessions.Session
	Users         *mysql.UserModel
	TemplateCache map[string]*template.Template
}

func (a *Application) SignupUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MatchesPattern("email", forms.EmailRX)
		form.MinLength("password", 10)
		if !form.Valid() {
			a.render(w, r, "signup.page.tmpl", &TemplateData{Form: form})
			return
		}
		err = a.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err == models.ErrDuplicateEmail {
			form.Errors.Add("email", "Address is already in use")
			a.render(w, r, "signup.page.tmpl", &TemplateData{Form: form})
			return
		} else if err != nil {
			a.serverError(w, err)
			return
		}
		a.Sessions.Put(r, "flash", "Your signup was successful. Please log in.")

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	} else {

		a.SignupUserForm(w, r)
		return
	}
}

func (a *Application) SignupUserForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "signup.page.tmpl", &TemplateData{Form: forms.New(nil)})
}
func (a *Application) LoginUserForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "login.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (a *Application) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		id, err := a.Users.Authenticate(form.Get("email"), form.Get("password"))
		if err == models.ErrInvalidCredentials {
			form.Errors.Add("generic", "Email or password is incorrect")
			a.render(w, r, "login.page.tmpl", &TemplateData{
				Form: form,
			})
			return
		} else if err != nil {
			a.serverError(w, err)
			return
		}
		a.Sessions.Put(r, "userID", id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		a.LoginUserForm(w, r)
		return
	}
}

func (a *Application) LogoutUser(w http.ResponseWriter, r *http.Request) {
	//remove the userid from the session
	a.Sessions.Remove(r,"userID")
	a.Sessions.Put(r,"flash","You have been logged out successfully!")
	http.Redirect(w,r,"/",303)

}
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	s, err := app.Snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &TemplateData{Snippets: s})
}
func (app *Application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	//	values := r.URL.Query()
	id, err := strconv.Atoi(r.URL.Path[len("/snippet/"):])
	//	fmt.Println(r.URL.Path[len("snippet/"):],id)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.Snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "show.page.tmpl", &TemplateData{
		Snippet: s,
	})
}
func (app *Application) CreateSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &TemplateData{Form: forms.New(nil)})
}
func (app *Application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("title", "content", "expires")
		form.MaxLength("title", 100)
		form.PermittedValues("expires", "365", "7", "1")

		//checking if the form is valid
		if !form.Valid() {
			app.render(w, r, "create.page.tmpl", &TemplateData{Form: form})
			return
		}

		id, err := app.Snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.Sessions.Put(r, "flash", "Snippet successfully created")
		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	} else {
		app.CreateSnippetForm(w, r)
		return
	}
}
