package handler

import (
	"context"

	"github.com/aim/aim/internal/pkg/database"
	"github.com/aim/aim/internal/user/model"
	"github.com/aim/aim/internal/user/service"
)

// UserHandler provides the gRPC-compatible handler for user operations.
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{svc: service.NewUserService(db)}
}

// RegisterRequest is the input for user registration.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// LoginRequest is the input for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse is the public user info returned to clients.
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// Register handles user registration.
func (h *UserHandler) Register(_ context.Context, req *RegisterRequest) (*UserResponse, error) {
	u, err := h.svc.Register(req.Username, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}
	return toUserResponse(u), nil
}

// Login handles user login.
func (h *UserHandler) Login(_ context.Context, req *LoginRequest) (*UserResponse, error) {
	u, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return toUserResponse(u), nil
}

// GetProfile returns user profile by ID.
func (h *UserHandler) GetProfile(_ context.Context, userID uint) (*UserResponse, error) {
	u, err := h.svc.GetProfile(userID)
	if err != nil {
		return nil, err
	}
	return toUserResponse(u), nil
}

// UpdateProfile updates user info.
func (h *UserHandler) UpdateProfile(_ context.Context, u *model.User) error {
	return h.svc.UpdateProfile(u)
}

// SetOnline updates the online status.
func (h *UserHandler) SetOnline(_ context.Context, userID uint, online bool) error {
	return h.svc.SetOnline(userID, online)
}

func toUserResponse(u *model.User) *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Email:    u.Email,
		Phone:    u.Phone,
	}
}
