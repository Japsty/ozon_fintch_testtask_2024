package resolvers

import (
	"Ozon_testtask/graph"
	"Ozon_testtask/internal/model"
	mocks "Ozon_testtask/tests/mocks/services"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"
)

var mockNewPost = model.NewPost{
	Title:           "mock title",
	Content:         "mock content",
	CommentsAllowed: true,
}

var mockPost = &model.Post{
	ID:              "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Title:           mockNewPost.Title,
	Content:         mockNewPost.Content,
	CommentsAllowed: mockNewPost.CommentsAllowed,
	CreatedAt:       time.Now().String(),
}

var mockComment = &model.Comment{
	ID:        "24ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "mock title comment",
	AuthorID:  "mock content comment",
	PostID:    "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  nil,
	Replies:   nil,
	CreatedAt: time.Now().String(),
}

type AddPostResponse struct {
	ID              string           `json:"id"`
	Title           string           `json:"title"`
	Content         string           `json:"content"`
	UserID          string           `json:"userID"`
	CommentsAllowed bool             `json:"commentsAllowed"`
	CreatedAt       string           `json:"createdAt"`
	Comments        []*model.Comment `json:"comments"`
}

var addPostResp struct {
	AddPost AddPostResponse `json:"addPost"`
}

func TestAddPost(t *testing.T) {
	testCases := []struct {
		id       int
		name     string
		mockPost *model.Post
		query    string
		response interface{}
		isError  bool
	}{{
		id:       1,
		name:     "Success",
		mockPost: mockPost,
		query: `
		mutation AddPost {
			addPost(post: { title: "mock title", content: "mock content", commentsAllowed: true }) {
				id
				title
				content
				userID
				commentsAllowed
				createdAt
				comments {
					id
					content
					authorID
					postID
					parentID
					createdAt
				}
			}
		}
		`,
		response: addPostResp,
		isError:  false,
	},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostService := new(mocks.PostServiceMock)
			mockCommentService := new(mocks.CommentServiceMock)
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			logger := zapLogger.Sugar()

			resolver := graph.NewResolver(mockPostService, mockCommentService, logger)
			c := client.New(handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver})))

			mockPostService.On("AddPost", mock.Anything, mockNewPost.Title, mockNewPost.Content, mockNewPost.CommentsAllowed).Return(*tc.mockPost, nil)

			err = c.Post(tc.query, &tc.response)
			if err != nil {
				return
			}
			mockPostService.AssertExpectations(t)
		})
	}
}

func TestAddComment(t *testing.T) {
	testCases := []struct {
		id           int
		name         string
		mockPost     *model.Post
		mockComments []*model.Comment
		query        string
		response     interface{}
		isError      bool
	}{{
		id:           1,
		name:         "Success",
		mockPost:     mockPost,
		mockComments: []*model.Comment{mockComment},
		query: `
		mutation AddComment {
			addComment(
				comment: {
					postID: "a2223908-47ff-4a9e-a775-8bea7ded7652"
					content: "mock comment"
					parentID: "\"\""
				}
			) {
				id
				title
				content
				userID
				commentsAllowed
				createdAt
				comments {
					id
					content
					authorID
					postID
					parentID
					createdAt
				}
			}
		}
		`,
		response: addPostResp,
		isError:  false,
	},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostService := new(mocks.PostServiceMock)
			mockCommentService := new(mocks.CommentServiceMock)
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			logger := zapLogger.Sugar()

			resolver := graph.NewResolver(mockPostService, mockCommentService, logger)
			c := client.New(handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver})))

			mockCommentService.On("CommentPost", mock.Anything, mockComment.PostID, mockComment.Content).Return(tc.mockComments, nil)
			mockPostService.On("GetPostByPostID", mock.Anything, mockComment.PostID).Return(*tc.mockPost, nil)

			err = c.Post(tc.query, &tc.response)
			if err != nil {
				return
			}

			mockPostService.AssertExpectations(t)
		})
	}
}

func TestPost(t *testing.T) {
	testCases := []struct {
		id           int
		name         string
		mockPost     *model.Post
		mockComments []*model.Comment
		query        string
		response     interface{}
		isError      bool
	}{{
		id:           1,
		name:         "Success",
		mockPost:     mockPost,
		mockComments: []*model.Comment{mockComment},
		query: `
		query Post {
			post(id: "14ad7024-7c45-4453-9fac-2dfae1ad2c96", limit: 1, offset: 0) {
				id
				title
				content
				userID
				commentsAllowed
				createdAt
				comments {
					id
					content
					authorID
					postID
					parentID
					createdAt
				}
			}
		}
		`,
		response: addPostResp,
		isError:  false,
	},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostService := new(mocks.PostServiceMock)
			mockCommentService := new(mocks.CommentServiceMock)
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			logger := zapLogger.Sugar()

			resolver := graph.NewResolver(mockPostService, mockCommentService, logger)
			c := client.New(handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver})))

			mockPostService.On("GetPostByPostID", mock.Anything, mockPost.ID).Return(*tc.mockPost, nil)
			mockCommentService.On("GetCommentsByPostID", mock.Anything, mockComment.PostID, 1, 0).Return(tc.mockComments, nil)

			err = c.Post(tc.query, &tc.response)
			if err != nil {
				return
			}
			mockPostService.AssertExpectations(t)
		})
	}
}
