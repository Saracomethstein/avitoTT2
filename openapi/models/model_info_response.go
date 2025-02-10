package models

type InfoResponse struct {

	// Количество доступных монет.
	Coins int32 `json:"coins,omitempty"`

	Inventory []InfoResponseInventoryInner `json:"inventory,omitempty"`

	CoinHistory InfoResponseCoinHistory `json:"coinHistory,omitempty"`
}
