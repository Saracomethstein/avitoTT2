package models

type AuthRequest struct {

	// Имя пользователя для аутентификации.
	Username string `json:"username"`

	// Пароль для аутентификации.
	Password string `json:"password"`
}
