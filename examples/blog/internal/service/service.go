package service

import (
	"context"

	"github.com/meiigo/gkit/examples/blog/api"
	"github.com/meiigo/gkit/log"
)

type BlogService struct {
	api.UnimplementedBlogServer

	log     *log.Helper
	article *ArticleBiz
}

func NewBlogService(article *ArticleBiz, l log.Logger) *BlogService {
	return &BlogService{
		article: article,
		log:     log.NewHelper(l),
	}
}

func (s *BlogService) CreateArticle(ctx context.Context, req *api.CreateArticleRequest) (*api.CreateArticleReply, error) {
	s.log.Infof("input data %v", req)
	//err := s.article.Create(ctx, &biz.Article{
	//	Title:   req.Title,
	//	Content: req.Content,
	//})
	return &api.CreateArticleReply{}, nil
}

func (s *BlogService) GetArticle(ctx context.Context, req *api.GetArticleRequest) (*api.GetArticleReply, error) {
	s.log.Infof("input data %v", req)
	//err := s.article.Create(ctx, &biz.Article{
	//	Title:   req.Title,
	//	Content: req.Content,
	//})
	return &api.GetArticleReply{
		Article: &api.Article{
			Id:      1,
			Title:   "",
			Content: "",
			Like:    999,
		},
	}, nil
}
