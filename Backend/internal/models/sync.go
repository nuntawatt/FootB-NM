package models

type SyncRequest struct {
	Season int    `json:"season"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}
