package queue

import (
	"context"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	ch         *amqp.Channel
	reqQ, repQ amqp.Queue
}

func MustConnectFromEnv() *Client {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	reqQ, err := ch.QueueDeclare("stock_requests", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	repQ, err := ch.QueueDeclare("stock_replies", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	return &Client{ch: ch, reqQ: reqQ, repQ: repQ}
}

func (c *Client) PublishRequest(code, room string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.ch.PublishWithContext(ctx, "", c.reqQ.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(room + "|" + code),
	})
}

func (c *Client) PublishReply(room, text string) error {
	return c.ch.Publish("", c.repQ.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(room + "|" + text),
	})
}

func (c *Client) ConsumeReplies(handle func(room, text string)) {
	msgs, _ := c.ch.Consume(c.repQ.Name, "", true, false, false, false, nil)
	go func() {
		for m := range msgs {
			body := string(m.Body)
			i := 0
			for i < len(body) && body[i] != '|' {
				i++
			}
			room, text := body[:i], body[i+1:]
			handle(room, text)
		}
	}()
}

func (c *Client) ConsumeRequests() (<-chan amqp.Delivery, error) {
	return c.ch.Consume(c.reqQ.Name, "", true, false, false, false, nil)
}
