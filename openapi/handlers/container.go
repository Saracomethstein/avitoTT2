package handlers

import "avitoTT/internal/service"

// Container will hold all dependencies for your application.
type Container struct {
	AuthService service.AuthService
}

// NewContainer returns an empty or an initialized container for your handlers.
func NewContainer() (Container, error) {
	c := Container{
		AuthService: &service.AuthServiceImpl{},
	}
	return c, nil
}
