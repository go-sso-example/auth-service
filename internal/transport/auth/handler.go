package auth

import "github.com/go-sso-example/auth-service/internal/service/auth"

type APIHandler struct {
	authService *auth.Service
}

func NewAPIHandler(authService *auth.Service) *APIHandler {
	return &APIHandler{
		authService: authService,
	}
}
