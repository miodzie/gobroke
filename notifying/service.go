package notifying

type Service interface {
	SendNotification(trig Trigger, price NewPrice) []error
}

type service struct {
	// TODO: fill with configs to make EmailNotifier.
	r Repository
}

type Repository interface {
	GetRecipientsByTriggerID(trigID int) ([]Recipient, error)
}

func NewService(repo Repository) Service {
	return &service{r: repo}
}

func (s service) SendNotification(trig Trigger, price NewPrice) []error {
	errs := []error{}
	recipients, err := s.r.GetRecipientsByTriggerID(trig.ID)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for _, r := range recipients {
		err = NotifyRecipient(r, price)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
