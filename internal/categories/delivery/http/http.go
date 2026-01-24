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

// HealthCheck godoc
// @Summary Get health status of categories API
// @Description Memeriksa status kesehatan API kategori
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/health [get]
func (d *CategoriesHandler) API(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse
	svcHealthCheckResult := d.service.API()

	if svcHealthCheckResult.IsHealthy {
		result.Code = constants.SuccessCode
		result.Message = fmt.Sprintf("%s is healthy", svcHealthCheckResult.Name)
		json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	} else {
		result.Code = constants.ErrorCode
		result.Message = fmt.Sprintf("%s is not healthy", svcHealthCheckResult.Name)
		json_wrapper.WriteJSONResponse(w, http.StatusServiceUnavailable, result)
	}

	return
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Mengambil semua data kategori
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/ [get]
func (d *CategoriesHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var result json_wrapper.APIResponse

	res := d.service.GetAllCategories()

	result.Code = constants.SuccessCode
	result.Message = "Success get all categories"
	result.Data = res

	json_wrapper.WriteJSONResponse(w, http.StatusOK, result)
	return
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Mengambil kategori berdasarkan ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories/{id} [get]
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

// InsertCategory godoc
// @Summary Create a new category
// @Description Membuat kategori baru
// @Tags categories
// @Accept json
// @Produce json
// @Param category body entity.Category true "Category Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/categories [post]
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

// UpdateCategory godoc
// @Summary Update category
// @Description Update kategori berdasarkan ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body entity.Category true "Category Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/categories/{id} [put]
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

// DeleteCategory godoc
// @Summary Delete category
// @Description Menghapus kategori berdasarkan ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/{id} [delete]
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
