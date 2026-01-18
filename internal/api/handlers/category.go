package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/auth"
	"github.com/igorschechtel/clearflow-backend/internal/services"
	u "github.com/igorschechtel/clearflow-backend/internal/utils"
)

type CategoryHandler struct {
	categoryService services.CategoryService
	validate        *validator.Validate
}

func NewCategoryHandler(categoryService services.CategoryService, validate *validator.Validate) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validate:        validate,
	}
}

func (h *CategoryHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	clerkID, ok := auth.GetUserID(r.Context())
	if !ok {
		u.WriteJSONError(w, http.StatusUnauthorized, u.ErrUnauthorized)
		return
	}

	type ListCategoriesRequest struct {
		Limit  int `json:"limit" validate:"min=1,max=100"`
		Offset int `json:"offset" validate:"min=0"`
	}
	queryParams := ListCategoriesRequest{
		Limit:  100,
		Offset: 0,
	}

	if err := u.ParseQueryParamInt(r, &queryParams.Limit, "limit", false); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}
	if err := u.ParseQueryParamInt(r, &queryParams.Offset, "offset", false); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// Validation
	if err := h.validate.Struct(queryParams); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Fetching
	categories, err := h.categoryService.ListByUser(r.Context(), clerkID, queryParams.Limit, queryParams.Offset)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	clerkID, ok := auth.GetUserID(r.Context())
	if !ok {
		u.WriteJSONError(w, http.StatusUnauthorized, u.ErrUnauthorized)
		return
	}

	type CreateCategoryRequest struct {
		Name        string `json:"name" validate:"required,min=1,max=255"`
		Description string `json:"description" validate:"max=255"`
		ColorHex    string `json:"colorHex" validate:"required,min=7,max=7"`
	}

	reqBody := CreateCategoryRequest{}
	if err := u.ParseJSON(r, &reqBody, true); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// Validation
	if err := h.validate.Struct(reqBody); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Creating
	modelCategory := &model.Category{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		ColorHex:    reqBody.ColorHex,
	}

	createdCategory, err := h.categoryService.Create(r.Context(), clerkID, modelCategory)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, createdCategory)
}
