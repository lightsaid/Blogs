package v1

import (
	"log/slog"
	"net/http"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/routers/forms"
	"github.com/lightsaid/blogs/service"
)

type PostsController struct {
	server *service.PostsServer
}

func NewPostsController(store dbrepo.PostsRepo) *PostsController {
	return &PostsController{
		server: service.NewPostsServer(store),
	}
}

func (p *PostsController) Add(w http.ResponseWriter, r *http.Request) {
	var req forms.NewPostsRequest

	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	newID, err := p.server.NewPosts(r.Context(), req)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("创建posts", slog.Int64("id", newID))
}
