package v1

import (
	"log/slog"
	"net/http"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/request"
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

func (p *PostsController) NewPosts(w http.ResponseWriter, r *http.Request) {
	var req forms.NewPostsRequest
	if err := request.ReadJSON(w, r, &req); err != nil {
		slog.ErrorContext(r.Context(), "->>> "+err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		slog.Error(err.Error())
		return
	}

	newID, err := p.server.NewPosts(r.Context(), req)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("创建posts", slog.Int64("id", newID))
}
