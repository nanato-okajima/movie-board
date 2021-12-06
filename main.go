package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// session variable
var sesName = "mvboard-session"
var cs = sessions.NewCookieStore([]byte("secret-key-1234"))

var envpath = "./.env"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/post", post)
	mux.HandleFunc("/group", group)

	//auth route
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}

	server.ListenAndServe()
}
