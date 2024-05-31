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

var parentID = ""
var mockComment = model.Comment{
	ID:        "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "Test Content",
	AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  &parentID,
	Replies:   nil,
	CreatedAt: "",
}

type CommentWithoutPointers struct {
	ID        string
	Content   string
	AuthorID  string
	PostID    string
	ParentID  string
	Replies   []CommentWithoutPointers
	CreatedAt string
}

func deletePointer(comment *model.Comment) CommentWithoutPointers {
	var parentID string
	if comment.ParentID != nil {
		parentID = *comment.ParentID
	}

	replies := []CommentWithoutPointers{}
	for _, reply := range comment.Replies {
		replies = append(replies, deletePointer(reply))
	}

	return CommentWithoutPointers{
		ID:        comment.ID,
		Content:   comment.Content,
		AuthorID:  comment.AuthorID,
		PostID:    comment.PostID,
		ParentID:  parentID,
		Replies:   replies,
		CreatedAt: comment.CreatedAt,
	}
}

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewCommentRepository(db)

	createdAt := time.Now()
	mock.ExpectExec(regexp.QuoteMeta(querries.CreateComment)).WithArgs(
		mockComment.ID,
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		mockComment.ParentID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByPostID)).WithArgs(
		mockComment.PostID,
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
					mockComment.ID,
					mockComment.Content,
					mockComment.AuthorID,
					mockComment.PostID,
					mockComment.ParentID,
					createdAt,
				))

	comment, err := repo.CreateComment(context.Background(), mockComment.ID, mockComment.Content, mockComment.AuthorID, mockPost.ID, parentID)
	if err != nil {
		t.Fatalf("CreateComment Error: %s", err)
	}

	expectedComment := model.Comment{
		ID:        "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Content:   "Test Content",
		AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		ParentID:  &parentID,
		Replies:   nil,
		CreatedAt: createdAt.String(),
	}

	if !reflect.DeepEqual(*comment[0], expectedComment) {
		t.Errorf("Unexpected comment data. Got %+v, expected %+v", comment, expectedComment)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCommentByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewCommentRepository(db)

	var childComment = mockComment
	var childChildComment = mockComment

	createdAt := time.Now()

	childCommParentID := "14ad7024-7c45-4453-9fac-2dfae1ad2c96"

	childComment.ID = "24ad7024-7c45-4453-9fac-2dfae1ad2c96"
	childComment.ParentID = &childCommParentID
	childComment.Replies = []*model.Comment{&childChildComment}
	childComment.CreatedAt = createdAt.String()

	childChildCommParentID := "24ad7024-7c45-4453-9fac-2dfae1ad2c96"

	childChildComment.ID = "34ad7024-7c45-4453-9fac-2dfae1ad2c96"
	childChildComment.ParentID = &childChildCommParentID
	childChildComment.CreatedAt = createdAt.String()

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByPostID)).WithArgs(
		mockComment.PostID,
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
					mockComment.ID,
					mockComment.Content,
					mockComment.AuthorID,
					mockComment.PostID,
					mockComment.ParentID,
					createdAt,
				))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByParentID)).WithArgs(
		mockComment.ID,
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
					childComment.ID,
					childComment.Content,
					childComment.AuthorID,
					childComment.PostID,
					childComment.ParentID,
					createdAt,
				))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByParentID)).WithArgs(
		childComment.ID,
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
					childChildComment.ID,
					childChildComment.Content,
					childChildComment.AuthorID,
					childChildComment.PostID,
					childChildComment.ParentID,
					createdAt,
				))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByParentID)).WithArgs(
		childChildComment.ID,
	).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"content",
			"author_id",
			"post_id",
			"parent_id",
			"created_at",
		}))

	comments, err := repo.GetCommentsByPostID(context.Background(), mockPost.ID)
	if err != nil {
		t.Fatalf("GetCommentsByPostID Error: %s", err)
	}

	expectedComment := &model.Comment{
		ID:        "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Content:   "Test Content",
		AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		ParentID:  &parentID,
		Replies:   []*model.Comment{&childComment},
		CreatedAt: createdAt.String(),
	}

	if !reflect.DeepEqual(deletePointer(comments[0]), deletePointer(expectedComment)) {
		t.Errorf("Unexpected comment data. Got %+v, expected %+v", deletePointer(comments[0]), deletePointer(expectedComment))
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCommentByPostIDPaginated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewCommentRepository(db)

	var childComment = mockComment
	var childChildComment = mockComment

	createdAt := time.Now()

	childCommParentID := "14ad7024-7c45-4453-9fac-2dfae1ad2c96"

	childComment.ID = "24ad7024-7c45-4453-9fac-2dfae1ad2c96"
	childComment.ParentID = &childCommParentID
	childComment.Replies = []*model.Comment{&childChildComment}
	childComment.CreatedAt = createdAt.String()

	childChildCommParentID := "24ad7024-7c45-4453-9fac-2dfae1ad2c96"

	childChildComment.ID = "34ad7024-7c45-4453-9fac-2dfae1ad2c96"
	childChildComment.ParentID = &childChildCommParentID
	childChildComment.CreatedAt = createdAt.String()

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByPostIDPaginated)).WithArgs(
		mockComment.PostID, 1, 0,
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
					mockComment.ID,
					mockComment.Content,
					mockComment.AuthorID,
					mockComment.PostID,
					mockComment.ParentID,
					createdAt,
				))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByParentID)).WithArgs(
		mockComment.ID,
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
					childComment.ID,
					childComment.Content,
					childComment.AuthorID,
					childComment.PostID,
					childComment.ParentID,
					createdAt,
				))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByParentID)).WithArgs(
		childComment.ID,
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
					childChildComment.ID,
					childChildComment.Content,
					childChildComment.AuthorID,
					childChildComment.PostID,
					childChildComment.ParentID,
					createdAt,
				))

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByParentID)).WithArgs(
		childChildComment.ID,
	).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"content",
			"author_id",
			"post_id",
			"parent_id",
			"created_at",
		}))

	comments, err := repo.GetCommentsByPostIDPaginated(context.Background(), mockPost.ID, 1, 0)
	if err != nil {
		t.Fatalf("TestGetCommentByPostIDPaginated Error: %s", err)
	}

	expectedComment := &model.Comment{
		ID:        "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
		Content:   "Test Content",
		AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
		PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
		ParentID:  &parentID,
		Replies:   []*model.Comment{&childComment},
		CreatedAt: createdAt.String(),
	}

	if !reflect.DeepEqual(deletePointer(comments[0]), deletePointer(expectedComment)) {
		t.Errorf("Unexpected comment data. Got %+v, expected %+v", deletePointer(comments[0]), deletePointer(expectedComment))
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCommentByPostIDPaginatedEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewCommentRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetCommentsByPostIDPaginated)).WithArgs(
		mockComment.PostID, 0, 0,
	).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"content",
			"author_id",
			"post_id",
			"parent_id",
			"created_at",
		}),
	)

	comments, err := repo.GetCommentsByPostIDPaginated(context.Background(), mockPost.ID, 0, 0)
	if err != nil {
		t.Fatalf("GetCommentsByPostIDPaginated Error: %s", err)
	}

	if len(comments) != 0 {
		t.Errorf("Expected empty comment, but got %+v", comments)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
