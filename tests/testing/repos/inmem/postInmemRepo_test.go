package repos

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/inmem"
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var mockPost = model.Post{
	ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Title:           "Test Title",
	Content:         "Test Content",
	UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	Comments:        []*model.Comment{},
	CommentsAllowed: true,
	CreatedAt:       "",
}

func TestGetAllPosts(t *testing.T) {
	repo := inmem.NewPostInMemoryRepository()

	post, err := repo.CreatePost(
		context.Background(),
		mockPost.ID,
		mockPost.Title,
		mockPost.Content,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	)
	if err != nil {
		t.Fatalf("CreatePost Error: %s", err)
	}

	posts, err := repo.GetAllPosts(context.Background())
	if err != nil {
		t.Fatalf("GetAllPosts Error: %s", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		Comments:        []*model.Comment{},
		CommentsAllowed: true,
		CreatedAt:       post.CreatedAt,
	}

	if !reflect.DeepEqual(post, expectedPost) {
		assert.Equal(t, post.ID, expectedPost.ID)
		t.Errorf("Unexpected post data. Got %+v, expected %+v", *posts[0], expectedPost)
	}
}

func TestCreatePost(t *testing.T) {
	repo := inmem.NewPostInMemoryRepository()

	post, err := repo.CreatePost(
		context.Background(),
		mockPost.ID,
		mockPost.Title,
		mockPost.Content,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	)
	if err != nil {
		t.Fatalf("CreatePost Error: %s", err)
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		Comments:        []*model.Comment{},
		CommentsAllowed: true,
		CreatedAt:       post.CreatedAt,
	}

	if !reflect.DeepEqual(post, expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", post, expectedPost)
	}
}

func TestGetPostByPostID(t *testing.T) {
	repo := inmem.NewPostInMemoryRepository()

	_, err := repo.CreatePost(
		context.Background(),
		mockPost.ID,
		mockPost.Title,
		mockPost.Content,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	)
	if err != nil {
		t.Fatalf("CreatePost Error: %s", err)
	}

	post, err := repo.GetPostByPostID(context.Background(), mockPost.ID)
	if err != nil {
		t.Fatalf("GetPostByPostID Error: %s", err)
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		Comments:        []*model.Comment{},
		CommentsAllowed: true,
		CreatedAt:       post.CreatedAt,
	}

	if !reflect.DeepEqual(post, expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", post, expectedPost)
	}
}

func TestUpdatePostCommentsStatus(t *testing.T) {
	repo := inmem.NewPostInMemoryRepository()

	_, err := repo.CreatePost(
		context.Background(),
		mockPost.ID,
		mockPost.Title,
		mockPost.Content,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	)
	if err != nil {
		t.Fatalf("CreatePost Error: %s", err)
	}

	post, err := repo.UpdatePostStatus(
		context.Background(),
		mockPost.ID,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	)
	if err != nil {
		t.Fatalf("GetPostByPostID Error: %s", err)
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		Comments:        []*model.Comment{},
		CommentsAllowed: true,
		CreatedAt:       post.CreatedAt,
	}

	if !reflect.DeepEqual(post, expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", post, expectedPost)
	}
}
