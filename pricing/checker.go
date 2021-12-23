package pricing

import (
	"github.com/rs/zerolog/log"
	"github.com/miodzie/gobroke/notifying"
	"strings"
	"time"
)

type Checker struct {
	feed     Feed
	triggers map[string][]*notifying.Trigger
	closing  chan chan error // for Stop
	notifier notifying.Service
}

// TODO: pass Storage interface instead of triggers.
func NewChecker(notifier notifying.Service, feed Feed, triggers []*notifying.Trigger) *Checker {
	c := make(map[string][]*notifying.Trigger)

	for _, trig := range triggers {
		c[trig.Symbol] = append(c[trig.Symbol], trig)
	}

	return &Checker{feed: feed, triggers: c, notifier: notifier}
}

// Run starts the Checker, it pulls the latest prices,
// checks them against active Triggers.
// Send alert when case is met.
// Disable Trigger, and remove from list.
func (p *Checker) Run() {
	go p.loop()
}

func (p *Checker) loop() {
	for {
		select {

		case errc := <-p.closing:
			errc <- p.feed.Close()
			return

			// Update the configs every 15 seconds.
		case <-time.After(time.Second * 15):
			// TODO: update configs lel

			// Consume feed.
		case msg := <-p.feed.Updates():
			// Check through list of Triggers.
			// FIXME: this sucks
			symbol := strings.Split(msg.ProductId, "-")[0]
			trigs := p.triggers[symbol]
			for i, trig := range trigs {
				// Send Alert if needed; Disable Trigger after if sent.
				if trig.Disabled {
					//TODO: MUTEX LOCK!!?!?!
					// Remove from configs, continue.
					p.triggers[symbol] = append(trigs[:i], trigs[i+1:]...)
					continue
				}
				if msg.Price >= trig.Threshold {
					trig.Disabled = true
					go func() {
						price := notifying.NewPrice{Symbol: msg.ProductId, Price: msg.Price, Time: time.Now()}
						errs := p.notifier.SendNotification(*trig, price)
						if len(errs) > 0 {
							for _, err := range errs {
								// TODO: Retry on error?
								log.Error().
									Err(err).
									Msg("Failed to send notification on price trigger.")
							}
						}
						// TODO: Save disabled status, so we don't pull the Trigger again.
					}()
				}
			}
		}
	}
}

// Close closes the feed, then sends a chan error to the closing channel,
// and waits to return the last error, if any.
func (n *Checker) Stop() error {
	errc := make(chan error)
	n.closing <- errc
	return <-errc
}
