package service

import "github.com/lightsaid/blogs/dbrepo"

type CommentServer struct {
	store dbrepo.CommentRepo
}

func NewCommentServer(store dbrepo.CommentRepo) *CommentServer {
	return &CommentServer{
		store: store,
	}
}
