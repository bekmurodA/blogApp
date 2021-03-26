package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"snippetbox/handlers"
	"snippetbox/mysql"
	"time"

	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql"
)

//dsn driver specific parameter
var dsn = flag.String("dsn", `web:3029455American.@/snippetbox?parseTime=true`, "MySQL database connection")
var addr = flag.String("addr", ":4000", "HTTP Network Address")
func main() {
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err.Error())
	}
	templateCache, err := handlers.NewTemplateCache("./ui/html/")
	if err != nil {

		errorLog.Fatal(err.Error())
	}
	session := sessions.New([]byte("ajywryw"))
	session.Lifetime = 12 * time.Hour
	defer db.Close()
	app := &handlers.Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Snippets:      &mysql.SnippetModel{DB: db},
		Sessions:      session,
		TemplateCache: templateCache,
		Users:&mysql.UserModel{DB:db},
	}
	//tls.Config struct to hold the non-default TLS settings
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:        *addr,
		ErrorLog:    errorLog,
		Handler:     app.Route(),
		TLSConfig:   tlsConfig,
		//timeouts
		IdleTimeout: time.Minute,
		ReadTimeout: 5*time.Second,
		WriteTimeout:10*time.Second,
	}

	infoLog.Printf("starting server on:%s", *addr)
	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
