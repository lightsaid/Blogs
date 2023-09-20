package service

import (
	"context"
	"log/slog"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/errs"
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
func (srv *PostsServer) NewPosts(ctx context.Context, req forms.NewPostsRequest) (int64, *errs.AppError) {

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
		CoverID:    &req.CoverID,
		Tags:       tags,
		Categories: categories,
	}

	newID, err := srv.store.Save(ctx, &posts)
	if err != nil {
		slog.Error(err.Error())
		return 0, errs.HandleSQLError(err)
	}
	return newID, nil
}

// UpdatePosts 更新 posts
func (srv *PostsServer) UpdatePosts(ctx context.Context, req forms.UpdatePostsRequest) *errs.AppError {
	tags := []*models.Tag{}
	for _, id := range req.TagIDs {
		tags = append(tags, &models.Tag{Model: models.Model{ID: id}})
	}

	categories := []*models.Category{}
	for _, id := range req.CategoryIDs {
		categories = append(categories, &models.Category{Model: models.Model{ID: id}})
	}

	posts, err := srv.store.Get(ctx, req.ID)
	if err != nil {
		return errs.HandleSQLError(err)
	}

	if req.Title != "" {
		posts.Title = req.Title
	}

	if req.Content != "" {
		posts.Content = req.Content
	}

	posts.Abstract = req.Abstract
	posts.Keyword = req.Keyword
	posts.CoverID = &req.CoverID
	posts.Slug = req.Slug
	posts.Tags = tags
	posts.Categories = categories

	_, err = srv.store.Save(ctx, posts)
	if err != nil {
		slog.Error(err.Error())
		return errs.HandleSQLError(err)
	}

	return nil
}

// GetPostsDetail 获取文章详情，包括 category、tags
func (srv *PostsServer) GetPostsDetail(ctx context.Context, id int64) (*models.Posts, *errs.AppError) {
	posts, err := srv.store.GetDetail(ctx, id)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}
	return posts, nil
}

// GetList 获取列表
func (srv *PostsServer) GetList(ctx context.Context, page, pageSize int) (*forms.ListRequest, *errs.AppError) {
	query := dbrepo.Filters{
		Page:     page,
		PageSize: pageSize,
	}
	list, meta, err := srv.store.List(ctx, query)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	data := &forms.ListRequest{
		Meta: meta,
		List: list,
	}

	return data, nil
}

// GetListByCategoryID 获取列表
func (srv *PostsServer) GetListByCategoryID(ctx context.Context, categoryID int64, page, pageSize int) (*forms.ListRequest, *errs.AppError) {
	query := dbrepo.Filters{
		Page:     page,
		PageSize: pageSize,
	}

	list, meta, err := srv.store.GetListByCategoryID(ctx, categoryID, query)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	data := &forms.ListRequest{
		Meta: meta,
		List: list,
	}

	return data, nil
}

// GetListByTagID 获取列表
func (srv *PostsServer) GetListByTagID(ctx context.Context, tagID int64, page, pageSize int) (*forms.ListRequest, *errs.AppError) {
	query := dbrepo.Filters{
		Page:     page,
		PageSize: pageSize,
	}

	list, meta, err := srv.store.GetListByTagID(ctx, tagID, query)
	if err != nil {
		return nil, errs.HandleSQLError(err)
	}

	data := &forms.ListRequest{
		Meta: meta,
		List: list,
	}

	return data, nil
}
