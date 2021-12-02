package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"youtube/config"
	"youtube/my"
)

var DB *gorm.DB

// db variable
func init() {
	conf := config.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})

}

// session variable
var sesName = "mvboard-session"
var cs = sessions.NewCookieStore([]byte("secret-key-1234"))

var envpath = "./.env"

// login check
func checkLogin(w http.ResponseWriter, rq *http.Request) *my.User {
	ses, _ := cs.Get(rq, sesName)
	if ses.Values["login"] == nil || !ses.Values["login"].(bool) {
		http.Redirect(w, rq, "/login", 302)
	}
	ac := ""
	ses.Values["accout"] = "flower"
	if ses.Values["accout"] != nil {
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

// get target Template.
func page(fname string) *template.Template {
	tmps, _ := template.ParseFiles("templates/"+fname+".html", "templates/head.html", "templates/footer.html")
	return tmps
}

// top page handler
func index(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)

	var pl []my.Post
	DB.Where("group_id > 0").Order("created_at desc").Limit(10).Find(&pl)
	var gl []my.Group
	DB.Order("created_at desc").Limit(10).Find(&gl)

	item := struct {
		Title   string
		Message string
		Name    string
		Account string
		Plist   []my.Post
		Glist   []my.Group
	}{
		Title:   "Index",
		Message: "This is Top page.",
		Name:    user.Name,
		Account: user.Account,
		Plist:   pl,
		Glist:   gl,
	}
	err := page("index").Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

// top page handler
func post(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)

	pid := rq.FormValue("pid")

	if rq.Method == "POST" {
		msg := rq.PostFormValue("message")
		pId, _ := strconv.Atoi(pid)
		cmt := my.Comment{
			UserId:  int(user.Model.ID),
			PostId:  pId,
			Message: msg,
		}
		DB.Create(&cmt)
	}

	var pst my.Post
	var cmts []my.CommentJoin

	DB.Where("id = ?", pid).First(&pst)
	DB.Table("comments").Select("comments.*, user.id, users.name").Joins("join users on users.id = comments.user_id").Where("comments.post_id = ?", pid).Order("created_at desc").Find(&cmts)

	item := struct {
		Title   string
		Message string
		Name    string
		Account string
		Post    my.Post
		Clist   []my.CommentJoin
	}{
		Title:   "Post",
		Message: "Post id=" + pid,
		Name:    user.Name,
		Account: user.Account,
		Post:    pst,
		Clist:   cmts,
	}
	err := page("post").Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

// home handler
func home(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)

	if rq.Method == "POST" {
		switch rq.PostFormValue("form") {
		case "post":
			ad := rq.PostFormValue("address")
			ad = strings.TrimSpace(ad)
			if strings.HasPrefix(ad, "http://youtu.be/") {
				ad = strings.TrimPrefix(ad, "https://youtu.be/")
			}

			pt := my.Post{
				UserId:  int(user.Model.ID),
				Address: ad,
				Message: rq.PostFormValue("message"),
			}
			DB.Create(&pt)
		case "group":
			gp := my.Group{
				UserId:  int(user.Model.ID),
				Name:    rq.PostFormValue("name"),
				Message: rq.PostFormValue("message"),
			}
			DB.Create(&gp)
		}
	}

	var pts []my.Post
	var gps []my.Group

	DB.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	DB.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&gps)

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		Plist   []my.Post
		Glist   []my.Group
	}{
		Title:   "Home",
		Message: "User account=\"" + user.Account + "\".",
		Name:    user.Name,
		Account: user.Account,
		Plist:   pts,
		Glist:   gps,
	}
	err := page("home").Execute(w, itm)
	if err != nil {
		log.Fatal(err)
	}
}

// group handler
func group(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)

	gid := rq.FormValue("gid")

	if rq.Method == "POST" {
		ad := rq.PostFormValue("address")
		ad = strings.TrimSpace(ad)
		if strings.HasPrefix(ad, "https://youtu.be/") {
			ad = strings.TrimPrefix(ad, "https://youtu.be/")
		}
		gId, _ := strconv.Atoi(gid)
		pt := my.Post{
			UserId:  int(user.Model.ID),
			Address: rq.PostFormValue("message"),
			GroupId: gId,
		}
		DB.Create(&pt)
	}

	var grp my.Group
	var pts []my.Post

	DB.Where("id=?", gid).First(&grp)
	DB.Order("created_at desc").Model(&grp).Association("posts").Find(&pts)

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		Group   my.Group
		Plist   []my.Post
	}{
		Title:   "Group",
		Message: "Group id=" + gid,
		Name:    user.Name,
		Account: user.Account,
		Group:   grp,
		Plist:   pts,
	}
	err := page("group").Execute(w, itm)
	if err != nil {
		log.Fatal(err)
	}
}

// login handler
func login(w http.ResponseWriter, rq *http.Request) {
	item := struct {
		Title   string
		Message string
		Account string
	}{
		Title:   "Login",
		Message: "type your account & password:",
		Account: "",
	}

	if rq.Method == "GET" {
		err := page("login").Execute(w, item)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if rq.Method == "POST" {
		usr := rq.PostFormValue("account")
		pass := rq.PostFormValue("pass")
		item.Account = usr

		// check account and password
		var re int64
		var user my.User
		DB.Where("account = ? and password = ?", usr, pass).Find(&user).Count(&re)

		if re <= 0 {
			item.Message = "Wrong account or password."
			page("login").Execute(w, item)
			return
		}

		// logined.
		ses, _ := cs.Get(rq, sesName)
		ses.Values["login"] = true
		ses.Values["account"] = usr
		ses.Values["name"] = user.Name
		ses.Save(rq, w)
		http.Redirect(w, rq, "/", 302)
	}

	err := page("login").Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

// logout handler.
func logout(w http.ResponseWriter, rq *http.Request) {
	ses, _ := cs.Get(rq, sesName)
	ses.Values["login"] = nil
	ses.Values["account"] = nil
	ses.Save(rq, w)
	http.Redirect(w, rq, "/login", 302)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		index(w, rq)
	})
	http.HandleFunc("/home", func(w http.ResponseWriter, rq *http.Request) {
		home(w, rq)
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, rq *http.Request) {
		post(w, rq)
	})
	http.HandleFunc("/group", func(w http.ResponseWriter, rq *http.Request) {
		group(w, rq)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, rq *http.Request) {
		login(w, rq)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, rq *http.Request) {
		logout(w, rq)
	})

	http.ListenAndServe(":8000", nil)
}
