package dbrepo

import (
	"context"
	"strings"

	"github.com/lightsaid/blogs/models"
)

// PostsRepo 定义Posts表操作方法
type PostsRepo interface {
	// Insert 保存一篇文章，如果有分类和tag也一同保存，posts.ID 存在，更新posts，若不存在则是新增一篇文章
	Save(ctx context.Context, posts *models.Posts) (int64, error)
	// InsertPosts 仅仅往 posts 表添加一条记录
	InsertPosts(ctx context.Context, posts *models.Posts) (int64, error)
	// UpdatePosts 仅仅更新 posts 表记录
	UpdatePosts(ctx context.Context, posts *models.Posts) error
	// Get 获取一篇文章
	Get(ctx context.Context, pid int64) (*models.Posts, error)
	// BlukInsertPostsCategory 批量插入记录到 posts_category 表
	BlukInsertPostsCategory(ctx context.Context, pcs []*models.PostsCategory) error
	// BlukInsertPostsTag 批量插入记录到 posts_tag 表
	BlukInsertPostsTag(ctx context.Context, pts []*models.PostsTag) error
	// DeletePostsCategory 删除文章分类关系记录
	DeletePostsCategory(ctx context.Context, pid int64) error
	// DeletePostsTag 删除文章标签关系记录
	DeletePostsTag(ctx context.Context, pid int64) error
	// List 获取posts列表, 包含tag关系
	List(ctx context.Context, filter Filters) ([]*models.Posts, Metadata, error)
	// GetDetail 获取一篇文章和对应的tags
	GetDetail(ctx context.Context, pid int64) (*models.Posts, error)
}

// 接口检查
var _ PostsRepo = (*postsRepo)(nil)

// postsRepo 实现 PostsRepo 接口
type postsRepo struct {
	DB Queryable
	utilRepo
}

// Insert 保存一篇文章，如果有分类和tag也一同保存，
// posts.ID 存在，更新posts，若不存在则是新增一篇文章
func (store *postsRepo) Save(ctx context.Context, posts *models.Posts) (int64, error) {
	var categoryIDs []int64
	for _, tag := range posts.Categories {
		categoryIDs = append(categoryIDs, tag.ID)
	}

	var tagIDs []int64
	for _, tag := range posts.Tags {
		tagIDs = append(tagIDs, tag.ID)
	}

	if posts.ID > 0 {
		// 更新操作
		err := store.update(ctx, posts, categoryIDs, tagIDs)
		return posts.ID, err
	}

	// 添加操作
	return store.create(ctx, posts, categoryIDs, tagIDs)
}

// List 获取posts列表, 包含tag关系
func (store *postsRepo) List(ctx context.Context, filter Filters) ([]*models.Posts, Metadata, error) {
	querySQL := `
	select 
		totalRecords,
		p.id, 
		p.title,
		p.content,
		p.keyword,
		p.slug,
		p.abstract,
		p.cover_image_id,
		p.views,
		p.likes,
		p.comments,
		p.created_at,
		p.updated_at,
		t.title,
		t.slug,
		t.created_at,
		t.updated_at
	from
	(
			select count(*) over() as totalRecords, * from posts limit $1 offset $2
	) p
	join 
		posts_tag pt on pt.posts_id = p.id
	join 
		tags t on  t.id = pt.tag_id
	where 
		p.deleted_at is null
	group by p.id, t.id
	order by p.created_at DESC;
	`
	var metaData Metadata
	list := make([]*models.Posts, 0)
	totalRecords := 0
	rows, err := store.DB.QueryContext(ctx, querySQL, filter.limit(), filter.offset())
	if err != nil {
		return list, metaData, err
	}
	defer rows.Close()

	// 查询会出现重posts复行，一个posts可以对应多个tag，定义一个根据posts.ID存在tag的变量
	postsMapTag := make(map[int64][]*models.Tag)

	for rows.Next() {
		var posts models.Posts
		var tag models.Tag
		err = rows.Scan(
			&totalRecords,
			&posts.ID,
			&posts.Title,
			&posts.Content,
			&posts.Keyword,
			&posts.Slug,
			&posts.Abstract,
			&posts.CoverID,
			&posts.Views,
			&posts.Likes,
			&posts.Comments,
			&posts.CreatedAt,
			&posts.UpdatedAt,
			&tag.ID,
			&tag.Title,
			&tag.Slug,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return list, metaData, err
		}

		postsMapTag[posts.ID] = append(postsMapTag[posts.ID], &tag)

		list = append(list, &posts)
	}

	// 整合数据，将 postsMapTag 数据提取出来，放到 list 对应的 posts.Tag
	for index := range list {
		tags, exists := postsMapTag[list[index].ID]
		if exists {
			list[index].Tags = tags
		}
	}

	metaData = calculateMetadata(totalRecords, filter.Page, filter.PageSize)
	return list, metaData, nil
}

// create 添加文章
func (store *postsRepo) create(ctx context.Context, posts *models.Posts, categoryIDs []int64, tagIDs []int64) (int64, error) {
	var newID int64

	//  执行事务操作
	err := store.execTx(ctx, store.DB, func(repo *Repository) error {
		// 插入文章
		pid, err := repo.PostsRepo.InsertPosts(ctx, posts)
		if err != nil {
			return err
		}
		newID = pid

		// 构造 添加关系表 数据（posts_category、posts_tag）
		pcs := make([]*models.PostsCategory, len(categoryIDs))
		pts := make([]*models.PostsTag, len(tagIDs))

		for index := range pcs {
			pcs[index] = &models.PostsCategory{PostsID: pid, CategoryID: categoryIDs[index]}
		}

		for index := range pts {
			pts[index] = &models.PostsTag{PostsID: pid, TagID: tagIDs[index]}
		}

		// 添加文章分类关系表
		err = repo.PostsRepo.BlukInsertPostsCategory(ctx, pcs)
		if err != nil {
			return err
		}

		// 添加文章标签关系表
		err = repo.PostsRepo.BlukInsertPostsTag(ctx, pts)
		if err != nil {
			return err
		}

		return nil
	})

	return newID, err
}

// update 更新文章
func (store *postsRepo) update(ctx context.Context, posts *models.Posts, categoryIDs []int64, tagIDs []int64) error {
	return store.execTx(ctx, store.DB, func(repo *Repository) error {
		pid := posts.ID

		// 更新 posts 记录
		err := repo.PostsRepo.UpdatePosts(ctx, posts)
		if err != nil {
			return err
		}

		// 删除 posts、category、tag 关系
		err = repo.PostsRepo.DeletePostsCategory(ctx, pid)
		if err != nil {
			return err
		}

		err = repo.PostsRepo.DeletePostsTag(ctx, pid)
		if err != nil {
			return err
		}

		// 构造 添加关系表 数据（posts_category、posts_tag）
		pcs := make([]*models.PostsCategory, len(categoryIDs))
		pts := make([]*models.PostsTag, len(tagIDs))

		for index := range pcs {
			pcs[index] = &models.PostsCategory{PostsID: pid, CategoryID: categoryIDs[index]}
		}

		for index := range pts {
			pts[index] = &models.PostsTag{PostsID: pid, TagID: tagIDs[index]}
		}

		// 重新添加文章分类关系表
		err = repo.PostsRepo.BlukInsertPostsCategory(ctx, pcs)
		if err != nil {
			return err
		}

		// 重新添加文章标签关系表
		err = repo.PostsRepo.BlukInsertPostsTag(ctx, pts)
		if err != nil {
			return err
		}

		return nil
	})
}

// InsertPosts 仅仅往 posts 表添加一条记录
func (store *postsRepo) InsertPosts(ctx context.Context, posts *models.Posts) (int64, error) {
	querySQL := `insert into posts(
		author_id, 
		title,
		content,
		keyword,
		slug,
		abstract,
		cover_image_id
	)values(
		$1, $2, $3, $4, $5, $6, $7
	);`

	result, err := store.DB.ExecContext(
		ctx,
		querySQL,
		posts.AuthorID,
		posts.Title,
		posts.Content,
		posts.Keyword,
		posts.Slug,
		posts.Abstract,
		posts.CoverID,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdatePosts 仅仅更新 posts 表记录
func (store *postsRepo) UpdatePosts(ctx context.Context, posts *models.Posts) error {
	querySQL := `update posts set
		title=$1,
		content=$2,
		keyword=$3,
		slug=$4,
		abstract=$5,
		cover_image_id=$6,
		views=$7,
		likes=$8,
		comments=$9,
		updated_at=datetime('now', 'localtime')
	where id=$10 and deleted_at is null;
	`
	result, err := store.DB.ExecContext(
		ctx,
		querySQL,
		posts.Title,
		posts.Content,
		posts.Keyword,
		posts.Slug,
		posts.Abstract,
		posts.CoverID,
		posts.Views,
		posts.Likes,
		posts.Comments,
		posts.ID,
	)
	if err != nil {
		return err
	}

	if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

// Get 获取一篇文章
func (store *postsRepo) Get(ctx context.Context, pid int64) (*models.Posts, error) {
	querySQL := `
		select
			id,
			author_id,
			title,
			content,
			keyword,
			slug,
			abstract,
			cover_image_id,
			views,
			likes,
			comments,
			created_at,
			updated_at
		from posts where id=$1 and deleted_at is null
	`
	var posts models.Posts
	err := store.DB.SelectContext(ctx, &posts, querySQL, pid)
	return &posts, err
}

// BlukInsertPostsCategory 批量插入记录到 posts_category 表
func (store *postsRepo) BlukInsertPostsCategory(ctx context.Context, pcs []*models.PostsCategory) error {
	querySQL := `insert into posts_category(posts_id, category_id) values`

	var values []string
	var params []interface{}
	for _, pc := range pcs {
		values = append(values, "($1, $2)")
		params = append(params, pc.PostsID, pc.CategoryID)
	}
	querySQL += strings.Join(values, ",")

	_, err := store.DB.ExecContext(ctx, querySQL, params...)
	if err != nil {
		return err
	}

	return nil
}

// BlukInsertPostsTag 批量插入记录到 posts_tag 表
func (store *postsRepo) BlukInsertPostsTag(ctx context.Context, pts []*models.PostsTag) error {
	querySQL := `insert into posts_tag(posts_id, tag_id) values`

	var values []string
	var params []interface{}
	for _, pt := range pts {
		values = append(values, "($1, $2)")
		params = append(params, pt.PostsID, pt.TagID)
	}
	querySQL += strings.Join(values, ",")

	_, err := store.DB.ExecContext(ctx, querySQL, params...)
	if err != nil {
		return err
	}

	return nil
}

// DeletePostsCategory 删除文章分类关系记录
func (store *postsRepo) DeletePostsCategory(ctx context.Context, pid int64) error {
	querySQL := `delete from posts_category where posts_id=$1`
	result, err := store.DB.ExecContext(ctx, querySQL, pid)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// DeletePostsTag 删除文章标签关系记录
func (store *postsRepo) DeletePostsTag(ctx context.Context, pid int64) error {
	querySQL := `delete from posts_tag where posts_id=$1`
	result, err := store.DB.ExecContext(ctx, querySQL, pid)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// GetDetail 获取一篇文章和对应的tags
func (store *postsRepo) GetDetail(ctx context.Context, pid int64) (*models.Posts, error) {
	posts, err := store.Get(ctx, pid)
	if err != nil {
		return nil, err
	}

	querySQL := `
	select 
		t.id, t.title, t.slug, t.created_at, t.updated_at 
	from posts_tag pt
	join posts p on p.id=pt.posts_id
	join tags t on t.id=pt.tag_id
	where p.id=$1 and p.deleted_at is null;
	`
	tags := []*models.Tag{}
	err = store.DB.SelectContext(ctx, tags, querySQL, posts.ID)
	if err != nil {
		return nil, err
	}

	posts.Tags = tags

	return posts, nil
}
