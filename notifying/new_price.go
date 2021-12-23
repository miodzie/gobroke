package notifying

import (
	"time"
)

type NewPrice struct {
	Symbol string
	Price  string
	Time   time.Time
}
