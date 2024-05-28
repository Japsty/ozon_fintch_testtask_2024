package models

import "context"

type Comment struct {
	ID              string     `json:"id"`
	Content         string     `json:"content"`
	AuthorID        string     `json:"AuthorID"`
	PostID          string     `json:"postId"`
	ParentCommentID string     `json:"parentCommentID,omitempty"`
	Replies         []*Comment `json:"replies"`
}

type CommentRepo interface {
	CreateComment(context.Context, string, string, string, string, string) ([]Comment, error)
	//UpdateComment(string, string) error
	DeleteCommentByID(context.Context, string, string, string) ([]Comment, error)
}
