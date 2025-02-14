package handlers

import (
	"avitoTT/internal/service"
)

// Container will hold all dependencies for your application.
type Container struct {
	AuthService service.AuthServiceImpl
	BuyService  service.BuyServiceImpl
	InfoService service.InfoServiceImpl
	SendService service.SendServiceImpl
}

// NewContainer returns an empty or an initialized container for your handlers.
func NewContainer(container service.ServiceContainer) (*Container, error) {
	return &Container{
		AuthService: *container.AuthService,
		BuyService:  *container.BuyService,
		InfoService: *container.InfoService,
		SendService: *container.SendService,
	}, nil
}
