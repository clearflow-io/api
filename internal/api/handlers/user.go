package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
	"github.com/igorschechtel/finance-tracker-backend/internal/utils"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type ListUsersRequest struct {
		Limit  int `json:"limit" validate:"min=1,max=100"`
		Offset int `json:"offset" validate:"min=0"`
	}

	// Query parameters with default values
	queryParams := ListUsersRequest{
		Limit:  10,
		Offset: 0,
	}

	// Parse and validate "limit" query parameter
	if errors, err := utils.ParseQueryParamInt(r, "limit", &queryParams.Limit); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	// Parse and validate "offset" query parameter
	if errors, err := utils.ParseQueryParamInt(r, "offset", &queryParams.Offset); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	// Validate the parsed query parameters
	if err := validate.Struct(queryParams); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, utils.FormatValidationErrors(err))
		return
	}

	// Fetch users from the repository
	users, err := h.userRepo.List(r.Context(), queryParams.Limit, queryParams.Offset)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to list users"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type CreateUserRequest struct {
		ID string `json:"id" validate:"required,uuid"`
	}

	var body CreateUserRequest
	if errors, err := utils.ParseJSON(r, &body, true); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	if err := validate.Struct(body); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, utils.FormatValidationErrors(err))
		return
	}

	user := model.User{
		ID: uuid.MustParse(body.ID),
	}

	createdUser, err := h.userRepo.Create(r.Context(), &user)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to create user"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdUser)
}
