package models

type InfoResponseCoinHistory struct {

	Received []InfoResponseCoinHistoryReceivedInner `json:"received,omitempty"`

	Sent []InfoResponseCoinHistorySentInner `json:"sent,omitempty"`
}
