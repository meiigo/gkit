package service

import (
	"context"
	"time"
)

type Article struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Like      int64
}

type ArticleRepo interface {
	// db
	ListArticle(ctx context.Context) ([]*Article, error)
	GetArticle(ctx context.Context, id int64) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) error
	UpdateArticle(ctx context.Context, id int64, article *Article) error
	DeleteArticle(ctx context.Context, id int64) error

	// redis
	GetArticleLike(ctx context.Context, id int64) (rv int64, err error)
	IncArticleLike(ctx context.Context, id int64) error
}

type ArticleBiz struct {
	repo ArticleRepo
}
