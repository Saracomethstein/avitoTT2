package models

type InfoResponse struct {
	Coins       int32                        `json:"coins,omitempty"`
	Inventory   []InfoResponseInventoryInner `json:"inventory,omitempty"`
	CoinHistory InfoResponseCoinHistory      `json:"coinHistory,omitempty"`
}
