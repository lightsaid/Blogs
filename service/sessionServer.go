package service

import "github.com/lightsaid/blogs/dbrepo"

type SessionServer struct {
	store dbrepo.SessionRepo
}

func NewSessionServer(store dbrepo.SessionRepo) *SessionServer {
	return &SessionServer{
		store: store,
	}
}
