package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	categoriesHandler "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/delivery/http"
	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
)

type fakeCategoriesService struct {
	apiResp entity.HealthResponse

	getAllResp []entity.Category

	getByIDResp entity.Category
	getByIDErr  error

	insertResp entity.Category

	updateResp entity.Category
	updateErr  error

	deleteResp int64
	deleteErr  error

	apiCalls     int
	getAllCalls  int
	getByIDCalls int
	insertCalls  int
	updateCalls  int
	deleteCalls  int

	lastGetByID int64
	lastInsert  entity.Category
	lastUpdate  entity.Category
	lastDelete  int64
}

func (f *fakeCategoriesService) GetAllCategories() []entity.Category {
	f.getAllCalls++
	return f.getAllResp
}

func (f *fakeCategoriesService) GetCategoryByID(categoryID int64) (entity.Category, error) {
	f.getByIDCalls++
	f.lastGetByID = categoryID
	if f.getByIDErr != nil {
		return entity.Category{}, f.getByIDErr
	}
	return f.getByIDResp, nil
}

func (f *fakeCategoriesService) InsertCategory(parameter entity.Category) entity.Category {
	f.insertCalls++
	f.lastInsert = parameter
	return f.insertResp
}

func (f *fakeCategoriesService) UpdateCategory(parameter entity.Category) (entity.Category, error) {
	f.updateCalls++
	f.lastUpdate = parameter
	if f.updateErr != nil {
		return entity.Category{}, f.updateErr
	}
	return f.updateResp, nil
}

func (f *fakeCategoriesService) DeleteCategory(categoryID int64) (int64, error) {
	f.deleteCalls++
	f.lastDelete = categoryID
	if f.deleteErr != nil {
		return 0, f.deleteErr
	}
	return f.deleteResp, nil
}

func (f *fakeCategoriesService) API() entity.HealthResponse {
	f.apiCalls++
	return f.apiResp
}

func TestNewRouter(t *testing.T) {
	svc := &fakeCategoriesService{}
	handler, err := categoriesHandler.NewCategoriesHandler(svc)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}

	cases := []struct {
		name    string
		handler *categoriesHandler.CategoriesHandler
	}{
		{name: "nil", handler: nil},
		{name: "set", handler: handler},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := NewRouter(tc.handler)
			if router == nil {
				t.Fatal("expected router")
			}
			if router.categories != tc.handler {
				t.Fatalf("expected handler to match")
			}
		})
	}
}

func TestRouter_RegisterRoutes(t *testing.T) {
	type callCounts struct {
		api     int
		getAll  int
		getByID int
		insert  int
		update  int
		del     int
	}

	type expectations struct {
		calls         callCounts
		getByID       *int64
		updateID      *int64
		insertName    *string
		deleteID      *int64
		bodyContains  string
		expectStatus  int
		needsRepoRoot bool
	}

	int64Ptr := func(v int64) *int64 { return &v }
	stringPtr := func(v string) *string { return &v }

	cases := []struct {
		name     string
		method   string
		path     string
		body     string
		setupSvc func(*fakeCategoriesService)
		expect   expectations
	}{
		{
			name:   "health ok",
			method: http.MethodGet,
			path:   "/categories/health",
			setupSvc: func(svc *fakeCategoriesService) {
				svc.apiResp = entity.HealthResponse{Name: "svc", IsHealthy: true}
			},
			expect: expectations{
				expectStatus: http.StatusOK,
				calls:        callCounts{api: 1},
			},
		},
		{
			name:   "health bad",
			method: http.MethodGet,
			path:   "/categories/health",
			setupSvc: func(svc *fakeCategoriesService) {
				svc.apiResp = entity.HealthResponse{Name: "svc", IsHealthy: false}
			},
			expect: expectations{
				expectStatus: http.StatusServiceUnavailable,
				calls:        callCounts{api: 1},
			},
		},
		{
			name:   "get all",
			method: http.MethodGet,
			path:   "/categories",
			setupSvc: func(svc *fakeCategoriesService) {
				svc.getAllResp = []entity.Category{{ID: 1, Name: "A"}}
			},
			expect: expectations{
				expectStatus: http.StatusOK,
				calls:        callCounts{getAll: 1},
			},
		},
		{
			name:   "get by id ok",
			method: http.MethodGet,
			path:   "/categories/7",
			setupSvc: func(svc *fakeCategoriesService) {
				svc.getByIDResp = entity.Category{ID: 7, Name: "A"}
			},
			expect: expectations{
				expectStatus: http.StatusOK,
				calls:        callCounts{getByID: 1},
				getByID:      int64Ptr(7),
			},
		},
		{
			name:   "get by id err",
			method: http.MethodGet,
			path:   "/categories/8",
			setupSvc: func(svc *fakeCategoriesService) {
				svc.getByIDErr = errors.New("boom")
			},
			expect: expectations{
				expectStatus: http.StatusInternalServerError,
				calls:        callCounts{getByID: 1},
				getByID:      int64Ptr(8),
			},
		},
		{
			name:   "get by id bad",
			method: http.MethodGet,
			path:   "/categories/abc",
			expect: expectations{
				expectStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "insert ok",
			method: http.MethodPost,
			path:   "/categories",
			body:   `{"name":"Books","description":"All books"}`,
			setupSvc: func(svc *fakeCategoriesService) {
				svc.insertResp = entity.Category{ID: 2, Name: "Books"}
			},
			expect: expectations{
				expectStatus: http.StatusCreated,
				calls:        callCounts{insert: 1},
				insertName:   stringPtr("Books"),
			},
		},
		{
			name:   "insert bad json",
			method: http.MethodPost,
			path:   "/categories",
			body:   `{`,
			expect: expectations{
				expectStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "update ok",
			method: http.MethodPut,
			path:   "/categories/10",
			body:   `{"name":"New","description":"D"}`,
			setupSvc: func(svc *fakeCategoriesService) {
				svc.updateResp = entity.Category{ID: 10, Name: "New"}
			},
			expect: expectations{
				expectStatus: http.StatusOK,
				calls:        callCounts{update: 1},
				updateID:     int64Ptr(10),
			},
		},
		{
			name:   "update bad id",
			method: http.MethodPut,
			path:   "/categories/xyz",
			body:   `{"name":"New","description":"D"}`,
			expect: expectations{
				expectStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "update bad json",
			method: http.MethodPut,
			path:   "/categories/11",
			body:   `{`,
			expect: expectations{
				expectStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "delete ok",
			method: http.MethodDelete,
			path:   "/categories/12",
			setupSvc: func(svc *fakeCategoriesService) {
				svc.deleteResp = 12
			},
			expect: expectations{
				expectStatus: http.StatusOK,
				calls:        callCounts{del: 1},
				deleteID:     int64Ptr(12),
			},
		},
		{
			name:   "delete bad id",
			method: http.MethodDelete,
			path:   "/categories/nope",
			expect: expectations{
				expectStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "reference",
			method: http.MethodGet,
			path:   "/categories/reference",
			expect: expectations{
				expectStatus:  http.StatusOK,
				bodyContains:  "Simple API",
				needsRepoRoot: true,
			},
		},
		{
			name:   "method mismatch",
			method: http.MethodPost,
			path:   "/categories/health",
			expect: expectations{
				expectStatus: http.StatusNotFound,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &fakeCategoriesService{}
			if tc.setupSvc != nil {
				tc.setupSvc(svc)
			}

			handler, err := categoriesHandler.NewCategoriesHandler(svc)
			if err != nil {
				t.Fatalf("unexpected handler error: %v", err)
			}

			router := NewRouter(handler)
			mux := router.RegisterRoutes()
			if mux == nil {
				t.Fatal("expected mux")
			}

			if tc.expect.needsRepoRoot {
				root := findRepoRoot(t)
				cwd, err := os.Getwd()
				if err != nil {
					t.Fatalf("getwd: %v", err)
				}
				if err := os.Chdir(root); err != nil {
					t.Fatalf("chdir: %v", err)
				}
				defer func() {
					if err := os.Chdir(cwd); err != nil {
						t.Fatalf("restore chdir: %v", err)
					}
				}()
			}

			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.expect.expectStatus {
				t.Fatalf("expected status %d, got %d", tc.expect.expectStatus, rec.Code)
			}

			if tc.expect.bodyContains != "" && !strings.Contains(rec.Body.String(), tc.expect.bodyContains) {
				t.Fatalf("expected body to contain %q", tc.expect.bodyContains)
			}

			if svc.apiCalls != tc.expect.calls.api {
				t.Fatalf("expected api calls %d, got %d", tc.expect.calls.api, svc.apiCalls)
			}
			if svc.getAllCalls != tc.expect.calls.getAll {
				t.Fatalf("expected getAll calls %d, got %d", tc.expect.calls.getAll, svc.getAllCalls)
			}
			if svc.getByIDCalls != tc.expect.calls.getByID {
				t.Fatalf("expected getByID calls %d, got %d", tc.expect.calls.getByID, svc.getByIDCalls)
			}
			if svc.insertCalls != tc.expect.calls.insert {
				t.Fatalf("expected insert calls %d, got %d", tc.expect.calls.insert, svc.insertCalls)
			}
			if svc.updateCalls != tc.expect.calls.update {
				t.Fatalf("expected update calls %d, got %d", tc.expect.calls.update, svc.updateCalls)
			}
			if svc.deleteCalls != tc.expect.calls.del {
				t.Fatalf("expected delete calls %d, got %d", tc.expect.calls.del, svc.deleteCalls)
			}

			if tc.expect.getByID != nil && svc.lastGetByID != *tc.expect.getByID {
				t.Fatalf("expected getByID %d, got %d", *tc.expect.getByID, svc.lastGetByID)
			}
			if tc.expect.updateID != nil && svc.lastUpdate.ID != *tc.expect.updateID {
				t.Fatalf("expected update id %d, got %d", *tc.expect.updateID, svc.lastUpdate.ID)
			}
			if tc.expect.insertName != nil && svc.lastInsert.Name != *tc.expect.insertName {
				t.Fatalf("expected insert name %q, got %q", *tc.expect.insertName, svc.lastInsert.Name)
			}
			if tc.expect.deleteID != nil && svc.lastDelete != *tc.expect.deleteID {
				t.Fatalf("expected delete id %d, got %d", *tc.expect.deleteID, svc.lastDelete)
			}
		})
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "docs", "swagger.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("docs/swagger.json not found")
		}
		dir = parent
	}
}
