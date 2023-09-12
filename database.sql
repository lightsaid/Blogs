-- 用户表
drop table if exists users;
create table users (
    id integer not null primary key autoincrement,
    email text not null unique,
    password text not null,
    username text not null,
    avatar text not null default '',
    role integer not null default 0, -- 角色：0:普通用户，1:管理员
    activated_at text,
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text
);

-- 文章表
drop table if exists posts;
create table posts (
    id integer not null primary key autoincrement,
    author_id integer not null,
    title text not null,
    content text not null,
    keyword text not null default '',  -- SEO 搜索
    slug text not null default '',   -- SEO 文章地址
    abstract text not null default '',  -- 文章摘要 
    cover_image_id integer, -- 封面图id
    views integer not null default 0,  -- 查看人数
    likes integer not null default 0, -- 点赞数
    comments integer not null default 0, -- 评论数
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text,
    foreign key(author_id) references users(id),
    foreign key(cover_image_id) references assets(id)
);


-- 分类表
drop table if exists category;
create table category (
    id integer not null primary key autoincrement,
    title text not null unique,
    slug text text not null default '', -- tag 对应文章列表链接，利于EO 
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text
);

-- 文章与分类关联表
drop table if exists posts_category;
create table posts_category (
    posts_id integer not null,
    category_id integer not null,
    foreign key (posts_id) references posts(id),
    foreign key (category_id) references tags(id)
);


-- 标签表
drop table if exists tags;
create table tags (
    id integer not null primary key autoincrement,
    title text not null unique,
    slug text not null default '', -- tag 对应文章列表链接，利于EO 
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text
);

-- 文章与标签关联表
drop table if exists posts_tag;
create table posts_tag (
    posts_id integer not null,
    tag_id integer not null,
    foreign key (posts_id) references posts(id),
    foreign key (tag_id) references tags(id)
);

-- 文件表
drop table if exists assets;
create table assets (
    id integer not null primary key autoincrement,
    user_id integer not null,
    posts_id integer,
    data blob not null,
    ext text, -- 文件后缀（如：.png）
    name text, -- 文件名
    size integer not null default 0, -- 文件大小
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text,
    foreign key(user_id) references users(id),
    foreign key (posts_id) references posts(id)
);

-- sessions 会话表
drop table if exists sessions;
CREATE TABLE sessions (
    id text not null primary key,
    user_id integer not null,
    refresh_token text not null unique,
    client_ip text not null default '',
    created_at text not null default (datetime('now', 'localtime')),
    expired_at text not null,
    foreign key(user_id) references users(id)
);

-- 评论表分两情况，一种有账号用户，无账号用户留言
drop table if exists comments;
create table comments (
    id integer not null primary key autoincrement,
    posts_id integer not null,
    user_id integer,
    parent_id integer,
    content text not null,
    nickname text,
    email text,
    created_at text not null default (datetime('now', 'localtime')),
    updated_at text not null default (datetime('now', 'localtime')),
    deleted_at text,
    foreign key(posts_id) references posts(id),
    foreign key(user_id) references users(id)
);


-- 创建2个用户，以后在开放注册 password: abc123
insert into users (email, password, username, avatar, role, activated_at) values(
    "xqq@qq.com", 
    "$2a$12$y2Yh9B.s.oqkzyaNXQ8ANO2kwlqyO7fJvQIXVGhkWhYlXQxce/Lfm",
    "xqq",
    "http://",
    1,
    datetime('now', 'localtime')
),(
    "xzz@qq.com", 
    "$2a$12$y2Yh9B.s.oqkzyaNXQ8ANO2kwlqyO7fJvQIXVGhkWhYlXQxce/Lfm",
    "xzz",
    "http://",
    0,
    datetime('now', 'localtime')
);

select *from users;

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
from -- 为了保证分页正确
   (
        select count(*) over() as totalRecords, * from posts limit 5 offset 0
   ) p
join 
    posts_tag pt on pt.posts_id = p.id
join 
    tags t on  t.id = pt.tag_id
where 
    p.deleted_at is null
group by p.id, t.id
order by p.created_at DESC;


select 
		t.id, t.title, t.slug, t.created_at, t.updated_at 
	from posts_tag pt
	join posts p on p.id=pt.posts_id
	join tags t on t.id=pt.tag_id
	where p.id = 1 and p.deleted_at is null;
