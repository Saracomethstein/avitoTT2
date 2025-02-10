package models

type AuthResponse struct {

	// JWT-токен для доступа к защищенным ресурсам.
	Token string `json:"token,omitempty"`
}
