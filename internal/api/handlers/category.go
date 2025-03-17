package handlers

import (
	"net/http"

	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/internal/auth"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
	u "github.com/igorschechtel/finance-tracker-backend/internal/utils"
)

type CategoryHandler struct {
	expenseRepo repositories.CategoryRepository
}

func NewCategoryHandler(expenseRepo repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{expenseRepo: expenseRepo}
}

func (h *CategoryHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	userID, err := auth.GetUserID(r)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
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
	if err := validate.Struct(queryParams); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Fetching
	categories, err := h.expenseRepo.ListByUser(r.Context(), userID, queryParams.Limit, queryParams.Offset)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	userID, err := auth.GetUserID(r)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
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
	if err := validate.Struct(reqBody); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Creating
	modelCategory := &model.Category{
		UserID:      userID,
		Name:        reqBody.Name,
		Description: reqBody.Description,
		ColorHex:    reqBody.ColorHex,
	}

	createdCategory, err := h.expenseRepo.Create(r.Context(), modelCategory)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, createdCategory)
}
