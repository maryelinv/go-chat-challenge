package tests

import (
	"testing"

	"github.com/maryelinv/go-chat-challenge/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "password123"
	hash, err := auth.Hash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	err = auth.Check(hash, password)
	assert.NoError(t, err)
	err = auth.Check(hash, "wrongpassword")
	assert.Error(t, err)
}
