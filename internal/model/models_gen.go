// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID              string     `json:"id"`
	Content         string     `json:"content"`
	Author          string     `json:"author"`
	PostID          string     `json:"postId"`
	ParentCommentID *string    `json:"parentCommentId,omitempty"`
	Replies         []*Comment `json:"replies,omitempty"`
	CreatedAt       string     `json:"createdAt"`
}

type Mutation struct {
}

type NewComment struct {
	PostID          string  `json:"postId"`
	Content         string  `json:"content"`
	Author          string  `json:"author"`
	ParentCommentID *string `json:"parentCommentId,omitempty"`
}

type NewPost struct {
	Title           string `json:"title"`
	Content         string `json:"content"`
	CommentsAllowed bool   `json:"commentsAllowed"`
}

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	UserID          string     `json:"userId"`
	Comments        []*Comment `json:"comments"`
	CommentsAllowed bool       `json:"commentsAllowed"`
	CreatedAt       string     `json:"createdAt"`
}

type Query struct {
}

type Subscription struct {
}
