package models

type InfoResponseCoinHistorySentInner struct {

	// Имя пользователя, которому отправлены монеты.
	ToUser string `json:"toUser,omitempty"`

	// Количество отправленных монет.
	Amount int32 `json:"amount,omitempty"`
}
