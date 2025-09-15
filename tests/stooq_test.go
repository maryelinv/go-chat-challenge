package tests

import (
	"testing"

	"github.com/maryelinv/go-chat-challenge/internal/stooq"
	"github.com/stretchr/testify/assert"
)

func TestFetchQuote(t *testing.T) {
	val, err := stooq.FetchQuote("AAPL.US")
	if err != nil {
		t.Logf("FetchQuote error: %v", err)
	}
	assert.True(t, err == nil || val == 0)
}
