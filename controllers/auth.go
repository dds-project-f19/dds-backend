package controllers

import (
	"crypto/sha256"
	"dds-backend/database"
	"dds-backend/models"
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const ( // token generator constants
	ExpirationDuration           = time.Hour * 12
	TokenGenerationPayloadLength = 10
)

// Claim types, greater values include previous claims
const (
	Worker int = iota
	Manager
)

type Auth struct {
	Username   string    `gorm:"unique;not null"`
	Claim      int       `gorm:"not null"`
	Token      string    `gorm:"unique;not null"`
	Expiration time.Time `gorm:"not null"`
}

// accept - username password
// check credentials, if not valid return error
// if already authorized return existing token, prolong expiration
// if authorization token expired return new token, delete existing token
// if no authorization existed, return new token
func Authorize(username, password string) (string, error) {
	user := models.User{Username: username, Password: Hash(password)}
	auth := Auth{Username: username}
	tx := database.DB.Begin() // begin transaction

	if err := tx.Find(&user).Error; err != nil { // check if credentials are valid
		return "", err
	}
	now := time.Now()
	if tx.Find(&auth).RecordNotFound() { // auth record not found
		auth.Claim = user.Claim
		auth.Expiration = now.Add(ExpirationDuration)
		auth.Token = GenerateNewToken(username)
		tx.Create(&auth)
	} else { // auth record found
		if now.After(auth.Expiration) { // auth record expired
			auth.Token = GenerateNewToken(username, auth.Token)
		}
		auth.Expiration = now.Add(ExpirationDuration)
		tx.Save(&auth)
	}
	if tx.Error != nil {
		tx.Rollback()
		return "", tx.Error
	}
	return auth.Token, tx.Commit().Error
}

type TokenExpirationError struct {
	err string
}

func (e *TokenExpirationError) Error() string {
	return e.err
}

// accept - token claim
// find auth record and check it's validity:
// compare requested claim with available claim
// check for expiration, in case of error return `TokenExpirationError`
func Authenticate(token string, requiredClaim int) error {
	auth := Auth{Token: token}
	if err := database.DB.Find(&auth).Error; err != nil { // auth record found
		return err
	}
	if time.Now().After(auth.Expiration) {
		return &TokenExpirationError{err: "token has expired"}
	}
	if auth.Claim < requiredClaim {
		return errors.New("claim has insufficient rights")
	}
	return nil
}

// accept `input` string
// return hash of `input` as string
func Hash(input string) string {
	ctx := sha256.New()
	ctx.Write([]byte(input))
	output := ctx.Sum(nil)
	return hex.EncodeToString(output)
}

// generate new token for authorization
// example arguments: username, previous token
func GenerateNewToken(input ...string) string {
	base := GenerateRandomString(TokenGenerationPayloadLength)
	for _, elem := range input {
		base += elem
	}
	return Hash(base)
}

// generate string of random chars of length `length`
func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
