package service

import (
	"context"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/errs"
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

func (srv *CategoryServer) Insert(ctx context.Context, title string) (int64, *errs.AppError) {
	category := models.Category{
		Title: title,
		Slug:  "",
	}

	newID, err := srv.store.Insert(ctx, &category)
	if err != nil {
		return 0, errs.HandleSQLError(err)
	}

	return newID, nil
}

func (srv *CategoryServer) Update(ctx context.Context, req forms.UpdateCategoryRequest) *errs.AppError {
	category := models.Category{
		Model: models.Model{ID: req.ID},
		Title: req.Title,
		Slug:  "",
	}

	err := srv.store.Update(ctx, &category)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	return nil
}

func (srv *CategoryServer) Delete(ctx context.Context, id int64) *errs.AppError {
	err := srv.store.Delete(ctx, id)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	return nil
}

func (srv *CategoryServer) List(ctx context.Context) ([]*models.Category, *errs.AppError) {
	list, err := srv.store.GetAll(ctx)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	return list, nil
}

func (srv *CategoryServer) Get(ctx context.Context, id int64) (*models.Category, *errs.AppError) {
	category, err := srv.store.Get(ctx, id)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	return category, nil
}
