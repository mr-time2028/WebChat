package validators

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode"
)

const (
	PasswordMinLength = 8
	PasswordMaxLength = 30

	PasswordMinLengthErrMsg = `this field must be a minimum length of %d characters`
	PasswordMaxLengthErrMsg = `this field must be a maximum length of %d characters`
	PasswordUppercaseErrMsg = `this field must contain at least one uppercase letter`
	PasswordLowercaseErrMsg = `this field must contain at least one lowercase letter`
	PasswordDigitErrMsg     = `this field must contain at least one digit`
	PasswordsNotMatchErrMsg = `passwords should be match`
)

// PasswordMatchesHashValidation check if two password are match, password from db (registered) and password from client
func (v *Validation) PasswordMatchesHashValidation(hashedDBPassword, ClientPassword string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedDBPassword), []byte(ClientPassword))
	if err != nil {
		v.Errors.Add("password", err.Error())
		if err == bcrypt.ErrMismatchedHashAndPassword {
			v.Errors.Code = http.StatusUnauthorized
		} else {
			v.Errors.Code = http.StatusInternalServerError
		}
	}
}

// PasswordsMatchesCharactersValidation check if user password are matches
func (v *Validation) PasswordsMatchesCharactersValidation(password1, password2 string) {
	if password1 != password2 {
		v.Errors.Add("password", PasswordsNotMatchErrMsg)
	}
}

// PasswordCharacterValidation check if password characters are valid (min length, max length and etc.)
func (v *Validation) PasswordCharacterValidation(password string) {
	hasUppercase := false
	hasLowercase := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		} else if unicode.IsLower(char) {
			hasLowercase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if len(password) < PasswordMinLength {
		v.Errors.Add("password", fmt.Sprintf(PasswordMinLengthErrMsg, PasswordMinLength))
	}
	if len(password) > PasswordMaxLength {
		v.Errors.Add("password", fmt.Sprintf(PasswordMaxLengthErrMsg, PasswordMaxLength))
	}
	if !hasUppercase {
		v.Errors.Add("password", PasswordUppercaseErrMsg)
	}
	if !hasLowercase {
		v.Errors.Add("password", PasswordLowercaseErrMsg)
	}
	if !hasDigit {
		v.Errors.Add("password", PasswordDigitErrMsg)
	}
}

// UserPasswordValidation check user password
func (v *Validation) UserPasswordValidation(password1, password2 string) {
	v.PasswordsMatchesCharactersValidation(password1, password2)
	if v.Valid() {
		v.PasswordCharacterValidation(password1)
	}
}
