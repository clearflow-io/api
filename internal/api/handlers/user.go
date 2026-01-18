package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/auth"
	"github.com/igorschechtel/clearflow-backend/internal/services"
	u "github.com/igorschechtel/clearflow-backend/internal/utils"
)

type UserHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService services.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
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
	users, err := h.userService.List(r.Context(), queryParams.Limit, queryParams.Offset)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parsing
	type CreateUserRequest struct {
		ClerkID   string  `json:"clerkId" validate:"required"`
		Email     string  `json:"email" validate:"required,email"`
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		ImageURL  *string `json:"imageUrl"`
	}

	var body CreateUserRequest
	if err := u.ParseJSON(r, &body, true); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// Validation
	if err := h.validate.Struct(body); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, u.FormatValidationErrors(err))
		return
	}

	// Security check: ensure user is only creating/updating their own profile
	authenticatedClerkID, ok := auth.GetUserID(r.Context())
	if !ok || authenticatedClerkID != body.ClerkID {
		u.WriteJSONError(w, http.StatusForbidden, errors.New("cannot create or update profile for another user"))
		return
	}

	user := model.User{
		ClerkID:   body.ClerkID,
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		ImageURL:  body.ImageURL,
	}

	// Creating (using Upsert for idempotency)
	createdUser, err := h.userService.Upsert(r.Context(), &user)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	u.WriteJSON(w, http.StatusCreated, createdUser)
}
