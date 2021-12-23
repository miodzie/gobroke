package memory

import (
	"time"

	"github.com/miodzie/gobroke/notifying"
)

type Recipient struct {
	ID          int
	Email       string
	PhoneNumber string
	JSONPayload string
	Type        notifying.DeliveryType
	Created     time.Time
}
