package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/maryelinv/go-chat-challenge/internal/queue"
	"github.com/maryelinv/go-chat-challenge/internal/stooq"
)

func main() {
	q := queue.MustConnectFromEnv()

	msgs, err := q.ConsumeRequests()
	if err != nil {
		log.Fatal("consume:", err)
	}

	log.Println("bot running")
	for m := range msgs {
		body := string(m.Body)
		i := strings.IndexByte(body, '|')

		if i <= 0 || i >= len(body)-1 {
			continue
		}
		room, code := body[:i], body[i+1:]

		quote, err := stooq.FetchQuote(code)
		reply := fmt.Sprintf("could not fetch %s", strings.ToUpper(code))
		if err == nil {
			reply = fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(code), quote)
		}

		if err := q.PublishReply(room, reply); err != nil {
			log.Println("publish:", err)
		}
	}
}
