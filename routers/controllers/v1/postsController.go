package v1

import (
	"net/http"
	"strconv"

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
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"id": newID, "msg": successText}

	successResponse(w, r, data)
}

// Update 更新文章，包括和category、tags的依赖关系
func (p *PostsController) Update(w http.ResponseWriter, r *http.Request) {
	var req forms.UpdatePostsRequest

	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	err := p.server.UpdatePosts(r.Context(), req)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"msg": successText}

	successResponse(w, r, data)
}

// GetDetail 获取文章详情包括 category、tags
func (p *PostsController) GetDetail(w http.ResponseWriter, r *http.Request) {
	id, ok := bindParamInt64(w, r, "id")
	if !ok {
		return
	}

	posts, err := p.server.GetPostsDetail(r.Context(), id)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"data": posts, "msg": successText}
	successResponse(w, r, data)
}

func (p *PostsController) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page")) // 第几页
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")) // 每页几条
	if err != nil {
		pageSize = 10
	}

	data, err2 := p.server.GetList(r.Context(), page, pageSize)
	if err2 != nil {
		errorResponse(w, r, err2.StatusCode(), err2, err2.Message())
		return
	}

	successResponse(w, r, envelop{"data": data, "msg": successText})
}

func (p *PostsController) ListByTagID(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page")) // 第几页
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")) // 每页几条
	if err != nil {
		pageSize = 10
	}
	tagID, ok := bindParamInt64(w, r, "id")
	if !ok {
		return
	}

	data, err2 := p.server.GetListByTagID(r.Context(), tagID, page, pageSize)
	if err2 != nil {
		errorResponse(w, r, err2.StatusCode(), err2, err2.Message())
		return
	}

	successResponse(w, r, envelop{"data": data, "msg": successText})
}

func (p *PostsController) ListByCategoryID(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page")) // 第几页
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")) // 每页几条
	if err != nil {
		pageSize = 10
	}
	categoryID, ok := bindParamInt64(w, r, "id")
	if !ok {
		return
	}

	data, err2 := p.server.GetListByCategoryID(r.Context(), categoryID, page, pageSize)
	if err2 != nil {
		errorResponse(w, r, err2.StatusCode(), err2, err2.Message())
		return
	}

	successResponse(w, r, envelop{"data": data, "msg": successText})
}
