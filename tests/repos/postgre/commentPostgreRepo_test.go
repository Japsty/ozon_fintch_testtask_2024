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
var createdAtTime = time.Now()

var mockPostForComments = model.Post{
	ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Title:           "Test Title",
	Content:         "Test Content",
	UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	Comments:        []*model.Comment{&mockComment, &mockChildComment, &mockChildChildComment},
	CommentsAllowed: true,
	CreatedAt:       "",
}

var mockComment = model.Comment{
	ID:       "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:  "Test Content",
	AuthorID: "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:   "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID: nil,
	Replies: []*model.Comment{
		&mockChildComment,
	},
	CreatedAt: createdAtTime.String(),
}

var childParentID = "14ad7024-7c45-4453-9fac-2dfae1ad2c96"
var mockChildComment = model.Comment{
	ID:        "24ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "Test Content",
	AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  &childParentID,
	Replies:   []*model.Comment{&mockChildChildComment},
	CreatedAt: createdAtTime.String(),
}

var childChildParentID = "24ad7024-7c45-4453-9fac-2dfae1ad2c96"
var mockChildChildComment = model.Comment{
	ID:        "34ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "Test Content",
	AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  &childChildParentID,
	Replies:   nil,
	CreatedAt: createdAtTime.String(),
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
	var parentIDStr string
	if comment.ParentID != nil {
		parentIDStr = *comment.ParentID
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
		ParentID:  parentIDStr,
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

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(querries.CreateComment)).WithArgs(
		mockComment.ID,
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		"",
	).WillReturnResult(
		sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = repo.CreateComment(
		context.Background(),
		mockComment.ID,
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		parentID,
	)
	if err != nil {
		t.Fatalf("CreateComment Error: %s", err)
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

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(querries.GetAllCommentsByPostID)).WithArgs(
		mockComment.PostID,
	).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
				AddRow(
					mockComment.ID,
					mockComment.Content,
					mockComment.AuthorID,
					mockComment.PostID,
					mockComment.ParentID,
					createdAtTime,
				).AddRow(
				mockChildComment.ID,
				mockChildComment.Content,
				mockChildComment.AuthorID,
				mockChildComment.PostID,
				mockChildComment.ParentID,
				createdAtTime,
			).AddRow(
				mockChildChildComment.ID,
				mockChildChildComment.Content,
				mockChildChildComment.AuthorID,
				mockChildChildComment.PostID,
				mockChildChildComment.ParentID,
				createdAtTime,
			),
		)

	mock.ExpectCommit()

	comments, err := repo.GetCommentsByPostID(context.Background(), mockPostForComments.ID)
	if err != nil {
		t.Fatalf("GetCommentsByPostID Error: %s", err)
	}

	expectedComment := mockComment

	for _, comm := range comments {
		if !reflect.DeepEqual(deletePointer(comm), deletePointer(&expectedComment)) {
			t.Errorf(
				"Unexpected comment data. Got %+v, expected %+v",
				deletePointer(comm),
				deletePointer(&expectedComment),
			)
		}
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCommentByPostIDPaginated(t *testing.T) {
	testCases := []struct {
		TestCaseID    int
		Name          string
		Limit         int
		Offset        int
		InputComment  model.Comment
		OutputStruct  []*model.Comment
		ExpectedError error
	}{
		{
			TestCaseID:    1,
			Name:          "Success",
			Limit:         1,
			Offset:        0,
			InputComment:  mockComment,
			OutputStruct:  []*model.Comment{&mockComment, &mockChildComment},
			ExpectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected ", err)
			}
			defer db.Close()

			repo := postgre.NewCommentRepository(db)

			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(querries.GetAllCommentsByPostID)).WithArgs(
				tc.InputComment.PostID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
					AddRow(
						tc.InputComment.ID,
						tc.InputComment.Content,
						tc.InputComment.AuthorID,
						tc.InputComment.PostID,
						tc.InputComment.ParentID,
						createdAtTime,
					).AddRow(
					mockChildComment.ID,
					mockChildComment.Content,
					mockChildComment.AuthorID,
					mockChildComment.PostID,
					mockChildComment.ParentID,
					createdAtTime,
				).AddRow(
					mockChildChildComment.ID,
					mockChildChildComment.Content,
					mockChildChildComment.AuthorID,
					mockChildChildComment.PostID,
					mockChildChildComment.ParentID,
					createdAtTime,
				),
			)
			mock.ExpectCommit()

			comments, err := repo.GetCommentsByPostIDPaginated(context.Background(), mockPostForComments.ID, 10, 0)
			if err != nil {
				t.Fatalf("TestGetCommentByPostIDPaginated Error: %s", err)
			}

			for idx, comm := range comments {
				if !reflect.DeepEqual(deletePointer(comm), deletePointer(tc.OutputStruct[idx])) {
					t.Errorf(
						"Unexpected comment data. Got %+v, expected %+v",
						deletePointer(comm),
						deletePointer(tc.OutputStruct[idx]),
					)
				}
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetCommentByPostIDPaginatedEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected ", err)
	}
	defer db.Close()

	repo := postgre.NewCommentRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querries.GetAllCommentsByPostID)).WithArgs(
		mockComment.PostID,
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
	mock.ExpectCommit()

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
