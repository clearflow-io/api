package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
	u "github.com/igorschechtel/finance-tracker-backend/internal/utils"
)

type ExpenseHandler struct {
	expenseRepo repositories.ExpenseRepository
}

func NewExpenseHandler(expenseRepo repositories.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{expenseRepo: expenseRepo}
}

func (h *ExpenseHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		u.WriteJSONError(w, http.StatusBadRequest, []string{"userId path parameter is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid userId format"})
		return
	}

	type ListExpensesRequest struct {
		Limit  int `json:"limit" validate:"min=1,max=100"`
		Offset int `json:"offset" validate:"min=0"`
	}
	queryParams := ListExpensesRequest{
		Limit:  100,
		Offset: 0,
	}

	if errors, err := u.ParseQueryParamInt(r, "limit", &queryParams.Limit); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}
	if errors, err := u.ParseQueryParamInt(r, "offset", &queryParams.Offset); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	// Validation
	if err := validate.Struct(queryParams); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Fetching
	expenses, err := h.expenseRepo.ListByUser(r.Context(), userID, queryParams.Limit, queryParams.Offset)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to list expenses"})
		return
	}

	fmt.Println("expenses", expenses)

	u.WriteJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		u.WriteJSONError(w, http.StatusBadRequest, []string{"UserID path parameter is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid UserID format"})
		return
	}

	type CreateExpenseRequest struct {
		Amount       float64 `json:"amount" validate:"required,min=0"`
		Description  string  `json:"description" validate:"required,min=1,max=255"`
		PurchaseDate string  `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
		BillDate     string  `json:"billDate" validate:"required,datetime=2006-01-02"`
		CategoryID   *int32  `json:"categoryId"`
	}

	reqBody := CreateExpenseRequest{
		Amount:       0,
		Description:  "",
		PurchaseDate: "",
		BillDate:     "",
		CategoryID:   nil,
	}

	if errors, err := u.ParseJSON(r, &reqBody, true); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	purchaseDate, err := time.Parse("2006-01-02", reqBody.PurchaseDate)
	if err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid purchaseDate format"})
		return
	}

	billDate, err := time.Parse("2006-01-02", reqBody.BillDate)
	if err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid billDate format"})
		return
	}

	// Validation
	if err := validate.Struct(reqBody); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Creating
	modelExpense := &model.Expense{
		UserID:       userID,
		Amount:       reqBody.Amount,
		Description:  reqBody.Description,
		PurchaseDate: purchaseDate,
		BillDate:     billDate,
		CategoryID:   reqBody.CategoryID,
	}

	createdExpense, err := h.expenseRepo.Create(r.Context(), modelExpense)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	u.WriteJSON(w, http.StatusOK, createdExpense)
}
