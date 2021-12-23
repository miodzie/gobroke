package memory

import "github.com/miodzie/gobroke/notifying"

type Storage struct {
	trigs      []*Trigger
	prices     []*NewPrice
	Recipients []Recipient
}

func (m *Storage) GetRecipientsByTriggerID(trigID int) ([]notifying.Recipient, error) {
	var recipients []notifying.Recipient
	for _, d := range m.Recipients {
		r := notifying.Recipient{
			Email:       d.Email,
			PhoneNumber: d.PhoneNumber,
			JSONPayload: d.JSONPayload,
			Type:        d.Type,
		}

		recipients = append(recipients, r)
	}

	return recipients, nil
}

func (m *Storage) SaveTrigger(triggers ...*Trigger) error {
	m.trigs = append(m.trigs, triggers...)

	return nil
}

func (m *Storage) SaveNewPrice(prices ...*NewPrice) error {
	m.prices = append(m.prices, prices...)

	return nil
}

func (m *Storage) GetTriggers() ([]*Trigger, error) {
	return []*Trigger{}, nil
}
