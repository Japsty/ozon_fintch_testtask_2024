package model

import (
	"context"
)

//type DBComment struct {
//	ID              string     `json:"id"`
//	Content         string     `json:"content"`
//	AuthorID        string     `json:"AuthorID"`
//	PostID          string     `json:"postId"`
//	ParentCommentID string     `json:"parentCommentID,omitempty"`
//	CreatedAt       time.Time  `json:"created_at,omitempty"`
//	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
//	Replies         []*Comment `json:"replies"`
//}

type CommentRepo interface {
	CreateComment(context.Context, string, string, string, string, string) ([]*Comment, error)
	GetCommentByParentID(context.Context, string) ([]*Comment, error)
	GetCommentsByPostID(context.Context, string) ([]*Comment, error)
	GetCommentsByPostIDPaginated(context.Context, string, int, int) ([]*Comment, error)
}

type CommentService interface {
	CommentPost(context.Context, string, string, string) ([]*Comment, error)
	GetCommentByParentID(context.Context, string) ([]*Comment, error)
	GetCommentsByPostID(context.Context, string, int, int) ([]*Comment, error)
}
