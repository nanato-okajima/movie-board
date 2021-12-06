package main

import "youtube/my"

type Item struct {
	Title   string
	Message string
	Name    string
	Account string
	Post    my.Post
	Clist   []my.CommentJoin
}
