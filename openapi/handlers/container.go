package handlers

import (
	"avitoTT/internal/service"
)

// Container will hold all dependencies for your application.
type Container struct {
	AuthService service.AuthServiceImpl
}

// NewContainer returns an empty or an initialized container for your handlers.
func NewContainer(container service.ServiceContainer) (*Container, error) {
	return &Container{
		AuthService: *container.AuthService,
	}, nil
}
