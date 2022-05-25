package model

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/sethvargo/go-password/password"
)

const staffPasswordLen = 12
const staffPasswordNumDigits = 4
const staffPasswordNumSymbols = 2
const passwordMinLength = 6
const passwordMaxLength = 32

var validNum = regexp.MustCompile(`[0-9]{1}`)
var validLowerCase = regexp.MustCompile(`[a-z]{1}`)
var validUpperCase = regexp.MustCompile(`[A-Z]{1}`)
var validSymbol = regexp.MustCompile(`[~!@#$%^&\*()_\+\-=|:<>\?,\.]{1}`)

func GenerateStaffPassword() (string, error) {
	pass, err := password.Generate(
		staffPasswordLen,
		staffPasswordNumDigits,
		staffPasswordNumSymbols,
		false,
		false)
	if err != nil {
		return "", fmt.Errorf("fail to generate staff password - %v", err)
	}

	return pass, nil
}

func IsStrongPassword(ps string) bool {
	if utf8.RuneCountInString(ps) < passwordMinLength {
		return false
	} else if utf8.RuneCountInString(ps) > passwordMaxLength {
		return false
	}

	if !validNum.MatchString(ps) {
		return false
	}
	if !validLowerCase.MatchString(ps) {
		return false
	}
	if !validUpperCase.MatchString(ps) {
		return false
	}
	if !validSymbol.MatchString(ps) {
		return false
	}
	return true
}
