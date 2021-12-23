package notifying

import (
	"testing"
)

type mockRepo struct {
	d []Recipient
}

func (m mockRepo) GetDeliverablesByTriggerID(trigId int) ([]Recipient, error) {
	return m.d, nil
}

var mR = &mockRepo{}
var testService = &service{r: mR}

// TODO: use mocks
func TestSendNotificationWillSendAllDeliverablesForAGivenTrigger(t *testing.T) {
	// Arrange
	trig := Trigger{ID: 1}
	price := NewPrice{}

	// Act
	testService.SendNotification(trig, price)

	// Assert

}
