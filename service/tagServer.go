package service

import (
	"context"

	"github.com/lightsaid/blogs/dbrepo"
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

func (srv *TagServer) Insert(ctx context.Context, title string) (int64, *dbrepo.DBError) {
	tag := models.Tag{
		Title: title,
		Slug:  "",
	}

	newID, err := srv.store.Insert(ctx, &tag)
	if err != nil {
		return 0, dbrepo.CheckError(err)
	}

	return newID, nil
}

func (srv *TagServer) Update(ctx context.Context, req forms.UpdateTagRequest) *dbrepo.DBError {
	tag := models.Tag{
		Model: models.Model{ID: req.ID},
		Title: req.Title,
		Slug:  "",
	}

	err := srv.store.Update(ctx, &tag)
	if err != nil {
		return dbrepo.CheckError(err)
	}

	return nil
}

func (srv *TagServer) Delete(ctx context.Context, id int64) *dbrepo.DBError {
	err := srv.store.Delete(ctx, id)
	if err != nil {
		return dbrepo.CheckError(err)
	}

	return nil
}

func (srv *TagServer) List(ctx context.Context) ([]*models.Tag, *dbrepo.DBError) {
	list, err := srv.store.GetAll(ctx)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	return list, nil
}

func (srv *TagServer) Get(ctx context.Context, id int64) (*models.Tag, *dbrepo.DBError) {
	tag, err := srv.store.Get(ctx, id)
	if err != nil {
		return nil, dbrepo.CheckError(err)
	}

	return tag, nil
}
