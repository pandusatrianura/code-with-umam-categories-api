package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
)

type mockRepository struct {
	getAllCategoriesFunc func() []entity.Category
	getCategoryByIDFunc  func(id int64) entity.Category
	insertCategoryFunc   func(category entity.Category) entity.Category
	updateCategoryFunc   func(category entity.Category) (entity.Category, error)
	deleteCategoryFunc   func(id int64) (int64, error)
}

func (m *mockRepository) GetAllCategories() []entity.Category {
	return m.getAllCategoriesFunc()
}

func (m *mockRepository) GetCategoryByID(categoryID int64) entity.Category {
	return m.getCategoryByIDFunc(categoryID)
}

func (m *mockRepository) InsertCategory(parameter entity.Category) entity.Category {
	return m.insertCategoryFunc(parameter)
}

func (m *mockRepository) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	return m.updateCategoryFunc(parameter)
}

func (m *mockRepository) DeleteCategory(categoryID int64) (int64, error) {
	return m.deleteCategoryFunc(categoryID)
}

func TestNewCategoriesService(t *testing.T) {
	repo := &mockRepository{}
	svc, err := NewCategoriesService(repo)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if svc.repo != repo {
		t.Errorf("expected repo to be set")
	}
}

func TestCategoriesService_API(t *testing.T) {
	svc := &CategoriesService{}
	expected := entity.HealthResponse{
		Name:      "Categories API",
		IsHealthy: true,
	}
	got := svc.API()
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestCategoriesService_GetAllCategories(t *testing.T) {
	tests := []struct {
		name     string
		mockData []entity.Category
		expected []entity.Category
	}{
		{
			name:     "Success",
			mockData: []entity.Category{{ID: 1, Name: "Cat 1"}},
			expected: []entity.Category{{ID: 1, Name: "Cat 1"}},
		},
		{
			name:     "Empty",
			mockData: []entity.Category{},
			expected: []entity.Category{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				getAllCategoriesFunc: func() []entity.Category {
					return tt.mockData
				},
			}
			svc := &CategoriesService{repo: repo}
			got := svc.GetAllCategories()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCategoriesService_GetCategoryByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		mockData  entity.Category
		expected  entity.Category
		expectErr bool
		errMsg    string
	}{
		{
			name:     "Found",
			id:       1,
			mockData: entity.Category{ID: 1, Name: "Cat 1"},
			expected: entity.Category{ID: 1, Name: "Cat 1"},
		},
		{
			name:      "Not Found",
			id:        99,
			mockData:  entity.Category{},
			expectErr: true,
			errMsg:    constants.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				getCategoryByIDFunc: func(id int64) entity.Category {
					return tt.mockData
				},
			}
			svc := &CategoriesService{repo: repo}
			got, err := svc.GetCategoryByID(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("expectErr %v, got error %v", tt.expectErr, err)
			}
			if tt.expectErr && err.Error() != tt.errMsg {
				t.Errorf("expected error message %v, got %v", tt.errMsg, err.Error())
			}
			if !tt.expectErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCategoriesService_InsertCategory(t *testing.T) {
	input := entity.Category{Name: "New"}
	expected := entity.Category{ID: 1, Name: "New"}

	repo := &mockRepository{
		insertCategoryFunc: func(c entity.Category) entity.Category {
			return expected
		},
	}
	svc := &CategoriesService{repo: repo}
	got := svc.InsertCategory(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestCategoriesService_UpdateCategory(t *testing.T) {
	tests := []struct {
		name      string
		input     entity.Category
		mockResp  entity.Category
		mockErr   error
		expected  entity.Category
		expectErr bool
	}{
		{
			name:     "Success",
			input:    entity.Category{ID: 1, Name: "Updated"},
			mockResp: entity.Category{ID: 1, Name: "Updated"},
			expected: entity.Category{ID: 1, Name: "Updated"},
		},
		{
			name:      "Error",
			input:     entity.Category{ID: 99},
			mockErr:   errors.New("update error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				updateCategoryFunc: func(c entity.Category) (entity.Category, error) {
					return tt.mockResp, tt.mockErr
				},
			}
			svc := &CategoriesService{repo: repo}
			got, err := svc.UpdateCategory(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("expectErr %v, got %v", tt.expectErr, err)
			}
			if !tt.expectErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCategoriesService_DeleteCategory(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		mockRows  int64
		mockErr   error
		expected  int64
		expectErr bool
	}{
		{
			name:     "Success",
			id:       1,
			mockRows: 1,
			expected: 1,
		},
		{
			name:      "Error",
			id:        99,
			mockErr:   errors.New("delete error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				deleteCategoryFunc: func(id int64) (int64, error) {
					return tt.mockRows, tt.mockErr
				},
			}
			svc := &CategoriesService{repo: repo}
			got, err := svc.DeleteCategory(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("expectErr %v, got %v", tt.expectErr, err)
			}
			if !tt.expectErr && got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
