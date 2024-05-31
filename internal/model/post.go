package model

import "context"

type PostRepo interface {
	GetAllPosts(context.Context) ([]*Post, error)
	CreatePost(context.Context, string, string, string, string, bool) (Post, error)
	GetPostByPostID(context.Context, string) (Post, error)
	UpdatePostStatus(context.Context, string, string, bool) (Post, error)
}

type PostService interface {
	GetAllPosts(context.Context) ([]*Post, error)
	AddPost(context.Context, string, string, bool) (Post, error)
	GetPostByPostID(context.Context, string) (Post, error)
	UpdatePostCommentsStatus(context.Context, string, bool) (Post, error)
}
