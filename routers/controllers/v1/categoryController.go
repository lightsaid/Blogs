package v1

import (
	"net/http"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/routers/forms"
	"github.com/lightsaid/blogs/service"
)

type CategoryController struct {
	server *service.CategoryServer
}

func NewCategoryController(store dbrepo.CategoryRepo) *CategoryController {
	return &CategoryController{
		server: service.NewCategoryServer(store),
	}
}

func (c *CategoryController) Add(w http.ResponseWriter, r *http.Request) {
	var req forms.AddCategoryRequest
	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	newID, err := c.server.Insert(r.Context(), req.Title)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"id": newID, "msg": successText}
	successResponse(w, r, data)
}

func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	var req forms.UpdateCategoryRequest
	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	err := c.server.Update(r.Context(), req)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"id": req.ID, "msg": successText}
	successResponse(w, r, data)
}

func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := bindParamInt64(w, r, "id")
	if !ok {
		return
	}

	if err := c.server.Delete(r.Context(), int64(id)); err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"id": id, "msg": successText}
	successResponse(w, r, data)
}

func (c *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	list, err := c.server.List(r.Context())
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"msg": successText, "data": list}
	successResponse(w, r, data)
}

func (c *CategoryController) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := bindParamInt64(w, r, "id")
	if !ok {
		return
	}

	category, err := c.server.Get(r.Context(), id)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	data := envelop{"msg": successText, "data": category}
	successResponse(w, r, data)
}
