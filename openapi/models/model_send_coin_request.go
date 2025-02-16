package models

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int32  `json:"amount" validate:"required"`
}
