package models

import (
	"github.com/jinzhu/gorm"
	"regexp"
	"sync"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Name     string
	Surname  string
	Phone    string `gorm:"unique"`
	Address  string
	Claim    int    `gorm:"not null;default:1"`
	GameType string `gorm:"not null"`
}

func (u *User) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["username"] = u.Username
	result["name"] = u.Name
	result["surname"] = u.Surname
	result["phone"] = u.Phone
	result["address"] = u.Address
	result["gametype"] = u.GameType
	return result
}

type regexChecker struct {
	usernamePattern *regexp.Regexp
	phonePattern    *regexp.Regexp
	emailPattern    *regexp.Regexp
}

var regexCheckerInstance *regexChecker
var once sync.Once

func regexCheckerGetInstance() *regexChecker {
	once.Do(func() {
		regexCheckerInstance = &regexChecker{}
		var err error
		regexCheckerInstance.usernamePattern, err = regexp.Compile(`[a-zA-Z0-9_]{3,15}`)
		if err != nil {
			panic(err)
		}
		regexCheckerInstance.phonePattern, err = regexp.Compile(`(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}`)
		if err != nil {
			panic(err)
		}
		regexCheckerInstance.emailPattern, err = regexp.Compile(`[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*`)
		if err != nil { // FIXME: add missing ` in regex
			panic(err)
		}
	})
	return regexCheckerInstance
}

func (u *User) IsValid() (bool, string) {
	if !regexCheckerGetInstance().usernamePattern.MatchString(u.Username) {
		return false, "username must consist of letters, numbers or underscores of length 3-15"
	}
	// TODO: add other checks and output message in compound form
	// FIXME: regex not working correctly !!!
	return true, ""
}
