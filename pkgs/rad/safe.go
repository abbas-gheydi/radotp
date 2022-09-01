package rad

import (
	"log"
	"regexp"
	"unicode"
)

var usernameRegexp = regexp.MustCompile(`^[0-9A-Za-z_.@\-]{1,30}$`)

func IsOtpCodeSafe(input string) bool {
	if len([]rune(input)) != 6 {
		return false
	}
	for _, digit := range input {
		if !unicode.IsDigit(digit) {
			return false
		}
	}
	return true
}

func isSafeInput(input string) bool {

	safe := usernameRegexp.MatchString(input)
	if !safe {
		log.Println("unexpected input string", input)
	}
	return safe

}
