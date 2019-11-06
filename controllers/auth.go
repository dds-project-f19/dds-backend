package controllers

import (
	"crypto/sha256"
	"dds-backend/database"
	"dds-backend/models"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

const ( // token generator constants
	ExpirationDuration           = time.Hour * 12
	TokenGenerationPayloadLength = 10
)

const ( // Claim types, greater values include previous claims
	Worker  = 1
	Manager = 5
	Admin   = 10
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
func Authorize(username, passwordHash string) (string, error) {
	user := models.User{Username: username, Password: passwordHash}
	auth := Auth{Username: username}
	tx := database.DB.Begin() // begin transaction

	var count int
	if tx.Model(&models.User{}).Where(&user).Count(&count); count <= 0 { // check if credentials are valid
		return "", errors.New("credentials invalid")
	}
	now := time.Now()
	if tx.Model(&Auth{}).Where(&auth).Count(&count); count <= 0 { // auth record not found
		auth.Claim = user.Claim
		auth.Expiration = now.Add(ExpirationDuration)
		auth.Token = GenerateNewToken(username)
		tx.Create(&auth)
	} else { // auth record found
		tx.Model(&Auth{}).Where(&auth).First(&auth)
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

type AuthenticationCondition func(auth *Auth) error

func HasEqualOrHigherClaim(claim int) AuthenticationCondition {
	return func(auth *Auth) error {
		if auth.Claim < claim {
			return errors.New("insufficient user rights")
		}
		return nil
	}
}

func HasSameUsername(username string) AuthenticationCondition {
	return func(auth *Auth) error {
		if auth.Username != username {
			return errors.New("access to another user data is forbidden")
		}
		return nil
	}
}

// accept - token claim
// find auth record and check it's validity:
// compare requested claim with available claim
// check for expiration, in case of error return `TokenExpirationError`
// prolong expiration on successful validation
func Authenticate(token string, condition AuthenticationCondition) error {
	auth := Auth{Token: token}
	var count int
	if database.DB.Model(&Auth{}).Where(&auth).Count(&count); count <= 0 { // auth record found
		return errors.New("token not found")
	}
	database.DB.Model(&Auth{}).Where(&auth).First(&auth)
	if time.Now().After(auth.Expiration) {
		return &TokenExpirationError{err: "token has expired"}
	}
	if err := condition(&auth); err != nil {
		return err
	}
	if !time.Now().After(auth.Expiration) {
		auth.Expiration = time.Now().Add(ExpirationDuration)
		if err := database.DB.Save(&auth).Error; err != nil {
			return err
		}
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

func fetchToken(c *gin.Context) (string, error) {
	authLine := c.Request.Header["Authorization"]
	if len(authLine) > 0 {
		return authLine[0], nil
	} else {
		return "", errors.New("authorization missing")
	}
}

func checkAuthConditional(c *gin.Context, condition AuthenticationCondition) error {
	token, err := fetchToken(c) // get authorization token from request header
	if err != nil {
		return err
	}
	err = Authenticate(token, condition) // check if token is valid
	if err != nil {
		return err
	}
	return nil
}
