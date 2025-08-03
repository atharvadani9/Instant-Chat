package api

import (
	"chat/internal/store"
	"chat/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	Store  store.UserStore
	logger *log.Logger
}

func NewUserHandler(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{Store: store, logger: logger}
}

func (h *UserHandler) GetMeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Envelope{"error": "Method not allowed"})
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Printf("ERROR: user_id query parameter is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "User ID query parameter is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Printf("ERROR: invalid user_id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	user, err := h.Store.GetUserByID(userID)
	if err != nil {
		h.logger.Printf("ERROR: user not found: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "User not found"})
		return
	}

	h.logger.Printf("INFO: User: %v retrieved successfully", user.Username)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
	return
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Envelope{"error": "Method not allowed"})
		return
	}

	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		h.logger.Printf("ERROR: username or password is empty")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username and password are required"})
		return
	}

	// Check if user already exists
	existingUser, err := h.Store.GetUserByUsername(req.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.logger.Printf("ERROR: checking existing user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if existingUser != nil {
		h.logger.Printf("ERROR: user already exists: %s", req.Username)
		utils.WriteJSON(w, http.StatusConflict, utils.Envelope{"error": "Username already exists"})
		return
	}

	// Hash password
	passwordHash, err := h.Store.HashPassword(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	// Create user
	user := &store.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	err = h.Store.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create user"})
		return
	}

	h.logger.Printf("INFO: user created successfully: %s", user.Username)
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Envelope{"error": "Method not allowed"})
		return
	}

	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username and password are required"})
		return
	}

	// Authenticate user
	user, err := h.Store.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Printf("INFO: login attempt with invalid username: %s", req.Username)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
			return
		}
		h.logger.Printf("ERROR: authenticating user: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	h.logger.Printf("INFO: user logged in successfully: %s", user.Username)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Envelope{"error": "Method not allowed"})
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Printf("ERROR: user_id query parameter is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "User ID query parameter is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Printf("ERROR: invalid user_id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	_, err = h.Store.GetUserByID(userID)
	if err != nil {
		h.logger.Printf("ERROR: user not found: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "User not found"})
		return
	}

	users, err := h.Store.GetUsersExcept(userID)
	if err != nil {
		h.logger.Printf("ERROR: getting users: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to get users"})
		return
	}

	h.logger.Printf("INFO: users retrieved successfully")
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"users": users})
}
