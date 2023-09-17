package service

import (
	"context"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/models"
	"github.com/lightsaid/blogs/routers/forms"
)

type CategoryServer struct {
	store dbrepo.CategoryRepo
}

func NewCategoryServer(store dbrepo.CategoryRepo) *CategoryServer {
	return &CategoryServer{
		store: store,
	}
}

func (srv *CategoryServer) Insert(ctx context.Context, title string) (int64, *dbrepo.DBError) {
	category := models.Category{
		Title: title,
		Slug:  "",
	}

	newID, err := srv.store.Insert(ctx, &category)
	if err != nil {
		return 0, dbrepo.CheckError(err)
	}

	return newID, nil
}

func (srv *CategoryServer) Update(ctx context.Context, req forms.UpdateCategoryRequest) *dbrepo.DBError {
	category := models.Category{
		Model: models.Model{ID: req.ID},
		Title: req.Title,
		Slug:  "",
	}

	err := srv.store.Update(ctx, &category)
	if err != nil {
		return dbrepo.CheckError(err)
	}

	return nil
}

func (srv *CategoryServer) Delete(ctx context.Context, id int64) *dbrepo.DBError {
	err := srv.store.Delete(ctx, id)
	if err != nil {
		return dbrepo.CheckError(err)
	}

	return nil
}

func (srv *CategoryServer) List(ctx context.Context) ([]*models.Category, *dbrepo.DBError) {
	list, err := srv.store.GetAll(ctx)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	return list, nil
}

func (srv *CategoryServer) Get(ctx context.Context, id int64) (*models.Category, *dbrepo.DBError) {
	category, err := srv.store.Get(ctx, id)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	return category, nil
}
