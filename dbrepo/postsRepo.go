package dbrepo

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
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
	// GetDetail 获取一篇文章详情，包括 category、tags
	GetDetail(ctx context.Context, pid int64) (*models.Posts, error)
	// GetListByCategoryID 根据分类ID查询文章列表
	GetListByCategoryID(ctx context.Context, categoryID int64, filter Filters) ([]*models.Posts, Metadata, error)
	// GetListByTagID 根据TagID查询文章列表
	GetListByTagID(ctx context.Context, categoryID int64, filter Filters) ([]*models.Posts, Metadata, error)
}

// 接口检查
var _ PostsRepo = (*postsRepo)(nil)

// postsRepo 实现 PostsRepo 接口
type postsRepo struct {
	DB Queryable
	utilRepo
}

// Save 保存一篇文章，如果有分类和tag也一同保存，
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
	list, metaData, err := store.getListByPage(ctx, "", filter)
	if err != nil {
		return list, metaData, err
	}
	if len(list) == 0 {
		return list, metaData, nil
	}
	// 查询 posts 列表与之对应的 tags、category 关系
	list, err = store.getListDetail(ctx, list)
	if err != nil {
		return list, metaData, nil
	}
	return list, metaData, nil

	// 已被拆分为getListByPage
	// 分页查询，先得到posts的列表
	// limitQuery := `
	// select
	// 	count(*) over() as totalRecords,
	// 	id,
	// 	title,
	// 	content,
	// 	keyword,
	// 	slug,
	// 	abstract,
	// 	cover_image_id,
	// 	views,
	// 	likes,
	// 	comments,
	// 	created_at,
	// 	updated_at
	// from posts limit $1 offset $2;`
	// var list []*models.Posts
	// var postsIDs []int64
	// var totalRecords int
	// var metaData Metadata
	// rows, err := store.DB.QueryContext(ctx, limitQuery, filter.limit(), filter.offset())
	// if err != nil {
	// 	return nil, metaData, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var posts models.Posts
	// 	err = rows.Scan(
	// 		&totalRecords,
	// 		&posts.ID,
	// 		&posts.Title,
	// 		&posts.Content,
	// 		&posts.Keyword,
	// 		&posts.Slug,
	// 		&posts.Abstract,
	// 		&posts.CoverID,
	// 		&posts.Views,
	// 		&posts.Likes,
	// 		&posts.Comments,
	// 		&posts.CreatedAt,
	// 		&posts.UpdatedAt,
	// 	)
	// 	if err != nil {
	// 		return nil, metaData, err
	// 	}
	// 	postsIDs = append(postsIDs, posts.ID)
	// 	list = append(list, &posts)
	// }
	// // 计算 metaData
	// metaData = calculateMetadata(totalRecords, filter.Page, filter.PageSize)
	// if len(list) == 0 {
	// 	return list, metaData, sql.ErrNoRows
	// }
	// 已被拆分 为 getListDetail()
	// queryDetail := `
	// select
	// 	p.id,
	// 	t.id,
	// 	t.title,
	// 	t.slug,
	// 	t.created_at,
	// 	t.updated_at,
	// 	c.id,
	// 	c.title,
	// 	c.slug,
	// 	c.created_at,
	// 	c.updated_at
	// from
	// 	posts p
	// left join
	// 	posts_tag pt on pt.posts_id = p.id
	// left join
	// 	tags t on  t.id = pt.tag_id
	// left join
	// 	posts_category pc on pc.posts_id = p.id
	// left join category c on c.id = pc.category_id
	// where p.id in
	// `
	// queryDetail += "(" + strings.Join(ids, ",") + ")" + " order by t.id, p.id "
	// rows2, err := store.DB.QueryContext(ctx, queryDetail)
	// if err != nil {
	// 	return list, metaData, err
	// }
	// defer rows2.Close()
	// // 保存对应关系
	// posts2Tags := make(map[int64][]*models.Tag)
	// posts2Category := make(map[int64][]*models.Category)
	// // 防重 map map[posts.id] -> map[tag.id]
	// type mmps map[int64]struct{}
	// posts2TagsExists := make(map[int64]mmps)
	// posts2CategoryExists := make(map[int64]mmps)
	// for rows2.Next() {
	// 	var posts models.Posts
	// 	// NOTE: 这里用SQLColumn来接受扫描而不是用models.Tag/Catgory是因为SQL 语句 Left join 会存在 tag 或者 category 表空记录的情况，
	// 	// 这时候返回都是 NULL，因此原来的 models.Tag models.Category 都是非空的，如果用来了就会panic，这是不对的。因此才重新定义SQLColumn来接受扫描
	// 	var categoryCol SQLColumn
	// 	var tagCol SQLColumn
	// 	var tag models.Tag
	// 	var category models.Category
	// 	err = rows2.Scan(
	// 		&posts.ID,
	// 		&tagCol.ID,
	// 		&tagCol.Title,
	// 		&tagCol.Slug,
	// 		&tagCol.CreatedAt,
	// 		&tagCol.UpdatedAt,
	// 		&categoryCol.ID,
	// 		&categoryCol.Title,
	// 		&categoryCol.Slug,
	// 		&categoryCol.CreatedAt,
	// 		&categoryCol.UpdatedAt,
	// 	)
	// 	if err != nil {
	// 		return list, metaData, err
	// 	}
	// 	if posts2TagsExists[posts.ID] == nil {
	// 		posts2TagsExists[posts.ID] = make(map[int64]struct{})
	// 	}
	// 	if posts2CategoryExists[posts.ID] == nil {
	// 		posts2CategoryExists[posts.ID] = make(map[int64]struct{})
	// 	}
	// 	tag = tagCol.ToTag()
	// 	category = categoryCol.ToCategory()
	// 	// 这里的到结果，tag 很有可能每个字段都是空（零）值，因此必须 tag.ID > 0
	// 	if _, ok := posts2TagsExists[posts.ID][tag.ID]; !ok && tag.ID > 0 {
	// 		posts2TagsExists[posts.ID][tag.ID] = struct{}{}
	// 		posts2Tags[posts.ID] = append(posts2Tags[posts.ID], &tag)
	// 	}
	// 	// 同上
	// 	if _, ok := posts2CategoryExists[posts.ID][category.ID]; !ok && category.ID > 0 {
	// 		posts2CategoryExists[posts.ID][category.ID] = struct{}{}
	// 		posts2Category[posts.ID] = append(posts2Category[posts.ID], &category)
	// 	}
	// }
	// for index := range list {
	// 	postsID := list[index].ID
	// 	// 如果没数据，返回 空JSON数组： []
	// 	if posts2Tags[postsID] == nil {
	// 		posts2Tags[postsID] = make([]*models.Tag, 0)
	// 	}
	// 	// 同上
	// 	if posts2Category[postsID] == nil {
	// 		posts2Category[postsID] = make([]*models.Category, 0)
	// 	}
	// 	// 匹配 Tags， 设置到 list 上
	// 	if _, ok := posts2Tags[postsID]; ok {
	// 		list[index].Tags = posts2Tags[postsID]
	// 	}
	// 	// 同上
	// 	if _, ok := posts2Category[postsID]; ok {
	// 		list[index].Categories = posts2Category[postsID]
	// 	}
	// }
	// return list, metaData, nil
}

// GetPostsTags 获取一篇文章详情，包括 category、tags
func (store *postsRepo) GetDetail(ctx context.Context, pid int64) (*models.Posts, error) {
	//  使用 left join, inner join 是不对的
	query := `
	select 
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
		t.id,
		t.title,
		t.slug,
		t.created_at,
		t.updated_at,
		c.id,
		c.title,
		c.slug,
		c.created_at,
		c.updated_at
	from posts p 
	left join posts_tag pt on pt.posts_id  = p.id
	left join tags t on t.id = pt.tag_id 
	left join posts_category pc on pc.posts_id =p.id 
	left join category c on c.id  = pc.posts_id 
	where p.id = $1 limit 1
	`
	// NOTE: 当没有数据不会返回 sql.ErrNoRows
	rows, err := store.DB.QueryContext(ctx, query, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var posts models.Posts
	var posts2Category map[int64]map[int64]*models.Category
	var posts2Tags map[int64]map[int64]*models.Tag

	for rows.Next() {
		var tagCol SQLColumn
		var categoryCol SQLColumn

		err = rows.Scan(
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
			&tagCol.ID,
			&tagCol.Title,
			&tagCol.Slug,
			&tagCol.CreatedAt,
			&tagCol.UpdatedAt,
			&categoryCol.ID,
			&categoryCol.Title,
			&categoryCol.Slug,
			&categoryCol.CreatedAt,
			&categoryCol.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if tagCol.ID != nil && *tagCol.ID > 0 {
			if posts2Tags == nil {
				posts2Tags = make(map[int64]map[int64]*models.Tag)
				posts2Tags[posts.ID] = make(map[int64]*models.Tag)
			}
			tag := tagCol.ToTag()
			posts2Tags[posts.ID][tag.ID] = &tag
		}

		if categoryCol.ID != nil && *categoryCol.ID > 0 {
			if posts2Category == nil {
				posts2Category = make(map[int64]map[int64]*models.Category)
				posts2Category[posts.ID] = make(map[int64]*models.Category)
			}
			category := categoryCol.ToCategory()
			posts2Category[posts.ID][category.ID] = &category
		}
	}

	// sqlx 有返回空行这种场景，在dbeaver里执行sql没有，故此判断一下
	if posts.ID == 0 {
		return nil, sql.ErrNoRows
	}

	if len(posts2Tags[posts.ID]) > 0 {
		for _, tag := range posts2Tags[posts.ID] {
			tmp := tag
			posts.Tags = append(posts.Tags, tmp)
		}
	} else {
		posts.Tags = make([]*models.Tag, 0)
	}

	if len(posts2Category[posts.ID]) > 0 {
		for _, cate := range posts2Category[posts.ID] {
			tmp := cate
			posts.Categories = append(posts.Categories, tmp)
		}
	} else {
		posts.Categories = make([]*models.Category, 0)
	}
	return &posts, nil
}

// GetListByCategoryID 根据分类ID查询文章列表
func (store *postsRepo) GetListByCategoryID(ctx context.Context, categoryID int64, filter Filters) ([]*models.Posts, Metadata, error) {
	query := `join posts_category pc on pc.posts_id = p.id where pc.category_id =` + strconv.Itoa(int(categoryID)) + " order by p.updated_at DESC"
	list, metaData, err := store.getListByPage(ctx, query, filter)
	if err != nil {
		return list, metaData, err
	}

	// 去重
	uniqueMap := make(map[int64]struct{})
	var newList []*models.Posts
	for _, v := range list {
		if _, ok := uniqueMap[v.ID]; !ok {
			newList = append(newList, v)
		}
		uniqueMap[v.ID] = struct{}{}
	}

	// 查询 posts 列表与之对应的 tags、category 关系
	list, err = store.getListDetail(ctx, newList)
	if err != nil {
		return list, metaData, nil
	}
	return list, metaData, nil
}

// GetListByTagID 根据TagID查询文章列表
func (store *postsRepo) GetListByTagID(ctx context.Context, tagID int64, filter Filters) ([]*models.Posts, Metadata, error) {
	query := `join posts_tag pt on pt.posts_id = p.id where pt.tag_id =` + strconv.Itoa(int(tagID)) + " order by p.updated_at DESC"
	list, metaData, err := store.getListByPage(ctx, query, filter)
	if err != nil {
		return list, metaData, err
	}

	if len(list) == 0 {
		return list, metaData, nil
	}

	// 去重
	uniqueMap := make(map[int64]struct{})
	var newList []*models.Posts
	for _, v := range list {
		if _, ok := uniqueMap[v.ID]; !ok {
			newList = append(newList, v)
		}
		uniqueMap[v.ID] = struct{}{}
	}

	// 查询 posts 列表与之对应的 tags、category 关系
	list, err = store.getListDetail(ctx, newList)
	if err != nil {
		return list, metaData, nil
	}
	return list, metaData, nil
}

// getListByPage 定义一个方法快速查询 posts 列表数据，其中sql语句仅写了一半，另一半至于是join还是where看调用方需求使用，最后会拼接分页 “limit $1 offset $2;”
// 如果 querySQL 为空，则是默认查询posts列表, 注意SQL拼接，防止SQL注入
func (store *postsRepo) getListByPage(ctx context.Context, querySQL string, filter Filters) ([]*models.Posts, Metadata, error) {
	limitQuery := `
	select  
		count(*) over() as totalRecords,
		id,
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
	from posts p `

	limitQuery += querySQL + " limit $1 offset $2; "

	var list []*models.Posts
	var totalRecords int
	var metaData Metadata

	rows, err := store.DB.QueryContext(ctx, limitQuery, filter.limit(), filter.offset())
	if err != nil {
		return nil, metaData, err
	}
	defer rows.Close()

	for rows.Next() {
		var posts models.Posts
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
		)
		if err != nil {
			return nil, metaData, err
		}
		list = append(list, &posts)
	}

	// 计算 metaData
	metaData = calculateMetadata(totalRecords, filter.Page, filter.PageSize)
	if len(list) == 0 {
		return list, metaData, sql.ErrNoRows
	}

	return list, metaData, nil
}

// getListDetail 根据指定的llist获取每个posts的tag、category
func (store *postsRepo) getListDetail(ctx context.Context, list []*models.Posts) ([]*models.Posts, error) {
	// 查询 posts 与之对应的 tags、category 关系
	queryDetail := `
		select 
			p.id, 
			t.id,
			t.title,
			t.slug,
			t.created_at,
			t.updated_at,
			c.id,
			c.title,
			c.slug,
			c.created_at,
			c.updated_at
		from
			posts p
		left join 
			posts_tag pt on pt.posts_id = p.id
		left join 
			tags t on  t.id = pt.tag_id
		left join 
			posts_category pc on pc.posts_id = p.id 
		left join category c on c.id = pc.category_id 
		where p.id in 
	`

	ids := make([]string, len(list))
	for i, v := range list {
		ids[i] = strconv.Itoa(int(v.ID))
	}

	queryDetail += "(" + strings.Join(ids, ",") + ")" + " order by t.id, p.id "
	rows, err := store.DB.QueryContext(ctx, queryDetail)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	// 保存对应关系
	posts2Tags := make(map[int64][]*models.Tag)
	posts2Category := make(map[int64][]*models.Category)

	// 防重 map map[posts.id] -> map[tag.id]
	type mmps map[int64]struct{}
	posts2TagsExists := make(map[int64]mmps)
	posts2CategoryExists := make(map[int64]mmps)
	for rows.Next() {
		var posts models.Posts
		// NOTE: 这里用SQLColumn来接受扫描而不是用models.Tag/Catgory是因为SQL 语句 Left join 会存在 tag 或者 category 表空记录的情况，
		// 这时候返回都是 NULL，因此原来的 models.Tag models.Category 都是非空的，如果用来了就会panic，这是不对的。因此才重新定义SQLColumn来接受扫描
		var categoryCol SQLColumn
		var tagCol SQLColumn

		var tag models.Tag
		var category models.Category

		err = rows.Scan(
			&posts.ID,
			&tagCol.ID,
			&tagCol.Title,
			&tagCol.Slug,
			&tagCol.CreatedAt,
			&tagCol.UpdatedAt,
			&categoryCol.ID,
			&categoryCol.Title,
			&categoryCol.Slug,
			&categoryCol.CreatedAt,
			&categoryCol.UpdatedAt,
		)
		if err != nil {
			return list, err
		}

		if posts2TagsExists[posts.ID] == nil {
			posts2TagsExists[posts.ID] = make(map[int64]struct{})
		}

		if posts2CategoryExists[posts.ID] == nil {
			posts2CategoryExists[posts.ID] = make(map[int64]struct{})
		}

		tag = tagCol.ToTag()
		category = categoryCol.ToCategory()

		// 这里的到结果，tag 很有可能每个字段都是空（零）值，因此必须 tag.ID > 0
		if _, ok := posts2TagsExists[posts.ID][tag.ID]; !ok && tag.ID > 0 {
			posts2TagsExists[posts.ID][tag.ID] = struct{}{}
			posts2Tags[posts.ID] = append(posts2Tags[posts.ID], &tag)
		}

		// 同上
		if _, ok := posts2CategoryExists[posts.ID][category.ID]; !ok && category.ID > 0 {
			posts2CategoryExists[posts.ID][category.ID] = struct{}{}
			posts2Category[posts.ID] = append(posts2Category[posts.ID], &category)

		}
	}

	for index := range list {
		postsID := list[index].ID

		// 如果没数据，返回 空JSON数组： []
		if posts2Tags[postsID] == nil {
			posts2Tags[postsID] = make([]*models.Tag, 0)
		}

		// 同上
		if posts2Category[postsID] == nil {
			posts2Category[postsID] = make([]*models.Category, 0)
		}

		// 匹配 Tags， 设置到 list 上
		if _, ok := posts2Tags[postsID]; ok {
			list[index].Tags = posts2Tags[postsID]
		}

		// 同上
		if _, ok := posts2Category[postsID]; ok {
			list[index].Categories = posts2Category[postsID]
		}
	}

	return list, nil
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

// Get 获取一篇文章, 不包括 tag、category
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
	err := store.DB.GetContext(ctx, &posts, querySQL, pid)
	return &posts, err
}

// BlukInsertPostsCategory 批量插入记录到 posts_category 表
func (store *postsRepo) BlukInsertPostsCategory(ctx context.Context, pcs []*models.PostsCategory) error {
	querySQL := `insert into posts_category(posts_id, category_id) values`

	var values []string
	var params []interface{}
	for _, pc := range pcs {
		values = append(values, "(?, ?)") // 这里不能使用 $ 模式
		params = append(params, pc.PostsID, pc.CategoryID)
	}
	querySQL += strings.Join(values, ",")

	slog.InfoContext(ctx, "BlukInsertPostsCategory",
		slog.String("sql", querySQL), slog.Any("params", params))

	// 没有可插入数据
	if len(values) == 0 {
		return nil
	}

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
		values = append(values, "(?, ?)") // 这里不能使用 $ 模式
		params = append(params, pt.PostsID, pt.TagID)
	}
	querySQL += strings.Join(values, ",")

	// 没有可插入数据
	if len(values) == 0 {
		return nil
	}

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
