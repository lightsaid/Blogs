package service

import (
	"context"
	"log/slog"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/models"
	"github.com/lightsaid/blogs/routers/forms"
)

type PostsServer struct {
	store dbrepo.PostsRepo
}

func NewPostsServer(store dbrepo.PostsRepo) *PostsServer {
	return &PostsServer{
		store: store,
	}
}

// NewPosts 创建 posts
func (srv *PostsServer) NewPosts(ctx context.Context, req forms.NewPostsRequest) (int64, error) {

	tags := []*models.Tag{}
	for _, id := range req.TagIDs {
		tags = append(tags, &models.Tag{Model: models.Model{ID: id}})
	}

	categories := []*models.Category{}
	for _, id := range req.TagIDs {
		categories = append(categories, &models.Category{Model: models.Model{ID: id}})
	}

	posts := models.Posts{
		AuthorID:   req.AuthorID,
		Title:      req.Title,
		Content:    req.Content,
		Keyword:    req.Keyword,
		Slug:       req.Slug,
		Abstract:   req.Abstract,
		CoverID:    req.CoverID,
		Tags:       tags,
		Categories: categories,
	}

	newID, err := srv.store.Save(ctx, &posts)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}
	return newID, nil
}

func (srv *PostsServer) UpdatePosts() {}
