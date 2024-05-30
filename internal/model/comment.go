package model

import (
	"context"
)

type CommentRepo interface {
	CreateComment(context.Context, string, string, string, string, string) ([]*Comment, error)
	GetCommentByParentID(context.Context, string) ([]*Comment, error)
	GetCommentsByPostID(context.Context, string) ([]*Comment, error)
	GetCommentsByPostIDPaginated(context.Context, string, int, int) ([]*Comment, error)
}

type CommentService interface {
	CommentPost(context.Context, string, string) ([]*Comment, error)
	GetCommentByParentID(context.Context, string) ([]*Comment, error)
	GetCommentsByPostID(context.Context, string, int, int) ([]*Comment, error)
}
