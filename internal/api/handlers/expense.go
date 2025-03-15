package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
	"github.com/igorschechtel/finance-tracker-backend/internal/utils"
)

type ExpenseHandler struct {
	expenseRepo repositories.ExpenseRepository
}

func NewExpenseHandler(expenseRepo repositories.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{expenseRepo: expenseRepo}
}

func (h *ExpenseHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Extract userId from the path parameters using chi
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, []string{"userId path parameter is required"})
		return
	}

	// Parse userId into a UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid userId format"})
		return
	}

	// Define the request struct for query parameters
	type ListExpensesRequest struct {
		Limit  int `json:"limit" validate:"min=1,max=100"`
		Offset int `json:"offset" validate:"min=0"`
	}

	queryParams := ListExpensesRequest{
		Limit:  10,
		Offset: 0,
	}

	// Parse and validate query parameters
	if errors, err := utils.ParseQueryParamInt(r, "limit", &queryParams.Limit); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	if errors, err := utils.ParseQueryParamInt(r, "offset", &queryParams.Offset); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	if err := validate.Struct(queryParams); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, utils.FormatValidationErrors(err))
		return
	}

	// Fetch expenses from the repository
	expenses, err := h.expenseRepo.ListByUser(r.Context(), userID, queryParams.Limit, queryParams.Offset)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to list expenses"})
		return
	}

	fmt.Println("expenses", expenses)

	// Write the response
	utils.WriteJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Extract userId from the path parameters using chi
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, []string{"userId path parameter is required"})
		return
	}

	// Parse userId into a UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid userId format"})
		return
	}

	// Define the request struct for query parameters
	type CreateExpenseRequest struct {
		Amount       float64 `json:"amount" validate:"required,min=0"`
		Description  string  `json:"description" validate:"required,min=1,max=255"`
		PurchaseDate string  `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
		BillDate     string  `json:"billDate" validate:"required,datetime=2006-01-02"`
		CategoryID   *int    `json:"categoryId"`
	}

	expenseReq := CreateExpenseRequest{
		Amount:       0,
		Description:  "",
		PurchaseDate: "",
		BillDate:     "",
		CategoryID:   nil,
	}

	// Parse the request body into the expense struct
	if errors, err := utils.ParseJSON(r, &expenseReq, true); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	// Validate the expense struct
	if err := validate.Struct(expenseReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, utils.FormatValidationErrors(err))
		return
	}

	purchaseDate, err := time.Parse("2006-01-02", expenseReq.PurchaseDate)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid purchaseDate format"})
		return
	}

	billDate, err := time.Parse("2006-01-02", expenseReq.BillDate)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, []string{"Invalid billDate format"})
		return
	}

	modelExpense := &model.Expense{
		Amount:       expenseReq.Amount,
		Description:  expenseReq.Description,
		PurchaseDate: purchaseDate,
		BillDate:     billDate,
		CategoryID:   int32(utils.IntPointerToValue(expenseReq.CategoryID)),
	}

	createdExpense, err := h.expenseRepo.Create(r.Context(), userID, modelExpense)
	if err != nil {
		fmt.Printf("%+v\n", modelExpense)
		fmt.Println("error", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to create expense"})
		return
	}

	// Write the response
	utils.WriteJSON(w, http.StatusOK, createdExpense)
}
