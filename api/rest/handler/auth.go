package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"horizonx-server/internal/config"
	"horizonx-server/internal/core/auth"
)

type AuthHandler struct {
	svc auth.AuthService
	cfg *config.Config
}

func NewAuthHandler(svc auth.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		svc: svc,
		cfg: cfg,
	}
}

type APIResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Message: "Invalid request body",
		})
		return
	}

	if err := h.svc.Register(r.Context(), req); err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			writeJSON(w, http.StatusConflict, APIResponse{
				Message: "Email already registered",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Message: "Something went wrong, please try again later",
		})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{
		Message: "User created successfully.",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Message: "Invalid request body",
		})
		return
	}

	res, err := h.svc.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, APIResponse{
				Message: "Invalid credentials",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Message: "Something went wrong, please try again later",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(h.cfg.JWTExpiry),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, APIResponse{
		Message: "Login successful",
		Data:    res.User,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusNoContent, APIResponse{})
}
