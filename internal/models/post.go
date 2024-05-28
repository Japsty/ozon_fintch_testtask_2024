package models

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Comments        []*Comment `json:"comments"`
	CommentsAllowed bool       `json:"commentsAllowed"`
}

type PostRepo interface {
	GetAllPosts() ([]Post, error)
	CreatePost(string, string, string, bool) (Post, error)
	GetPostByPostID(string) (Post, error)
	UpdatePostComments(string, bool) (Post, error)
	DeletePostByID(string) error
}
