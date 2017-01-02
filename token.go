package token

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"time"
)

type Token struct {
	Id   string
	Hash []byte
}

var LETTERS = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateTokenNow() []byte {
	return GenerateToken(time.Now())
}

func GenerateToken(baseTime time.Time) []byte {
	// generate first 8 bytes (timestamp)
	tBytes := make([]byte, 16)
	t := baseTime.UnixNano()
	binary.PutVarint(tBytes, t)
	// generate random noise here; the 8 is the length of the random string,
	// not the length of the hash; the hash will always be 32 bytes long (sha-256)
	hashBytes := RandomHash(8)
	// combine together
	var totalBytes []byte
	totalBytes = append(totalBytes, tBytes...)
	totalBytes = append(totalBytes, hashBytes...)
	return totalBytes
}

func GetTokenTime(token []byte) time.Time {
	// get first 8 bytes of the slice and convert to 64 bit int
	timeInt, _ := binary.Varint(token[:16])
	t := time.Unix(0, timeInt)
	return t
}

func RandomHash(size int) []byte {
	letters := RandomLetters(size)
	hash := sha256.New()
	hash.Write(letters)
	return hash.Sum(nil)
}

func RandomLetters(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = LETTERS[rand.Intn(size)]
	}
	return bytes
}

func (t *Token) String() string {
	return "ID  : " + t.Id + "\nHash: " + hex.EncodeToString(t.Hash) + "\nTime: " + GetTokenTime(t.Hash).String()
}
