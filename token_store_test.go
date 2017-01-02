package token

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenStore(t *testing.T) {
	cleanUp("NewTokenStore")
	store, err := NewTokenStore("", "NewTokenStore", time.Second*1)
	assert.NotNil(t, store)
	assert.Nil(t, err)
	//_, err = os.Stat("/tmp/auth.db")
	//assert.Nil(t, err)
}

func TestAddUser(t *testing.T) {
	cleanUp("AddUser")
	store, _ := NewTokenStore("", "AddUser", time.Second*1)
	err := store.AddUser("test", "password")
	assert.Nil(t, err)
}

func TestNewToken(t *testing.T) {
	cleanUp("NewToken")
	store, _ := NewTokenStore("", "NewToken", time.Second*1)
	store.AddUser("test", "password")
	token, err := store.NewToken("test", "password")
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestIsValidToken(t *testing.T) {
	cleanUp("IsValidToken")
	store, _ := NewTokenStore("", "IsValidToken", time.Second*1)
	store.AddUser("test", "password")
	token, err := store.NewToken("test", "password")
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.True(t, store.IsValid(token))
	time.Sleep(time.Second * 3)
	assert.False(t, store.IsValid(token))
}

func cleanUp(dbName string) {
	os.Remove("/tmp/" + dbName + ".db")
}
