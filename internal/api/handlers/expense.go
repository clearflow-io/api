package handlers

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/auth"
	"github.com/igorschechtel/clearflow-backend/internal/services"
	u "github.com/igorschechtel/clearflow-backend/internal/utils"
)

type ExpenseHandler struct {
	expenseService services.ExpenseService
	validate       *validator.Validate
}

func NewExpenseHandler(expenseService services.ExpenseService, validate *validator.Validate) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
		validate:       validate,
	}
}

func (h *ExpenseHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		u.WriteJSONError(w, http.StatusUnauthorized, u.ErrUnauthorized)
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
	expenses, err := h.expenseService.ListByUser(r.Context(), userID, queryParams.Limit, queryParams.Offset)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		u.WriteJSONError(w, http.StatusUnauthorized, u.ErrUnauthorized)
		return
	}

	type CreateExpenseRequest struct {
		Amount       float64 `json:"amount" validate:"required,min=0"`
		Description  string  `json:"description" validate:"required,min=1,max=255"`
		PurchaseDate string  `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
		BillDate     string  `json:"billDate" validate:"required,datetime=2006-01-02"`
		CategoryID   *int32  `json:"categoryId"`
	}

	reqBody := CreateExpenseRequest{}
	if err := u.ParseJSON(r, &reqBody, true); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	var purchaseDate, billDate time.Time
	if err := u.ParseIsoDate(reqBody.PurchaseDate, &purchaseDate); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}
	if err := u.ParseIsoDate(reqBody.BillDate, &billDate); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// Validation
	if err := h.validate.Struct(reqBody); err != nil {
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

	createdExpense, err := h.expenseService.Create(r.Context(), modelExpense)
	if err != nil {
		if err == u.ErrNotFound {
			u.WriteJSONError(w, http.StatusNotFound, err)
			return
		}
		if err == u.ErrForbidden {
			u.WriteJSONError(w, http.StatusForbidden, err)
			return
		}
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, createdExpense)
}
