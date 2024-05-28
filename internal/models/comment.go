package models

type Comment struct {
	ID              string     `json:"id"`
	Content         string     `json:"content"`
	AuthorID        string     `json:"AuthorID"`
	PostID          string     `json:"postId"`
	ParentCommentID string     `json:"parentCommentID,omitempty"`
	Replies         []*Comment `json:"replies"`
}

type CommentRepo interface {
	CreateComment(string, string, string, string, string) ([]Comment, error)
	//UpdateComment(string, string) error
	DeleteCommentByID(string, string, string) ([]Comment, error)
}
