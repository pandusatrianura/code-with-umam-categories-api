package service

import (
	"errors"
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAllCategories() []entity.Category {
	args := m.Called()
	return args.Get(0).([]entity.Category)
}

func (m *MockRepository) GetCategoryByID(categoryID int64) entity.Category {
	args := m.Called(categoryID)
	return args.Get(0).(entity.Category)
}

func (m *MockRepository) InsertCategory(parameter entity.Category) entity.Category {
	args := m.Called(parameter)
	return args.Get(0).(entity.Category)
}

func (m *MockRepository) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	args := m.Called(parameter)
	return args.Get(0).(entity.Category), args.Error(1)
}

func (m *MockRepository) DeleteCategory(categoryID int64) (int64, error) {
	args := m.Called(categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetAllCategories(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &CategoriesService{repo: mockRepo}

	mockCategories := []entity.Category{
		{ID: 1, Name: "Category1", Description: "Description1"},
		{ID: 2, Name: "Category2", Description: "Description2"},
	}
	mockRepo.On("GetAllCategories").Return(mockCategories)

	actual := service.GetAllCategories()
	assert.Equal(t, mockCategories, actual)
	mockRepo.AssertExpectations(t)
}

func TestGetCategoryByID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &CategoriesService{repo: mockRepo}

	t.Run("Valid ID", func(t *testing.T) {
		expectedCategory := entity.Category{ID: 1, Name: "Category1", Description: "Description1"}
		mockRepo.On("GetCategoryByID", int64(1)).Return(expectedCategory)

		actual, err := service.GetCategoryByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, actual)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		mockRepo.On("GetCategoryByID", int64(2)).Return(entity.Category{})
		actual, err := service.GetCategoryByID(2)
		assert.Error(t, err)
		assert.EqualError(t, err, constants.ErrCategoryNotFound)
		assert.Equal(t, entity.Category{}, actual)
		mockRepo.AssertExpectations(t)
	})
}

func TestInsertCategory(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &CategoriesService{repo: mockRepo}

	newCategory := entity.Category{Name: "NewCategory", Description: "NewDescription"}
	expectedCategory := entity.Category{ID: 1, Name: "NewCategory", Description: "NewDescription"}
	mockRepo.On("InsertCategory", newCategory).Return(expectedCategory)

	actual := service.InsertCategory(newCategory)
	assert.Equal(t, expectedCategory, actual)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCategory(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &CategoriesService{repo: mockRepo}

	t.Run("Valid Update", func(t *testing.T) {
		updateCategory := entity.Category{ID: 1, Name: "UpdatedCategory", Description: "UpdatedDescription"}
		mockRepo.On("UpdateCategory", updateCategory).Return(updateCategory, nil)

		actual, err := service.UpdateCategory(updateCategory)
		assert.NoError(t, err)
		assert.Equal(t, updateCategory, actual)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		updateCategory := entity.Category{ID: 2, Name: "InvalidCategory", Description: "InvalidDescription"}
		mockRepo.On("UpdateCategory", updateCategory).Return(entity.Category{}, errors.New(constants.ErrCategoryNotFound))

		actual, err := service.UpdateCategory(updateCategory)
		assert.Error(t, err)
		assert.EqualError(t, err, constants.ErrCategoryNotFound)
		assert.Equal(t, entity.Category{}, actual)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &CategoriesService{repo: mockRepo}

	t.Run("Valid Delete", func(t *testing.T) {
		mockRepo.On("DeleteCategory", int64(1)).Return(int64(1), nil)

		id, err := service.DeleteCategory(1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		mockRepo.On("DeleteCategory", int64(2)).Return(int64(0), errors.New(constants.ErrCategoryNotFound))

		id, err := service.DeleteCategory(2)
		assert.Error(t, err)
		assert.EqualError(t, err, constants.ErrCategoryNotFound)
		assert.Equal(t, int64(0), id)
		mockRepo.AssertExpectations(t)
	})
}
