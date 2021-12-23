package pricing

import (
	"fmt"
	"log"

	ws "github.com/gorilla/websocket"
)

type JSON map[string]interface{}

type MessageChannel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type Message struct {
	Type      string           `json:"type"`
	ProductId string           `json:"product_id"`
	Channels  []MessageChannel `json:"channels"`
	Price     string           `json:"price"`
}

// Feed delivers Messages over a channel, Close cancels
// the feed, closes the Updates channel, and returns the last error,
// if any.
type Feed interface {
	Updates() <-chan Message
	Close() error
}

// NewCoinbaseFeed creates a new CoinBase WebSocket feed for the given currency and
// cryptos.
func NewCoinbaseFeed(currency string, symbols ...string) (Feed, error) {
	var wsDialer ws.Dialer
	con, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		return nil, err
	}

	var productIds []string
	for _, symbol := range symbols {
		productIds = append(productIds, fmt.Sprintf("%s-%s", symbol, currency))
	}

	subscribe := Message{
		Type: "subscribe",
		Channels: []MessageChannel{
			{
				Name:       "ticker",
				ProductIds: productIds,
			},
		},
	}

	if err := con.WriteJSON(subscribe); err != nil {
		return nil, err
	}

	f := &feed{
		con:     con,
		updates: make(chan Message),
		closing: make(chan chan error),
	}

	go f.loop()

	return f, nil
}

// feed implements the Feed interface.
type feed struct {
	con     *ws.Conn
	updates chan Message    // sends Messages to the user
	closing chan chan error // for Close
}

func (f *feed) Updates() <-chan Message {
	return f.updates
}

// Close sends a chan error to the closing channel,
// and waits to return the last error, if any.
func (f *feed) Close() error {
	errc := make(chan error)
	f.closing <- errc
	return <-errc
}

func (f *feed) loop() {
	var err error
	for {
		select {
		case errc := <-f.closing:
			if e := f.con.Close(); e != nil {
				err = e
			}
			close(f.updates)
			errc <- err
			return

		default: // TODO: Change to non-blocking.
			msg := Message{}
			if err = f.con.ReadJSON(&msg); err != nil {
				log.Fatalf(err.Error())
				fmt.Println(err)
				continue
			}

			// Don't send heartbeat types.
			if msg.Type == "ticker" {
				f.updates <- msg
			}
		}
	}
}
