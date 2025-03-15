package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/internal/repositories"
	u "github.com/igorschechtel/finance-tracker-backend/internal/utils"
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

	// Parsing
	type ListUsersRequest struct {
		Limit  int `json:"limit" validate:"min=1,max=100"`
		Offset int `json:"offset" validate:"min=0"`
	}

	queryParams := ListUsersRequest{
		Limit:  10,
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
	users, err := h.userRepo.List(r.Context(), queryParams.Limit, queryParams.Offset)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to list users"})
		return
	}

	u.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	type CreateUserRequest struct {
		ID string `json:"id" validate:"required,uuid"`
	}

	var body CreateUserRequest
	if errors, err := u.ParseJSON(r, &body, true); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, errors)
		return
	}

	// Validation
	if err := validate.Struct(body); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	user := model.User{
		ID: uuid.MustParse(body.ID),
	}

	// Creating
	createdUser, err := h.userRepo.Create(r.Context(), &user)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, []string{"Failed to create user"})
		return
	}

	u.WriteJSON(w, http.StatusCreated, createdUser)
}
