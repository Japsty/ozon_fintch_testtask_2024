package repos

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
	"time"
)

var mockPost = model.Post{
	ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Title:           "Test Title",
	Content:         "Test Content",
	UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	Comments:        nil,
	CommentsAllowed: true,
	CreatedAt:       "",
}

func TestGetAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewPostRepository(db)

	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetAllPosts)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
			AddRow(mockPost.ID, mockPost.Title, mockPost.Content, mockPost.UserID, true, createdAt))

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
		CommentsAllowed: true,
		CreatedAt:       createdAt.String(),
	}

	if !reflect.DeepEqual(*posts[0], expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", *posts[0], expectedPost)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewPostRepository(db)

	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(querries.CreatePost)).WithArgs(
		mockPost.ID,
		mockPost.Title,
		mockPost.Content,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"title",
					"content",
					"user_id",
					"comments_allowed",
					"created_at",
				}).
				AddRow(
					mockPost.ID,
					mockPost.Title,
					mockPost.Content,
					mockPost.UserID,
					true,
					createdAt,
				))

	post, err := repo.CreatePost(context.Background(), mockPost.ID, mockPost.Title, mockPost.Content, mockPost.UserID, mockPost.CommentsAllowed)
	if err != nil {
		t.Fatalf("CreatePost Error: %s", err)
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		CommentsAllowed: true,
		CreatedAt:       createdAt.String(),
	}

	if !reflect.DeepEqual(post, expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", post, expectedPost)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPostByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewPostRepository(db)

	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetPostByID)).
		WithArgs(mockPost.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
			AddRow(mockPost.ID, mockPost.Title, mockPost.Content, mockPost.UserID, true, createdAt))

	posts, err := repo.GetPostByPostID(context.Background(), mockPost.ID)
	if err != nil {
		t.Fatalf("GetPostByPostI Error: %s", err)
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		CommentsAllowed: true,
		CreatedAt:       createdAt.String(),
	}

	if !reflect.DeepEqual(posts, expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", posts, expectedPost)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdatePostCommentsStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewPostRepository(db)

	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetPostByID)).
		WithArgs(mockPost.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
			AddRow(mockPost.ID, mockPost.Title, mockPost.Content, mockPost.UserID, true, createdAt))

	mock.ExpectQuery(regexp.QuoteMeta(querries.UpdatePostCommentsStatus)).
		WithArgs(mockPost.ID, mockPost.CommentsAllowed).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
			AddRow(mockPost.ID, mockPost.Title, mockPost.Content, mockPost.UserID, true, createdAt))

	posts, err := repo.UpdatePostCommentsStatus(context.Background(), mockPost.ID, mockPost.UserID, mockPost.CommentsAllowed)
	if err != nil {
		t.Fatalf("GetPostByPostI Error: %s", err)
	}

	expectedPost := model.Post{
		ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Title:           "Test Title",
		Content:         "Test Content",
		UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		CommentsAllowed: true,
		CreatedAt:       createdAt.String(),
	}

	if !reflect.DeepEqual(posts, expectedPost) {
		t.Errorf("Unexpected post data. Got %+v, expected %+v", posts, expectedPost)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
