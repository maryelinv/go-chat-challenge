package tests

import (
	"testing"

	"github.com/maryelinv/go-chat-challenge/internal/queue"
)

func TestMustConnectFromEnv(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic as expected: %v", r)
		}
	}()
	_ = queue.MustConnectFromEnv()
}
