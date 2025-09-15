package tests

import (
	"testing"

	"github.com/maryelinv/go-chat-challenge/internal/chat"
)

func TestHubJoinLeave(t *testing.T) {
	hub := chat.NewHub()
	client := &chat.Client{Room: "room1", Send: make(chan []byte, 1)}
	hub.Join(client)
	hub.Leave(client)
}
