package memory

import (
	"github.com/miodzie/gobroke/notifying"
)

type Trigger struct {
	Symbol    string
	Threshold string
}

type TriggerRepository struct {
	t []notifying.Trigger
}

func NewTriggerRepo() *TriggerRepository {
	t := &TriggerRepository{}
	return t
}

func (t *TriggerRepository) GetAll() ([]notifying.Trigger, error) {
	return t.t, nil
}

func (t *TriggerRepository) Save(trig *notifying.Trigger) error {
	t.t = append(t.t, *trig)

	return nil
}
