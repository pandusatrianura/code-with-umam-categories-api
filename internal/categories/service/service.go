package service

import (
	"fmt"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/repository"
)

// ICategoriesService provides methods for managing category entities.
// GetAllCategories retrieves all categories from the storage.
// GetCategoryByID retrieves a category by its unique identifier.
// InsertCategory creates a new category in the storage.
// UpdateCategory updates an existing category's details.
// DeleteCategory removes a category from storage using its ID.
type ICategoriesService interface {
	GetAllCategories() []entity.Category
	GetCategoryByID(categoryID int64) (entity.Category, error)
	InsertCategory(parameter entity.Category) entity.Category
	UpdateCategory(parameter entity.Category) (entity.Category, error)
	DeleteCategory(categoryID int64) (int64, error)
}

// CategoriesService provides methods to manage and manipulate category data using the ICategoriesRepository abstraction.
type CategoriesService struct {
	repo repository.ICategoriesRepository
}

// NewCategoriesService initializes a new CategoriesService instance with the provided ICategoriesRepository implementation.
func NewCategoriesService(repo repository.ICategoriesRepository) (*CategoriesService, error) {
	return &CategoriesService{
		repo: repo,
	}, nil
}

// GetAllCategories retrieves all categories from the repository and returns them as a slice of Category entities.
func (s *CategoriesService) GetAllCategories() []entity.Category {
	return s.repo.GetAllCategories()
}

// GetCategoryByID retrieves a category by its ID from the repository. Returns an error if the category is not found.
func (s *CategoriesService) GetCategoryByID(categoryID int64) (entity.Category, error) {
	cat := s.repo.GetCategoryByID(categoryID)

	if cat.ID == 0 {
		return entity.Category{}, fmt.Errorf(constants.ErrCategoryNotFound)
	}

	return cat, nil
}

// InsertCategory adds a new category to the repository and returns the created category.
func (s *CategoriesService) InsertCategory(parameter entity.Category) entity.Category {
	return s.repo.InsertCategory(parameter)
}

// UpdateCategory updates an existing category in the data source and returns the updated category or an error if any occurs.
func (s *CategoriesService) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	return s.repo.UpdateCategory(parameter)
}

// DeleteCategory removes a category by its ID and returns the number of rows affected or an error if the operation fails.
func (s *CategoriesService) DeleteCategory(categoryID int64) (int64, error) {
	return s.repo.DeleteCategory(categoryID)
}
