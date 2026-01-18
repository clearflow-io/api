package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	svix "github.com/svix/svix-webhooks/go"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/services"
	u "github.com/igorschechtel/clearflow-backend/internal/utils"
)

type ClerkWebhookHandler struct {
	userService   services.UserService
	webhookSecret string
}

func NewClerkWebhookHandler(userService services.UserService, webhookSecret string) *ClerkWebhookHandler {
	return &ClerkWebhookHandler{
		userService:   userService,
		webhookSecret: webhookSecret,
	}
}

// Clerk Webhook Data Structures
type clerkEmailAddress struct {
	EmailAddress string `json:"email_address"`
	ID           string `json:"id"`
}

type clerkUserData struct {
	ID                    string              `json:"id"`
	EmailAddresses        []clerkEmailAddress `json:"email_addresses"`
	PrimaryEmailAddressID string              `json:"primary_email_address_id"`
	FirstName             *string             `json:"first_name"`
	LastName              *string             `json:"last_name"`
	ImageURL              *string             `json:"image_url"`
}

type clerkEvent struct {
	Data clerkUserData `json:"data"`
	Type string        `json:"type"`
}

func (h *ClerkWebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if h.webhookSecret == "" {
		u.WriteJSONError(w, http.StatusInternalServerError, errors.New("Webhook secret is not configured"))
		return
	}

	// Verify signature
	headers := r.Header
	wh, err := svix.NewWebhook(h.webhookSecret)
	if err != nil {
		u.WriteJSONError(w, http.StatusInternalServerError, errors.New("Failed to initialize webhook verifier"))
		return
	}

	// Limit payload size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, errors.New("Failed to read request body or payload too large"))
		return
	}

	if err := wh.Verify(payload, headers); err != nil {
		u.WriteJSONError(w, http.StatusUnauthorized, errors.New("Invalid webhook signature"))
		return
	}

	// Parse Clerk Event
	var event clerkEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		u.WriteJSONError(w, http.StatusBadRequest, errors.New("Failed to parse webhook payload"))
		return
	}

	switch event.Type {
	case "user.created", "user.updated":
		// Extract primary email
		var primaryEmail string
		for _, email := range event.Data.EmailAddresses {
			if email.ID == event.Data.PrimaryEmailAddressID {
				primaryEmail = email.EmailAddress
				break
			}
		}

		if primaryEmail == "" && len(event.Data.EmailAddresses) > 0 {
			primaryEmail = event.Data.EmailAddresses[0].EmailAddress
		}

		user := model.User{
			ClerkID:   event.Data.ID,
			Email:     primaryEmail,
			FirstName: event.Data.FirstName,
			LastName:  event.Data.LastName,
			ImageURL:  event.Data.ImageURL,
		}

		_, err := h.userService.Upsert(r.Context(), &user)
		if err != nil {
			u.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

	case "user.deleted":
		err := h.userService.DeleteByClerkID(r.Context(), event.Data.ID)
		if err != nil {
			u.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

	default:
		// Unknown event type, just return 200 as per Clerk requirements
	}

	w.WriteHeader(http.StatusOK)
}
