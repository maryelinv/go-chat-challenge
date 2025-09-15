package tests

import (
	"testing"

	"github.com/maryelinv/go-chat-challenge/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestSaveMessageAndLastNMessages(t *testing.T) {
	g, err := db.Open(":memory:")
	assert.NoError(t, err)
	err = db.Migrate(g)
	assert.NoError(t, err)
	m := &db.Message{Room: "test", Username: "user", Text: "Hi"}
	err = db.SaveMessage(g, m)
	assert.NoError(t, err)
	msgs, err := db.LastNMessages(g, "test", 10)
	assert.NoError(t, err)
	assert.NotNil(t, msgs)
}
