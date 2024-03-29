package common

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

func StringClaim(claim int) string {
	switch claim {
	case Worker:
		return "worker"
	case Manager:
		return "manager"
	case Admin:
		return "admin"
	default:
		return "unknown"
	}
}

// accept - username password
// check credentials, if not valid return error
// if already authorized return existing token, prolong expiration
// if authorization token expired return new token, delete existing token
// if no authorization existed, return new token
func Authorize(username, passwordHash string) (models.Auth, error) {
	user := models.User{Username: username, Password: passwordHash}
	auth := models.Auth{Username: username}
	tx := database.DB.Begin() // begin transaction

	if tx.Model(&models.User{}).Where(&user).First(&user).RecordNotFound() { // check if credentials are valid
		return models.Auth{}, errors.New("credentials invalid")
	}
	now := time.Now().UTC()
	if tx.Model(&models.Auth{}).Where(&auth).First(&auth).RecordNotFound() { // auth record not found
		auth.Claim = user.Claim
		auth.GameType = user.GameType
		auth.Expiration = now.Add(ExpirationDuration)
		auth.Token = GenerateNewToken(username)
		tx.Create(&auth)
	} else { // auth record found
		auth.Token = GenerateNewToken(username, auth.Token) // holder of previous token loses access
		auth.Expiration = now.Add(ExpirationDuration)
		tx.Save(&auth)
	}
	if tx.Error != nil {
		tx.Rollback()
		return models.Auth{}, tx.Error
	}
	return auth, tx.Commit().Error
}

type TokenExpirationError struct {
	err string
}

func (e *TokenExpirationError) Error() string {
	return e.err
}

type AuthenticationCondition func(auth *models.Auth) error

func HasEqualOrHigherClaim(claim int) AuthenticationCondition {
	return func(auth *models.Auth) error {
		if auth.Claim < claim {
			return errors.New("insufficient user rights")
		}
		return nil
	}
}

func HasSameUsername(username string) AuthenticationCondition {
	return func(auth *models.Auth) error {
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
func Authenticate(token string, conditions ...AuthenticationCondition) (*models.Auth, error) {
	auth := models.Auth{Token: token}
	res := database.DB.Model(&models.Auth{}).Where(&auth).First(&auth)
	if res.RecordNotFound() {
		return nil, errors.New("token not found")
	} else if res.Error != nil {
		return nil, res.Error
	}
	if time.Now().UTC().After(auth.Expiration) { // check if token has not expired
		return nil, &TokenExpirationError{err: "token has expired"}
	}
	for _, condition := range conditions { // check if all conditions hold
		if err := condition(&auth); err != nil {
			return nil, err
		}
	}
	auth.Expiration = time.Now().UTC().Add(ExpirationDuration) // prolong expiry date
	if err := database.DB.Model(&models.Auth{}).Save(&auth).Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

func InvalidateToken(token string) error {
	auth := models.Auth{Token: token}
	res := database.DB.Model(&models.Auth{}).Where(&auth).First(&auth)
	if res.RecordNotFound() {
		return errors.New("token not found")
	} else if res.Error != nil {
		return res.Error
	}
	auth.Expiration = time.Now().UTC()
	if err := database.DB.Model(&models.Auth{}).Save(&auth).Error; err != nil {
		return err
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
	rand.Seed(time.Now().UTC().UnixNano())
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
	if authLine, err := c.Cookie("dds-auth-token"); err != nil {
		return "", err
	} else {
		return authLine, nil
	}
}

func CheckAuthConditional(c *gin.Context, conditions ...AuthenticationCondition) (*models.Auth, error) {
	token, err := fetchToken(c) // get authorization token from request header
	if err != nil {
		return nil, err
	}
	auth, err := Authenticate(token, conditions...) // check if token is valid
	if err != nil {
		CleanDDSCookies(c)
	}
	return auth, err
}

const (
	COOKIE_NAME_GAMETYPE = "dds-auth-gametype"
	COOKIE_NAME_CLAIM    = "dds-auth-claim"
	COOKIE_NAME_TOKEN    = "dds-auth-token"
)

func SetDDSCookies(c *gin.Context, auth models.Auth) {
	c.SetCookie(COOKIE_NAME_GAMETYPE, auth.GameType, 60*60*12, "/", "", false, false)
	c.SetCookie(COOKIE_NAME_CLAIM, StringClaim(auth.Claim), 60*60*12, "/", "", false, false)
	c.SetCookie(COOKIE_NAME_TOKEN, auth.Token, 60*60*12, "/", "", false, false)
}

func CleanDDSCookies(c *gin.Context) {
	c.SetCookie(COOKIE_NAME_GAMETYPE, "", -1, "/", "", false, false)
	c.SetCookie(COOKIE_NAME_CLAIM, "", -1, "/", "", false, false)
	c.SetCookie(COOKIE_NAME_TOKEN, "", -1, "/", "", false, false)
}
