package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
	"github.com/pandusatrianura/code-with-umam-categories-api/pkg/json_wrapper"
)

type MockCategoriesService struct{}

func (m *MockCategoriesService) GetAllCategories() []entity.Category {
	return []entity.Category{
		{ID: 1, Name: "Food", Description: "Food category"},
		{ID: 2, Name: "Beverage", Description: "Beverage category"},
	}
}

func (m *MockCategoriesService) GetCategoryByID(categoryID int64) (entity.Category, error) {
	if categoryID == 1 {
		return entity.Category{ID: 1, Name: "Food", Description: "Food category"}, nil
	}
	return entity.Category{}, errors.New(constants.ErrInvalidCategoryID)
}

func (m *MockCategoriesService) InsertCategory(parameter entity.Category) entity.Category {
	return entity.Category{ID: 3, Name: parameter.Name, Description: parameter.Description}
}

func (m *MockCategoriesService) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	if parameter.ID == 1 {
		return parameter, nil
	}
	return entity.Category{}, errors.New("category not found")
}

func (m *MockCategoriesService) DeleteCategory(categoryID int64) (int64, error) {
	if categoryID == 1 {
		return 1, nil
	}
	return 0, errors.New("category not found")
}

func TestGetAllCategories(t *testing.T) {
	handler := &CategoriesHandler{service: &MockCategoriesService{}}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/categories/", nil)

	handler.GetAllCategories(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var resp json_wrapper.APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}

func TestGetCategoryByID(t *testing.T) {
	tests := []struct {
		path           string
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{path: "/categories/1", expectedStatus: http.StatusOK, expectedCode: constants.SuccessCode, expectedMsg: "Success get category by id"},
		{path: "/categories/x", expectedStatus: http.StatusBadRequest, expectedCode: constants.ErrorCode, expectedMsg: constants.ErrInvalidCategoryID},
		{path: "/categories/99", expectedStatus: http.StatusInternalServerError, expectedCode: constants.ErrorCode, expectedMsg: constants.ErrInvalidCategoryID},
	}

	handler := &CategoriesHandler{service: &MockCategoriesService{}}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)

			handler.GetCategoryByID(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			var resp json_wrapper.APIResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}

			if resp.Code != tc.expectedCode || resp.Message != tc.expectedMsg {
				t.Errorf("Expected code %s and message %s, got code %s and message %s", tc.expectedCode, tc.expectedMsg, resp.Code, resp.Message)
			}
		})
	}
}

func TestInsertCategory(t *testing.T) {
	handler := &CategoriesHandler{service: &MockCategoriesService{}}
	tests := []struct {
		body           string
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{`{"name": "Electronics", "description": "Category for electronic items"}`, http.StatusCreated, constants.SuccessCode, "Success insert new category"},
		{`invalid json`, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidRequest},
	}

	for _, tc := range tests {
		t.Run(tc.body, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/categories/", bytes.NewBufferString(tc.body))

			handler.InsertCategory(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			var resp json_wrapper.APIResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}

			if resp.Code != tc.expectedCode || resp.Message != tc.expectedMsg {
				t.Errorf("Expected code %s and message %s, got code %s and message %s", tc.expectedCode, tc.expectedMsg, resp.Code, resp.Message)
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	handler := &CategoriesHandler{service: &MockCategoriesService{}}
	tests := []struct {
		id             string
		body           string
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{"1", `{"name": "Updated Name", "description": "Updated Description"}`, http.StatusOK, constants.SuccessCode, "Success update existing category"},
		{"x", `{"name": "Updated Name", "description": "Updated Description"}`, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidCategoryID},
		{"99", `{"name": "Updated Name", "description": "Updated Description"}`, http.StatusInternalServerError, constants.ErrorCode, "category not found"},
		{"1", `invalid json`, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidRequest},
	}

	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/categories/"+tc.id, bytes.NewBufferString(tc.body))

			handler.UpdateCategory(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			var resp json_wrapper.APIResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}

			if resp.Code != tc.expectedCode || resp.Message != tc.expectedMsg {
				t.Errorf("Expected code %s and message %s, got code %s and message %s", tc.expectedCode, tc.expectedMsg, resp.Code, resp.Message)
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	handler := &CategoriesHandler{service: &MockCategoriesService{}}
	tests := []struct {
		id             string
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{"1", http.StatusOK, constants.SuccessCode, "Success delete category with id 1"},
		{"x", http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidCategoryID},
		{"99", http.StatusInternalServerError, constants.ErrorCode, "category not found"},
	}

	for _, tc := range tests {
		t.Run(tc.id, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/categories/"+tc.id, nil)

			handler.DeleteCategory(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			var resp json_wrapper.APIResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}

			strText, _ := resp.Message.(string)
			if !strings.Contains(strText, tc.expectedMsg) {
				t.Errorf("Expected message containing %q, got %q", tc.expectedMsg, resp.Message)
			}
		})
	}
}
