package notifying

import "time"

type Trigger struct {
	ID        int       `json:"id"`
	Symbol    string    `json:"symbol"`
	Threshold string    `json:"threshold"`
	Disabled  bool      `json:"disabled"`
	Created   time.Time `json:"created_at"`
}

type TriggerRepository interface {
	GetAll() ([]Trigger, error)
	Save(*Trigger) error
}
