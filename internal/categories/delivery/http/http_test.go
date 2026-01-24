package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
	"github.com/pandusatrianura/code-with-umam-categories-api/pkg/json_wrapper"
)

type mockService struct {
	GetAllCategoriesFunc func() []entity.Category
	GetCategoryByIDFunc  func(categoryID int64) (entity.Category, error)
	InsertCategoryFunc   func(parameter entity.Category) entity.Category
	UpdateCategoryFunc   func(parameter entity.Category) (entity.Category, error)
	DeleteCategoryFunc   func(categoryID int64) (int64, error)
	APIFunc              func() entity.HealthResponse
}

func (m *mockService) GetAllCategories() []entity.Category {
	return m.GetAllCategoriesFunc()
}
func (m *mockService) GetCategoryByID(categoryID int64) (entity.Category, error) {
	return m.GetCategoryByIDFunc(categoryID)
}
func (m *mockService) InsertCategory(parameter entity.Category) entity.Category {
	return m.InsertCategoryFunc(parameter)
}
func (m *mockService) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	return m.UpdateCategoryFunc(parameter)
}
func (m *mockService) DeleteCategory(categoryID int64) (int64, error) {
	return m.DeleteCategoryFunc(categoryID)
}
func (m *mockService) API() entity.HealthResponse {
	return m.APIFunc()
}

func TestNewCategoriesHandler(t *testing.T) {
	svc := &mockService{}
	handler, err := NewCategoriesHandler(svc)
	if err != nil {
		t.Errorf("NewCategoriesHandler() error = %v, wantErr nil", err)
	}
	if handler.service != svc {
		t.Errorf("NewCategoriesHandler() handler.service = %v, want %v", handler.service, svc)
	}
}

func TestCategoriesHandler_API(t *testing.T) {
	tests := []struct {
		name       string
		mockRes    entity.HealthResponse
		wantStatus int
		wantBody   json_wrapper.APIResponse
	}{
		{
			name: "healthy",
			mockRes: entity.HealthResponse{
				Name:      "Service",
				IsHealthy: true,
			},
			wantStatus: http.StatusOK,
			wantBody: json_wrapper.APIResponse{
				Code:    constants.SuccessCode,
				Message: "Service is healthy",
			},
		},
		{
			name: "unhealthy",
			mockRes: entity.HealthResponse{
				Name:      "Service",
				IsHealthy: false,
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody: json_wrapper.APIResponse{
				Code:    constants.ErrorCode,
				Message: "Service is not healthy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				APIFunc: func() entity.HealthResponse {
					return tt.mockRes
				},
			}
			h := &CategoriesHandler{service: svc}
			req := httptest.NewRequest(http.MethodGet, "/api", nil)
			w := httptest.NewRecorder()

			h.API(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("API() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var gotBody json_wrapper.APIResponse
			json.Unmarshal(w.Body.Bytes(), &gotBody)
			if gotBody.Code != tt.wantBody.Code || gotBody.Message != tt.wantBody.Message {
				t.Errorf("API() body = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestCategoriesHandler_GetAllCategories(t *testing.T) {
	tests := []struct {
		name       string
		mockRes    []entity.Category
		wantStatus int
	}{
		{
			name: "success",
			mockRes: []entity.Category{
				{ID: 1, Name: "Cat 1"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty",
			mockRes:    []entity.Category{},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				GetAllCategoriesFunc: func() []entity.Category {
					return tt.mockRes
				},
			}
			h := &CategoriesHandler{service: svc}
			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			w := httptest.NewRecorder()

			h.GetAllCategories(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetAllCategories() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var gotBody json_wrapper.APIResponse
			json.Unmarshal(w.Body.Bytes(), &gotBody)
			if gotBody.Code != constants.SuccessCode {
				t.Errorf("GetAllCategories() code = %v, want %v", gotBody.Code, constants.SuccessCode)
			}
		})
	}
}

func TestCategoriesHandler_GetCategoryByID(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		mockRes    entity.Category
		mockErr    error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "success",
			path:       "/categories/1",
			mockRes:    entity.Category{ID: 1, Name: "Cat 1"},
			wantStatus: http.StatusOK,
			wantMsg:    "Success get category by id",
		},
		{
			name:       "invalid id",
			path:       "/categories/abc",
			wantStatus: http.StatusBadRequest,
			wantMsg:    constants.ErrInvalidCategoryID,
		},
		{
			name:       "not found",
			path:       "/categories/99",
			mockErr:    errors.New(constants.ErrCategoryNotFound),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    constants.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				GetCategoryByIDFunc: func(id int64) (entity.Category, error) {
					return tt.mockRes, tt.mockErr
				},
			}
			h := &CategoriesHandler{service: svc}
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			h.GetCategoryByID(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetCategoryByID() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var gotBody json_wrapper.APIResponse
			json.Unmarshal(w.Body.Bytes(), &gotBody)
			if gotBody.Message != tt.wantMsg {
				t.Errorf("GetCategoryByID() message = %v, want %v", gotBody.Message, tt.wantMsg)
			}
		})
	}
}

func TestCategoriesHandler_InsertCategory(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		mockRes    entity.Category
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "success",
			body:       entity.Category{Name: "New Cat"},
			mockRes:    entity.Category{ID: 1, Name: "New Cat"},
			wantStatus: http.StatusCreated,
			wantMsg:    "Success insert new category",
		},
		{
			name:       "invalid json",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
			wantMsg:    constants.ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				InsertCategoryFunc: func(c entity.Category) entity.Category {
					return tt.mockRes
				},
			}
			h := &CategoriesHandler{service: svc}

			var bodyReader *bytes.Reader
			if s, ok := tt.body.(string); ok {
				bodyReader = bytes.NewReader([]byte(s))
			} else {
				b, _ := json.Marshal(tt.body)
				bodyReader = bytes.NewReader(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/categories", bodyReader)
			w := httptest.NewRecorder()

			h.InsertCategory(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("InsertCategory() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var gotBody json_wrapper.APIResponse
			json.Unmarshal(w.Body.Bytes(), &gotBody)
			if gotBody.Message != tt.wantMsg {
				t.Errorf("InsertCategory() message = %v, want %v", gotBody.Message, tt.wantMsg)
			}
		})
	}
}

func TestCategoriesHandler_UpdateCategory(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		body       interface{}
		mockRes    entity.Category
		mockErr    error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "success",
			path:       "/categories/1",
			body:       entity.Category{Name: "Updated Cat"},
			mockRes:    entity.Category{ID: 1, Name: "Updated Cat"},
			wantStatus: http.StatusOK,
			wantMsg:    "Success update existing category",
		},
		{
			name:       "invalid json",
			path:       "/categories/1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
			wantMsg:    constants.ErrInvalidRequest,
		},
		{
			name:       "service error",
			path:       "/categories/1",
			body:       entity.Category{Name: "Updated Cat"},
			mockErr:    errors.New("some error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				UpdateCategoryFunc: func(c entity.Category) (entity.Category, error) {
					return tt.mockRes, tt.mockErr
				},
			}
			h := &CategoriesHandler{service: svc}

			var bodyReader *bytes.Reader
			if s, ok := tt.body.(string); ok {
				bodyReader = bytes.NewReader([]byte(s))
			} else if tt.body != nil {
				b, _ := json.Marshal(tt.body)
				bodyReader = bytes.NewReader(b)
			}

			req := httptest.NewRequest(http.MethodPut, tt.path, bodyReader)
			w := httptest.NewRecorder()

			h.UpdateCategory(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdateCategory() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var gotBody json_wrapper.APIResponse
			json.Unmarshal(w.Body.Bytes(), &gotBody)
			if gotBody.Message != tt.wantMsg {
				t.Errorf("UpdateCategory() message = %v, want %v", gotBody.Message, tt.wantMsg)
			}
		})
	}
}

func TestCategoriesHandler_DeleteCategory(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		mockRes    int64
		mockErr    error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "success",
			path:       "/categories/1",
			mockRes:    1,
			wantStatus: http.StatusOK,
			wantMsg:    "Success delete category with id 1",
		},
		{
			name:       "invalid id",
			path:       "/categories/abc",
			wantStatus: http.StatusBadRequest,
			wantMsg:    constants.ErrInvalidCategoryID,
		},
		{
			name:       "service error",
			path:       "/categories/1",
			mockErr:    errors.New("some error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				DeleteCategoryFunc: func(id int64) (int64, error) {
					return tt.mockRes, tt.mockErr
				},
			}
			h := &CategoriesHandler{service: svc}
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			w := httptest.NewRecorder()

			h.DeleteCategory(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteCategory() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var gotBody json_wrapper.APIResponse
			json.Unmarshal(w.Body.Bytes(), &gotBody)
			if gotBody.Message != tt.wantMsg {
				t.Errorf("DeleteCategory() message = %v, want %v", gotBody.Message, tt.wantMsg)
			}
		})
	}
}
