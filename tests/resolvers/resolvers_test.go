package resolvers

import (
	"Ozon_testtask/internal/services"
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var mockPost = models.Post{
	Score: 1,
	Views: 400,
	Type:  "text",
	Title: "MockPost",
	Author: models.Author{
		Username: "testuser",
		UserID:   "127b3a39-e061-434b-bcbd-d2362abeb4bc",
	},
	Category: "music",
	Text:     "MockPost",
	Votes: []models.PostVote{
		{
			UserID: "127b3a39-e061-434b-bcbd-d2362abeb4bc",
			Vote:   1,
		},
	},
	Comments: []models.Comment{
		{
			Created: time.Now(),
			Author: models.Author{
				Username: "testuser",
				UserID:   "127b3a39-e061-434b-bcbd-d2362abeb4bc",
			},
			Body:   "MockComment",
			ID:     "782a81b7-fbcb-4959-a818-3a86cfaa5686",
			PostID: "1e8250d4-9500-4f91-b1eb-de666774daa6",
		},
	},
	Created:          time.Now(),
	UpvotePercentage: 100,
	URL:              "",
	ID:               "1e8250d4-9500-4f91-b1eb-de666774daa6",
}

var mockSession = &models.Session{
	ID:       1,
	Username: "testuser",
	UserID:   "127b3a39-e061-434b-bcbd-d2362abeb4bc",
	Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQ4Mzg1MDAsInVzZXIiOnsiaWQiOiIxMmVmYWExZi04MjgyLTQzODctOWJkYi03OWU2NDBjYzAxNGMiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.IHd9A9IK9OR5RDWyyJMMZVoyH-ji3tf5yydTjsArHvc",
	Exp:      20000000000,
}

type mockRepoResp struct {
	mockPost  interface{}
	mockError error
}
type mockSessionResp struct {
	mockSession *models.Session
	mockError   error
}
type mockRequest struct {
	mockRequestMethod        string
	mockRequestURL           string
	mockRequestBody          *strings.Reader
	mockIncorrectRequestBody *bytes.Buffer
}

type testCase struct {
	id              int
	name            string
	mockRequest     mockRequest
	mockRepoResp    mockRepoResp
	mockSessionResp mockSessionResp
	callRepo        bool
	breakWrite      bool
	expectedStatus  int
}

func TestPostHandlerIndex(t *testing.T) {
	template1 := template.Must(template.ParseGlob("../../static/html/*"))
	template2 := template.Must(template.ParseGlob("../mockStatic/html/*"))
	test := []struct {
		id             int
		name           string
		template       *template.Template
		expectedStatus int
	}{
		{
			id:             1,
			name:           "Success",
			template:       template1,
			expectedStatus: http.StatusOK,
		},
		{
			id:             2,
			name:           "Error",
			template:       template2,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	mockPostRepo := new(repomocks.MockPostRepo)
	mockUserRepo := new(repomocks.MockUserRepo)
	mockCommentRepo := new(repomocks.MockCommentRepo)

	mockPostService := services.NewPostService(mockUserRepo, mockPostRepo, mockCommentRepo)
	mockSessionService := new(servicemocks.SessionServiceMock)

	zapLogger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}
	logger := zapLogger.Sugar()

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			postHandler := handlers.NewPostHandler(mockSessionService, mockPostService, logger, tc.template)

			req, err := http.NewRequest("GET", "/", strings.NewReader(""))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			postHandler.Index(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}

func TestPostHandlerGetAllPosts(t *testing.T) {
	mockPost2 := mockPost
	mockPost2.Score = 2

	testCases := []testCase{
		{
			id:   1,
			name: "Success",
			mockRequest: mockRequest{
				mockRequestMethod: "GET",
				mockRequestURL:    "/api/posts/",
				mockRequestBody:   strings.NewReader(""),
			},
			mockRepoResp: mockRepoResp{
				mockPost:  []models.Post{mockPost, mockPost2},
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Error",
			mockRequest: mockRequest{
				mockRequestMethod: "GET",
				mockRequestURL:    "/api/posts/",
				mockRequestBody:   strings.NewReader(""),
			},
			mockRepoResp: mockRepoResp{
				mockPost:  []models.Post{},
				mockError: errors.New("Error"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   3,
			name: "Encode Error",
			mockRequest: mockRequest{
				mockRequestMethod: "GET",
				mockRequestURL:    "/api/posts/",
				mockRequestBody:   strings.NewReader(""),
			},
			mockRepoResp: mockRepoResp{
				mockPost:  []models.Post{},
				mockError: nil,
			},
			callRepo:       true,
			breakWrite:     true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostRepo := new(repomocks.MockPostRepo)
			mockUserRepo := new(repomocks.MockUserRepo)
			mockCommentRepo := new(repomocks.MockCommentRepo)

			mockPostService := services.NewPostService(mockUserRepo, mockPostRepo, mockCommentRepo)
			mockSessionService := new(servicemocks.SessionServiceMock)

			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			logger := zapLogger.Sugar()

			templates := template.Must(template.ParseGlob("../../static/html/*"))

			mockPostRepo.On("GetAllPosts", mock.AnythingOfType("*context.timerCtx")).Return(tc.mockRepoResp.mockPost, tc.mockRepoResp.mockError)
			postHandler := handlers.NewPostHandler(mockSessionService, mockPostService, logger, templates)

			req, err := http.NewRequest(tc.mockRequest.mockRequestMethod, tc.mockRequest.mockRequestURL, tc.mockRequest.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}
			rr := httptest.NewRecorder()

			if tc.breakWrite {
				postHandler.GetAllPosts(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				postHandler.GetAllPosts(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockPostRepo.AssertCalled(t, "GetAllPosts", mock.Anything)
			}
		})
	}
}
