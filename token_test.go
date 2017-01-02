package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenNow(t *testing.T) {
	token := GenerateTokenNow()
	assert.Equal(t, 48, len(token))
}

func TestGenerateToken(t *testing.T) {
	nowTime := time.Now()
	token := GenerateToken(nowTime)
	tokenTime := GetTokenTime(token)
	assert.Equal(t, nowTime, tokenTime)
}

func TestRandomLetters(t *testing.T) {
	assert.Equal(t, 16, len(RandomLetters(16)))
}

func TestRandomHash(t *testing.T) {
	assert.Equal(t, 32, len(RandomHash(16)))
}
