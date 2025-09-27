package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"lab1-rsoi/internal/dto"
	v1 "lab1-rsoi/internal/handlers/http/v1"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPersonService struct {
	CreateFn func(ctx context.Context, req *dto.CreatePersonRequest) (uint64, error)
	ListFn   func(ctx context.Context) ([]dto.PersonResponse, error)
	GetFn    func(ctx context.Context, id uint64) (*dto.PersonResponse, error)
	DeleteFn func(ctx context.Context, id uint64) error
	UpdateFn func(ctx context.Context, id uint64, req dto.PersonResponse) (*dto.PersonResponse, error)
}

func (m *mockPersonService) Create(ctx context.Context, req *dto.CreatePersonRequest) (uint64, error) {
	return m.CreateFn(ctx, req)
}
func (m *mockPersonService) List(ctx context.Context) ([]dto.PersonResponse, error) {
	return m.ListFn(ctx)
}
func (m *mockPersonService) Get(ctx context.Context, id uint64) (*dto.PersonResponse, error) {
	return m.GetFn(ctx, id)
}
func (m *mockPersonService) Update(ctx context.Context, id uint64, req dto.PersonResponse) (*dto.PersonResponse, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, req)
	}
	return &req, nil
}
func (m *mockPersonService) Delete(ctx context.Context, id uint64) error {
	return m.DeleteFn(ctx, id)
}

func setupRouter(h *v1.PersonHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

func ptr[T any](v T) *T { return &v }

func TestCreatePerson_Success(t *testing.T) {
	mockSvc := &mockPersonService{
		CreateFn: func(ctx context.Context, req *dto.CreatePersonRequest) (uint64, error) {
			return 42, nil
		},
	}
	handler := v1.New(mockSvc)
	router := setupRouter(handler)

	body := `{"name":"John","age":30,"address":"NY","work":"Dev"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/persons", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/api/v1/persons/42", w.Header().Get("Location"))
}

func TestCreatePerson_BadRequest(t *testing.T) {
	mockSvc := &mockPersonService{}
	handler := v1.New(mockSvc)
	router := setupRouter(handler)

	body := `{"name":123}` // некорректный тип
	req := httptest.NewRequest(http.MethodPost, "/api/v1/persons", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListPersons_Success(t *testing.T) {
	mockSvc := &mockPersonService{
		ListFn: func(ctx context.Context) ([]dto.PersonResponse, error) {
			return []dto.PersonResponse{
				{ID: 1, Name: ptr("John"), Age: ptr(25)},
				{ID: 2, Name: ptr("Jane"), Age: ptr(30)},
			}, nil
		},
	}
	handler := v1.New(mockSvc)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/persons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var persons []dto.PersonResponse
	err := json.Unmarshal(w.Body.Bytes(), &persons)
	assert.NoError(t, err)
	assert.Len(t, persons, 2)
	assert.Equal(t, "John", *persons[0].Name)
}

func TestGetPerson_NotFound(t *testing.T) {
	mockSvc := &mockPersonService{
		GetFn: func(ctx context.Context, id uint64) (*dto.PersonResponse, error) {
			return nil, nil
		},
	}
	handler := v1.New(mockSvc)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/persons/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePerson_Error(t *testing.T) {
	mockSvc := &mockPersonService{
		DeleteFn: func(ctx context.Context, id uint64) error {
			return errors.New("db error")
		},
	}
	handler := v1.New(mockSvc)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/persons/"+strconv.Itoa(1), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
