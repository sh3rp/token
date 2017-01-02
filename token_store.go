package token

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

// name of the password table
const PASSWORD_BUCKET = "user_auth"

// name of the token table
const TOKEN_BUCKET = "token"

// TokenStore
type TokenStore struct {
	DataStore      *bolt.DB
	ExpirationTime time.Duration
}

// NewTokenStore creates a TokenStore struct and initializes the database if
// necessary with the password and token tables
func NewTokenStore(dirname string, dbName string, expirationTime time.Duration) (*TokenStore, error) {
	if dirname == "" {
		dirname = "/tmp"
	}

	db, err := bolt.Open(dirname+"/"+dbName+".db", 0600, nil)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, bucketErr := tx.CreateBucketIfNotExists([]byte(PASSWORD_BUCKET))

		return bucketErr
	})

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, bucketErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET))

		return bucketErr
	})

	if err != nil {
		return nil, err
	}

	return &TokenStore{
		DataStore:      db,
		ExpirationTime: expirationTime,
	}, nil
}

// AddUser creates a new user in the password table
func (ts *TokenStore) AddUser(id string, password string) error {
	return ts.storeUser(id, HashPassword(password))
}

// DeleteUser deletes a user from the password table
func (ts *TokenStore) DeleteUser(id string) error {
	return ts.deleteUser(id)
}

// NewToken generates a new token if the specified user exists in the
// password database then generates a token and stores the token in the
// token table
func (ts *TokenStore) NewToken(id string, password string) (*Token, error) {
	if !ts.isValidPassword(id, HashPassword(password)) {
		return nil, errors.New("Invalid password")
	}

	token := GenerateTokenNow()

	err := ts.DataStore.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TOKEN_BUCKET))
		err := bucket.Put([]byte(id), token)
		return err
	})

	return &Token{Id: id, Hash: token}, err
}

// IsValid returns true if the specified token has 1) a valid user id, 2) the
// token associated with that user matches the token specified, 3) the token
// time is not stale.  Staleness is determined by comparing the
// (tokenTime + expirationTime) < currentTime
func (ts *TokenStore) IsValid(token *Token) bool {
	var dbToken []byte

	err := ts.DataStore.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TOKEN_BUCKET))
		dbToken = bucket.Get([]byte(token.Id))

		if dbToken == nil {
			return errors.New("No such token")
		}

		return nil
	})

	if err != nil {
		return false
	}

	if !bytes.Equal(dbToken, token.Hash) {
		return false
	}

	if GetTokenTime(token.Hash).Add(ts.ExpirationTime).Before(time.Now()) {
		ts.DataStore.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(TOKEN_BUCKET))
			bucket.Delete([]byte(token.Id))
			return nil
		})
		return false
	}

	return true
}

// HashPassword generates a SHA-256 based on the specified string
func HashPassword(str string) []byte {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hash.Sum(nil)
}

// storeUser stores the specified id and byte hash in the password table
func (ts *TokenStore) storeUser(id string, passwordHash []byte) error {
	return ts.DataStore.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PASSWORD_BUCKET))
		err := bucket.Put([]byte(id), passwordHash)
		return err
	})
}

// deleteUser deletes the specified user id from the password database
func (ts *TokenStore) deleteUser(id string) error {
	return ts.DataStore.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PASSWORD_BUCKET))
		err := bucket.Delete([]byte(id))
		return err
	})
}

func (ts *TokenStore) isValidPassword(id string, passwordHash []byte) bool {
	var hash []byte
	err := ts.DataStore.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PASSWORD_BUCKET))

		hash = bucket.Get([]byte(id))
		return nil
	})

	if err != nil {
		return false
	}

	return bytes.Equal(hash, passwordHash)
}
