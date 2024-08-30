package validations

import "regexp"

func IsValidEmail(email string) bool {
	//validar email
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

	re := regexp.MustCompile(regex)

	return re.MatchString(email)
}
