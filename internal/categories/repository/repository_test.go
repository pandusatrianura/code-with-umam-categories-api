package repository

import (
	"fmt"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
)

func TestCategoriesRepository_GetAllCategories(t *testing.T) {
	repo := &CategoriesRepository{}
	expectedCategories := categoriesList

	categories := repo.GetAllCategories()

	if len(categories) != len(expectedCategories) {
		t.Errorf("expected %d categories, got %d", len(expectedCategories), len(categories))
	}
	for i, category := range categories {
		if category != expectedCategories[i] {
			t.Errorf("expected category %+v, got %+v", expectedCategories[i], category)
		}
	}
}

func TestCategoriesRepository_GetCategoryByID(t *testing.T) {
	repo := &CategoriesRepository{}

	tests := []struct {
		name       string
		categoryID int64
		expected   entity.Category
	}{
		{"existing ID", 1, categoriesList[0]},
		{"non-existing ID", 99, entity.Category{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.GetCategoryByID(tt.categoryID)
			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestCategoriesRepository_InsertCategory(t *testing.T) {
	repo := &CategoriesRepository{}

	tests := []struct {
		name     string
		input    entity.Category
		validate func(entity.Category) error
	}{
		{
			"new category with ID 0",
			entity.Category{Name: "Books", Description: "Category for books"},
			func(cat entity.Category) error {
				if cat.ID <= 0 {
					return fmt.Errorf("expected a valid ID, got %d", cat.ID)
				}
				if cat.Name != "Books" || cat.Description != "Category for books" {
					return fmt.Errorf("unexpected category data: %+v", cat)
				}
				return nil
			},
		},
		{
			"new category with provided ID",
			entity.Category{ID: 10, Name: "Toys", Description: "Category for toys"},
			func(cat entity.Category) error {
				if cat.ID != 10 {
					return fmt.Errorf("expected ID 10, got %d", cat.ID)
				}
				if cat.Name != "Toys" || cat.Description != "Category for toys" {
					return fmt.Errorf("unexpected category data: %+v", cat)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.InsertCategory(tt.input)
			if err := tt.validate(result); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCategoriesRepository_UpdateCategory(t *testing.T) {
	repo := &CategoriesRepository{}

	tests := []struct {
		name        string
		input       entity.Category
		expectedErr error
	}{
		{
			"update existing category",
			entity.Category{ID: 1, Name: "Updated Electronics", Description: "Updated description"},
			nil,
		},
		{
			"update non-existing category",
			entity.Category{ID: 99, Name: "Non-existent", Description: "Description"},
			fmt.Errorf(constants.ErrCategoryNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.UpdateCategory(tt.input)
			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %+v, got %+v", tt.expectedErr, err)
			}
		})
	}
}

func TestCategoriesRepository_DeleteCategory(t *testing.T) {
	repo := &CategoriesRepository{}

	tests := []struct {
		name        string
		categoryID  int64
		expectedID  int64
		expectedErr error
	}{
		{"delete existing category", 1, 1, nil},
		{"delete non-existing category", 99, 0, fmt.Errorf(constants.ErrCategoryNotFound)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.DeleteCategory(tt.categoryID)
			if id != tt.expectedID {
				t.Errorf("expected ID %d, got %d", tt.expectedID, id)
			}
			if (err == nil) != (tt.expectedErr == nil) || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error %+v, got %+v", tt.expectedErr, err)
			}
		})
	}
}
