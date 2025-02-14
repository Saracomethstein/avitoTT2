package models

type SendCoinRequest struct {

	// Имя пользователя, которому нужно отправить монеты.
	ToUser string `json:"toUser" validate:"required"`

	// Количество монет, которые необходимо отправить.
	Amount int32 `json:"amount" validate:"required"`
}
