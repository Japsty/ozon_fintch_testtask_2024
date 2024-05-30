// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	AuthorID  string     `json:"authorID"`
	PostID    string     `json:"postID"`
	ParentID  string     `json:"parentID"`
	Replies   []*Comment `json:"replies,omitempty"`
	CreatedAt string     `json:"createdAt"`
}

type Mutation struct {
}

type NewComment struct {
	PostID   string `json:"postID"`
	ParentID string `json:"parentID"`
	Content  string `json:"content"`
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
	UserID          string     `json:"userID"`
	Comments        []*Comment `json:"comments"`
	CommentsAllowed bool       `json:"commentsAllowed"`
	CreatedAt       string     `json:"createdAt"`
}

type Query struct {
}

type Subscription struct {
}
