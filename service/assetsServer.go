package service

import "github.com/lightsaid/blogs/dbrepo"

type AssetsServer struct {
	store dbrepo.AssetsRepo
}

func NewAssetsServer(store dbrepo.AssetsRepo) *AssetsServer {
	return &AssetsServer{
		store: store,
	}
}
