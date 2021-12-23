package email

import (
	"fmt"
	"net/smtp"

	"github.com/rs/zerolog/log"
	"github.com/miodzie/gobroke/config"
	"github.com/miodzie/gobroke/notifying"
)

type EmailNotifier struct {
	conf config.EmailNotifierConfig
}

func NewEmailNotifier(conf config.EmailNotifierConfig) EmailNotifier {
	return EmailNotifier{conf: conf}
}

func (e EmailNotifier) Notify(recipient notifying.Recipient, price notifying.NewPrice) error {
	msg := []byte(fmt.Sprintf("%s is now %s!", price.Symbol, price.Price))

	log.Debug().
		Str("Email", recipient.Email).
		Msg("Attempting to email new price.")

	err := smtp.SendMail(e.conf.Host+":"+e.conf.Port, e.conf.Auth(), e.conf.Username, []string{recipient.Email}, msg)

	if err != nil {
		log.Err(err).Msg("Failed to send email.")
	}

	return err
}
