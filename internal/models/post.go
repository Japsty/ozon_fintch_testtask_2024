package models

import "context"

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Comments        []*Comment `json:"comments"`
	CommentsAllowed bool       `json:"commentsAllowed"`
}

type PostRepo interface {
	GetAllPosts(context.Context) ([]Post, error)
	CreatePost(context.Context, string, string, string, bool) (Post, error)
	GetPostByPostID(context.Context, string) (Post, error)
	UpdatePostComments(context.Context, string, bool) (Post, error)
	DeletePostByID(context.Context, string) error
}
