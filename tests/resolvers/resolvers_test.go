package resolvers

import (
	"Ozon_testtask/graph"
	"Ozon_testtask/internal/model"
	mocks "Ozon_testtask/tests/mocks/services"
	"bytes"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"strings"
	"testing"
	"time"
)

var mockNewPost = model.NewPost{
	Title:           "mockTitle",
	Content:         "mockContent",
	CommentsAllowed: true,
}

var mockPost = model.Post{
	ID:              "1",
	Title:           mockNewPost.Title,
	Content:         mockNewPost.Content,
	CommentsAllowed: mockNewPost.CommentsAllowed,
	CreatedAt:       time.Now().String(),
}

type mockRepoResp struct {
	mockPost  interface{}
	mockError error
}

//	type mockSessionResp struct {
//		mockSession *models.Session
//		mockError   error
//	}
type mockRequest struct {
	mockRequestMethod        string
	mockRequestURL           string
	mockRequestBody          *strings.Reader
	mockIncorrectRequestBody *bytes.Buffer
}

type testCase struct {
	id             int
	name           string
	mockRequest    mockRequest
	mockRepoResp   mockRepoResp
	expectedStatus int
}

func TestShouldAddPost(t *testing.T) {
	mockPostService := new(mocks.PostServiceMock)
	mockCommentService := new(mocks.CommentServiceMock)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	logger := zapLogger.Sugar()

	resolver := graph.NewResolver(mockPostService, mockCommentService, logger)

	newPost := model.NewPost{
		Title:           "mockTitle",
		Content:         "mockContent",
		CommentsAllowed: true,
	}

	createdPost := model.Post{
		ID:              "1",
		Title:           newPost.Title,
		Content:         newPost.Content,
		CommentsAllowed: newPost.CommentsAllowed,
		CreatedAt:       time.Now().String(),
	}

	mockPostService.On("AddPost", mock.Anything, newPost.Title, newPost.Content, newPost.CommentsAllowed).Return(createdPost, nil)

	c := client.New(handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver})))

	q := `mutation AddPost {
		addPost(
			post: { title: "mockTitle", content: "mockContent", commentsAllowed: true }
		) {
			id
			title
			content
			userID
			commentsAllowed
			createdAt
		}
	}
	`

	var response struct {
		Data struct {
			AddPost model.Post `json:"addPost"`
		} `json:"data"`
	}

	err = c.Post(q, &response)
	if err != nil {
		t.Fatalf("q: %s, response: %v", q, response)
		return
	}

	expectedPost := model.Post{
		ID:              "1",
		Title:           "mockTitle",
		Content:         "mockContent",
		CommentsAllowed: true,
		CreatedAt:       createdPost.CreatedAt,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedPost.ID, response.Data.AddPost.ID)
	assert.Equal(t, expectedPost.Title, response.Data.AddPost.Title)
	assert.Equal(t, expectedPost.Content, response.Data.AddPost.Content)
	assert.Equal(t, expectedPost.CommentsAllowed, response.Data.AddPost.CommentsAllowed)
	assert.Equal(t, expectedPost.CreatedAt, response.Data.AddPost.CreatedAt)
}
