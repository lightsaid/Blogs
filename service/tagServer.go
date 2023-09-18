package service

import (
	"context"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/errs"
	"github.com/lightsaid/blogs/models"
	"github.com/lightsaid/blogs/routers/forms"
)

type TagServer struct {
	store dbrepo.TagRepo
}

func NewTagServer(store dbrepo.TagRepo) *TagServer {
	return &TagServer{
		store: store,
	}
}

func (srv *TagServer) Insert(ctx context.Context, title string) (int64, *errs.AppError) {
	tag := models.Tag{
		Title: title,
		Slug:  "",
	}

	newID, err := srv.store.Insert(ctx, &tag)
	if err != nil {
		return 0, errs.HandleSQLError(err)
	}

	return newID, nil
}

func (srv *TagServer) Update(ctx context.Context, req forms.UpdateTagRequest) *errs.AppError {
	tag := models.Tag{
		Model: models.Model{ID: req.ID},
		Title: req.Title,
		Slug:  "",
	}

	err := srv.store.Update(ctx, &tag)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	return nil
}

func (srv *TagServer) Delete(ctx context.Context, id int64) *errs.AppError {
	err := srv.store.Delete(ctx, id)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	return nil
}

func (srv *TagServer) List(ctx context.Context) ([]*models.Tag, *errs.AppError) {
	list, err := srv.store.GetAll(ctx)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	return list, nil
}

func (srv *TagServer) Get(ctx context.Context, id int64) (*models.Tag, *errs.AppError) {
	tag, err := srv.store.Get(ctx, id)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	return tag, nil
}
