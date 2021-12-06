package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"youtube/config"
	"youtube/constants"
	"youtube/my"
)

var DB *gorm.DB

// db variable
func init() {
	conf := config.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// get target Template.
func page(fname string) *template.Template {
	tmps, _ := template.ParseFiles(constants.TEMPLATES_DIR+fname+".html", constants.TEMPLATE_HEAD, constants.TEMPLATE_FOOT)
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
func post(w http.ResponseWriter, r *http.Request) {
	user := checkLogin(w, r)

	pid := r.FormValue("pid")

	if r.Method == "POST" {
		msg := r.PostFormValue("message")
		pId, _ := strconv.Atoi(pid)
		my.CreatePost(user, msg, pId)
	}

	var pst my.Post
	var cmts []my.CommentJoin

	my.FindPostById(pid, &pst)
	my.SelectJoinedTable(&cmts, "comments", "comments.*, user.id, users.name", "join users on users.id = comments.user_id", "comments.post_id = ?", pid, "desc")

	item := Item{
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
