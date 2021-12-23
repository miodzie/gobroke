package notifying

type Recipient struct {
	Email       string
	PhoneNumber string
	JSONPayload string
	Type        DeliveryType
}
