package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"youtube/my"
)

// db variable
var host = os.Getenv("DB_HOST")
var user = os.Getenv("POSTGRES_USER")
var password = os.Getenv("POSTGRES_PASSWORD")
var dbname = os.Getenv("POSTGRES_DB")
var port = os.Getenv("POSTGRES_PORT")
var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", host, user, password, dbname, port)

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
	if ses.Values["accout"] != nil {
		ac = ses.Values["account"].(string)
	}

	var user my.User
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.Where("account = ?", ac).First(&user)

	return &user
}

// Template for no-template
func notemp() *template.Template {
	tmp, _ := template.New("index").Parse("NO PAGE.")
	return tmp
}

// get target Template.
func page(fname string) *template.Template {
	tmps, _ := template.ParseFiles("templates/"+fname+".html", "template/head.html", "templates/foot.html")
	return tmps
}

// top page handler
func index(w http.ResponseWriter, rq *http.Request) {
	user := chackLogin(w, rq)

	db, er := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	var pl []my.Post
	db.Where("group_id > 0").Order("created_at desc").Limit(10).Find(&pl)
	var gl []my.Group
	db.Order("created_at desc").Limit(10).Find(&gl)

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
	err := page("index").Excute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

// top page handler
func post(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)

	pid := rq.FormValue("pid")
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if rq.Method == "POST" {
		msg := rq.PostFormValue("message")
		pId, _ := strconv.Atoi(pid)
		cmt := my.Commet{
			UserId:  int(user.Model.ID),
			PostId:  pId,
			Message: msg,
		}
		db.Create(&cmt)
	}

	var pst my.Post
	var cmts []my.CommentJoin

	db.Where("id = ?", pid).First(&pst)
	db.Table("comments").Select("comments.*, user.id, users.name").Joins("join users on users.id = comments.user_id").Where("comments.post_id = ?", pid).Order("created_at desc").Find(&cmts)

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

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

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
			db.Create(&pt)
		case "group":
			gp := my.Group{
				UserId:  int(user.Model.ID),
				Name:    rq.PostFormValue("name"),
				Message: rq.PostFormValue("message"),
			}
			db.Create(&gp)
		}
	}

	var pts []my.Post
	var gps []my.Group

	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&gps)

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
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

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
		db.Create(&pt)
	}

	var grp my.Group
	var pts []my.Post

	db.Where("id=?", gid).First(&grp)
	db.Order("created_at desc").Model(&grp).Related(&pts)

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
	}
	if rq.Method == "POST" {
		db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		usr := rq.PostFormValue("account")
		pass := rq.PostFormValue("pass")
		item.Account = usr

		// check account and password
		var re int64
		var user my.User

		db.Where("account = ? and password = ?", usr, pass).Find(&user).Count(&re)

		if re <= 0 {
			item.Message = "Wrong account or password."
			page("login").Execute(w, item)
			return
		}

		// logined.
		ses, _ := cs.Get(rq, sesName)
		ses.Values["login"] = true
		ses.Values["account"] = usr
		ses.Valuse["name"] = user.Name
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
	ses.Value["login"] = nil
	ses.Value["account"] = nil
	ses.Save(rq, w)
	http.Redirect(w, rq, "/login", 302)
}

func main() {
	loadEnv(envpath)

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
	http.HandleFunc("/logout", func(w hgit ttp.ResponseWriter, rq *http.Request) {
		logout(w, rq)
	})

	http.ListenAndServe("", nil)
}

func loadEnv(envpath string) {
	err := godotenv.Load(envpath)
	if err != nil {
		fmt.Println("error")
	}
}
