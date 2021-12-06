package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"youtube/my"
)

// login check
func checkLogin(w http.ResponseWriter, r *http.Request) *my.User {
	ses, _ := cs.Get(r, sesName)
	if ses.Values["login"] == nil || !ses.Values["login"].(bool) {
		http.Redirect(w, r, "/login", 302)
	}
	ac := ""

	if ses.Values["account"] != nil {
		ac = ses.Values["account"].(string)
	}

	var user my.User

	DB.Where("account = ?", ac).First(&user)

	return &user
}

// Template for no-template
func notemp() *template.Template {
	tmp, _ := template.New("index").Parse("NO PAGE.")
	return tmp
}

// login handler
func login(w http.ResponseWriter, r *http.Request) {
	item := struct {
		Title   string
		Message string
		Account string
	}{
		Title:   "Login",
		Message: "type your account & password:",
		Account: "",
	}

	if r.Method == "GET" {
		err := page("login").Execute(w, item)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if r.Method == "POST" {
		usr := r.PostFormValue("account")
		pass := r.PostFormValue("pass")
		item.Account = usr

		// check account and password
		var re int64
		var user my.User
		DB.Where("account = ? and password = ?", usr, pass).Find(&user).Count(&re)
		fmt.Println("na", usr, pass)
		if re <= 0 {
			item.Message = "Wrong account or password."
			page("login").Execute(w, item)
			return
		}

		// logined.
		ses, _ := cs.Get(r, sesName)
		ses.Values["login"] = true
		ses.Values["account"] = usr
		ses.Values["name"] = user.Name
		ses.Save(r, w)
		http.Redirect(w, r, "/", 302)
	}

	err := page("login").Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

// logout handler.
func logout(w http.ResponseWriter, r *http.Request) {
	ses, _ := cs.Get(r, sesName)
	ses.Values["login"] = nil
	ses.Values["account"] = nil
	ses.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}
