package handlers

import (
	"avitoTT/internal/service"
)

// Container will hold all dependencies for your application.
type Container struct {
	AuthService service.AuthServiceImpl
	BuyService  service.BuyServiceImpl
}

// NewContainer returns an empty or an initialized container for your handlers.
func NewContainer(container service.ServiceContainer) (*Container, error) {
	return &Container{
		AuthService: *container.AuthService,
		BuyService:  *container.BuyService,
	}, nil
}
