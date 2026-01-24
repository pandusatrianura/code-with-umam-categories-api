package repository

import (
	"reflect"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
)

func withCategories(t *testing.T, categories []entity.Category, fn func()) {
	t.Helper()
	original := append([]entity.Category(nil), categoriesList...)
	categoriesList = append([]entity.Category(nil), categories...)
	defer func() {
		categoriesList = original
	}()
	fn()
}

func TestNewCategoriesRepository(t *testing.T) {
	repo, err := NewCategoriesRepository()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Fatalf("expected repository instance")
	}
}

func TestCategoriesRepository_GetAllCategories(t *testing.T) {
	tests := []struct {
		name       string
		categories []entity.Category
	}{
		{
			name:       "empty",
			categories: nil,
		},
		{
			name: "some",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
				{ID: 2, Name: "B", Description: "D2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCategories(t, tt.categories, func() {
				repo := &CategoriesRepository{}
				got := repo.GetAllCategories()
				if !reflect.DeepEqual(got, tt.categories) {
					t.Fatalf("expected %v, got %v", tt.categories, got)
				}
			})
		})
	}
}

func TestCategoriesRepository_GetCategoryByID(t *testing.T) {
	tests := []struct {
		name       string
		categories []entity.Category
		id         int64
		want       entity.Category
	}{
		{
			name: "found",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
				{ID: 2, Name: "B", Description: "D2"},
			},
			id:   2,
			want: entity.Category{ID: 2, Name: "B", Description: "D2"},
		},
		{
			name:       "missing",
			categories: []entity.Category{{ID: 1, Name: "A", Description: "D1"}},
			id:         3,
			want:       entity.Category{},
		},
		{
			name:       "empty",
			categories: nil,
			id:         1,
			want:       entity.Category{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCategories(t, tt.categories, func() {
				repo := &CategoriesRepository{}
				got := repo.GetCategoryByID(tt.id)
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
			})
		})
	}
}

func TestCategoriesRepository_InsertCategory(t *testing.T) {
	tests := []struct {
		name       string
		categories []entity.Category
		input      entity.Category
		want       entity.Category
		wantList   []entity.Category
	}{
		{
			name:       "autoEmpty",
			categories: nil,
			input:      entity.Category{Name: "A", Description: "D1"},
			want:       entity.Category{ID: 1, Name: "A", Description: "D1"},
			wantList:   []entity.Category{{ID: 1, Name: "A", Description: "D1"}},
		},
		{
			name: "autoExisting",
			categories: []entity.Category{
				{ID: 2, Name: "A", Description: "D1"},
				{ID: 5, Name: "B", Description: "D2"},
			},
			input:    entity.Category{Name: "C", Description: "D3"},
			want:     entity.Category{ID: 6, Name: "C", Description: "D3"},
			wantList: []entity.Category{{ID: 2, Name: "A", Description: "D1"}, {ID: 5, Name: "B", Description: "D2"}, {ID: 6, Name: "C", Description: "D3"}},
		},
		{
			name: "givenID",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
			},
			input:    entity.Category{ID: 10, Name: "B", Description: "D2"},
			want:     entity.Category{ID: 10, Name: "B", Description: "D2"},
			wantList: []entity.Category{{ID: 1, Name: "A", Description: "D1"}, {ID: 10, Name: "B", Description: "D2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCategories(t, tt.categories, func() {
				repo := &CategoriesRepository{}
				got := repo.InsertCategory(tt.input)
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
				if !reflect.DeepEqual(categoriesList, tt.wantList) {
					t.Fatalf("expected list %v, got %v", tt.wantList, categoriesList)
				}
			})
		})
	}
}

func TestCategoriesRepository_UpdateCategory(t *testing.T) {
	tests := []struct {
		name       string
		categories []entity.Category
		input      entity.Category
		want       entity.Category
		wantList   []entity.Category
		wantErr    string
	}{
		{
			name: "found",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
				{ID: 2, Name: "B", Description: "D2"},
			},
			input:    entity.Category{ID: 2, Name: "BB", Description: "DD"},
			want:     entity.Category{ID: 2, Name: "BB", Description: "DD"},
			wantList: []entity.Category{{ID: 1, Name: "A", Description: "D1"}, {ID: 2, Name: "BB", Description: "DD"}},
		},
		{
			name: "missing",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
			},
			input:    entity.Category{ID: 2, Name: "B", Description: "D2"},
			want:     entity.Category{},
			wantList: []entity.Category{{ID: 1, Name: "A", Description: "D1"}},
			wantErr:  constants.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCategories(t, tt.categories, func() {
				repo := &CategoriesRepository{}
				got, err := repo.UpdateCategory(tt.input)
				if tt.wantErr == "" && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr != "" {
					if err == nil {
						t.Fatalf("expected error")
					}
					if err.Error() != tt.wantErr {
						t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
					}
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
				if !reflect.DeepEqual(categoriesList, tt.wantList) {
					t.Fatalf("expected list %v, got %v", tt.wantList, categoriesList)
				}
			})
		})
	}
}

func TestCategoriesRepository_DeleteCategory(t *testing.T) {
	tests := []struct {
		name       string
		categories []entity.Category
		id         int64
		wantID     int64
		wantList   []entity.Category
		wantErr    string
	}{
		{
			name: "found",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
				{ID: 2, Name: "B", Description: "D2"},
			},
			id:       1,
			wantID:   1,
			wantList: []entity.Category{{ID: 2, Name: "B", Description: "D2"}},
		},
		{
			name: "missing",
			categories: []entity.Category{
				{ID: 1, Name: "A", Description: "D1"},
			},
			id:       2,
			wantID:   0,
			wantList: []entity.Category{{ID: 1, Name: "A", Description: "D1"}},
			wantErr:  constants.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCategories(t, tt.categories, func() {
				repo := &CategoriesRepository{}
				gotID, err := repo.DeleteCategory(tt.id)
				if tt.wantErr == "" && err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr != "" {
					if err == nil {
						t.Fatalf("expected error")
					}
					if err.Error() != tt.wantErr {
						t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
					}
				}
				if gotID != tt.wantID {
					t.Fatalf("expected id %d, got %d", tt.wantID, gotID)
				}
				if !reflect.DeepEqual(categoriesList, tt.wantList) {
					t.Fatalf("expected list %v, got %v", tt.wantList, categoriesList)
				}
			})
		})
	}
}
