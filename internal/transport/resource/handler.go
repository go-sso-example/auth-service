package resource

import (
	"github.com/go-sso-example/auth-service/internal/service/resource"
	"github.com/go-sso-example/auth-service/internal/service/service"
	"github.com/go-sso-example/auth-service/internal/service/user"
)

type APIHandler struct {
	serviceService  *service.Service
	resourceService *resource.Service
	userService     *user.Service
}

func NewAPIHandler(serviceService *service.Service, resourceService *resource.Service, userService *user.Service) *APIHandler {
	return &APIHandler{
		serviceService:  serviceService,
		resourceService: resourceService,
		userService:     userService,
	}
}
