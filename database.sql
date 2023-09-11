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
    foreign key (tag_id) references tags(id)
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
    created_at text not null,
    expired_at text not null,
    foreign key(user_id) references users(id)
);

