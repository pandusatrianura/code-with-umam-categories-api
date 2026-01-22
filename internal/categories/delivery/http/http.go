package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pandusatrianura/code-with-umam-categories-api/constants"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/service"
	"github.com/pandusatrianura/code-with-umam-categories-api/pkg/json_wrapper"
)

// CategoriesHandler serves as an HTTP handler that processes category-related requests with the help of ICategoriesService.
type CategoriesHandler struct {
	service service.ICategoriesService
}

// NewCategoriesHandler initializes and returns a new CategoriesHandler instance with the provided ICategoriesService implementation.
func NewCategoriesHandler(service service.ICategoriesService) (*CategoriesHandler, error) {
	delegate := &CategoriesHandler{
		service: service,
	}

	return delegate, nil
}

// GetAllCategories handles the HTTP request to retrieve all category records and sends a successful JSON response.
func (d *CategoriesHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse

	res := d.service.GetAllCategories()

	result.Code = constants.SuccessCode
	result.Message = "Success get all categories"
	result.Data = res

	json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	return
}

// GetCategoryByID retrieves a specific category by its ID, handles errors, and sends an appropriate HTTP response.
func (d *CategoriesHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = constants.ErrInvalidCategoryID
		json_wrapper.WriteJSONResponse(w, http.StatusBadRequest, result)
		return
	}

	category, err := d.service.GetCategoryByID(int64(id))
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = err.Error()
		json_wrapper.WriteJSONResponse(w, http.StatusInternalServerError, result)
		return
	}

	result.Code = constants.SuccessCode
	result.Message = "Success get category by id"
	result.Data = category
	json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	return
}

// InsertCategory handles the HTTP POST request for creating a new category by parsing the request body into a Category entity.
func (d *CategoriesHandler) InsertCategory(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse

	var categoryNew entity.Category
	err := json_wrapper.ParseJSON(r, &categoryNew)
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = constants.ErrInvalidRequest
		json_wrapper.WriteJSONResponse(w, http.StatusBadRequest, result)
		return
	}

	data := d.service.InsertCategory(categoryNew)
	result.Code = constants.SuccessCode
	result.Message = "Success insert new category"
	result.Data = data
	json_wrapper.WriteJSONResponse(w, http.StatusCreated, result)
	return
}

// UpdateCategory handles HTTP PUT requests to update an existing category by its ID with the provided data.
func (d *CategoriesHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = constants.ErrInvalidCategoryID
		json_wrapper.WriteJSONResponse(w, http.StatusBadRequest, result)
		return
	}

	var categoryExisting entity.Category
	err = json_wrapper.ParseJSON(r, &categoryExisting)
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = constants.ErrInvalidRequest
		json_wrapper.WriteJSONResponse(w, http.StatusBadRequest, result)
		return
	}

	categoryExisting.ID = int64(id)
	res, err := d.service.UpdateCategory(categoryExisting)
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = err.Error()
		json_wrapper.WriteJSONResponse(w, http.StatusInternalServerError, result)
		return
	}

	result.Code = constants.SuccessCode
	result.Message = "Success update existing category"
	result.Data = res
	json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	return
}

// DeleteCategory handles HTTP requests to delete a category by its ID, responding with success or error details.
func (d *CategoriesHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = constants.ErrInvalidCategoryID
		json_wrapper.WriteJSONResponse(w, http.StatusBadRequest, result)
		return
	}

	res, err := d.service.DeleteCategory(int64(id))
	if err != nil {
		result.Code = constants.ErrorCode
		result.Message = err.Error()
		json_wrapper.WriteJSONResponse(w, http.StatusInternalServerError, result)
		return
	}

	result.Code = constants.SuccessCode
	result.Message = fmt.Sprintf("Success delete category with id %d", res)
	json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	return
}
