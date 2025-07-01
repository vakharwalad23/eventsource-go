package domain

import "time"

type EventType string

const (
	AccountCreated EventType = "AccountCreated"
	MoneyDeposited EventType = "MoneyDeposited"
	MoneyWithdrawn EventType = "MoneyWithdrawn"
)

type Event struct {
	Type      EventType `json:"type"`
	AccountID string    `json:"account_id"`
	Amount    float64   `json:"amount,omitempty"`
	Time      time.Time `json:"time"`
}
