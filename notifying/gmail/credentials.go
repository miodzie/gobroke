package gmail

type GmailNotifier struct {
}

func New() GmailNotifier {
	return GmailNotifier{}
}

func (g GmailNotifier) Notify(recipient notifying.Recipient, price notifying.NewPrice) error {

	return nil
}
