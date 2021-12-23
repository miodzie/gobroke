package notifying

import (
	"fmt"
)

// Notifier notifies the recipient of the NewPrice.
type Notifier interface {
	// Send sends a message to the recipient.
	Notify(recipient Recipient, price NewPrice) error
}

var notifiers map[DeliveryType]Notifier

func init() {
	notifiers = make(map[DeliveryType]Notifier)
}

// TODO: What if I want different mailer platforms?
// e.g. raw SMTP vs an API?
func RegisterNotifier(t DeliveryType, n Notifier) {
	notifiers[t] = n
}

func NotifyRecipient(recipient Recipient, price NewPrice) error {
	// Create Notifer based on Deliverable.Type
	// Send it off
	notifier, ok := notifiers[recipient.Type]
	if !ok {
		return fmt.Errorf("notifier for DeliverableType: %s is not configured", recipient.Type)
	}

	return notifier.Notify(recipient, price)
}
