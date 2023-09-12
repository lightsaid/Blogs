package service

import "github.com/lightsaid/blogs/dbrepo"

type UserServer struct {
	store dbrepo.UserRepo
}

func NewUserServer(store dbrepo.UserRepo) *UserServer {
	return &UserServer{
		store: store,
	}
}

func (srv *UserServer) Login() {}
