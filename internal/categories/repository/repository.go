package repository

import (
	"fmt"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/utils"
)

// ICategoriesRepository defines an abstraction for performing CRUD operations on Category entities.
type ICategoriesRepository interface {
	GetAllCategories() []entity.Category
	GetCategoryByID(categoryID int64) entity.Category
	InsertCategory(parameter entity.Category) entity.Category
	UpdateCategory(parameter entity.Category) (entity.Category, error)
	DeleteCategory(categoryID int64) (int64, error)
}

// categoriesList holds a predefined list of categories represented as entity.Category objects for initial data setup.
var categoriesList = []entity.Category{
	{
		ID:          1,
		Name:        "Elektronik",
		Description: "Kategori Elektronik",
	},
	{
		ID:          2,
		Name:        "Komputer",
		Description: " Kategori Komputer",
	},
	{
		ID:          3,
		Name:        "Handphone",
		Description: "Kategori Handphone",
	},
}

// CategoriesRepository manages CRUD operations for Category entity.
type CategoriesRepository struct{}

// NewCategoriesRepository initializes and returns a new instance of CategoriesRepository or an error if creation fails.
func NewCategoriesRepository() (*CategoriesRepository, error) {
	r := &CategoriesRepository{}
	return r, nil
}

// GetAllCategories retrieves all categories from the repository and returns them as a slice of entity.Category.
func (r *CategoriesRepository) GetAllCategories() []entity.Category {
	categories := categoriesList
	return categories
}

// GetCategoryByID retrieves a category from the list based on the provided category ID. Returns an empty category if not found.
func (r *CategoriesRepository) GetCategoryByID(categoryID int64) entity.Category {
	categories := categoriesList

	for _, category := range categories {
		if category.ID == categoryID {
			return category
		}
	}

	return entity.Category{}
}

// InsertCategory adds a new category to the categories list. It assigns a new ID if the given ID is 0 and returns the category.
func (r *CategoriesRepository) InsertCategory(parameter entity.Category) entity.Category {
	var cat entity.Category
	if parameter.ID == 0 {
		cat.ID = int64(utils.GetMaxID(categoriesList)) + 1
	} else {
		cat.ID = parameter.ID
	}

	cat.Name = parameter.Name
	cat.Description = parameter.Description

	categories := categoriesList
	categories = append(categories, cat)
	categoriesList = categories

	return cat
}

// UpdateCategory updates an existing category with new data or returns an error if the category is not found.
func (r *CategoriesRepository) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	var cat entity.Category
	cat.ID = parameter.ID
	cat.Name = parameter.Name
	cat.Description = parameter.Description

	categories := categoriesList
	for i, category := range categories {
		if category.ID == parameter.ID {
			categories[i] = cat
			return cat, nil
		}
	}

	return entity.Category{}, fmt.Errorf(constants.ErrCategoryNotFound)
}

// DeleteCategory removes a category by its ID and returns the ID of the deleted category or an error if not found.
func (r *CategoriesRepository) DeleteCategory(categoryID int64) (int64, error) {
	categories := categoriesList
	for i, category := range categories {
		if category.ID == categoryID {
			categories = append(categories[:i], categories[i+1:]...)
			categoriesList = categories
			return category.ID, nil
		}
	}

	return 0, fmt.Errorf(constants.ErrCategoryNotFound)
}
