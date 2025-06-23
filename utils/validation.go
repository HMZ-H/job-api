package utils

import (
	"regexp"
	"unicode"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("containsuppercase", containsUppercase)
	validate.RegisterValidation("containslowercase", containsLowercase)
	validate.RegisterValidation("containsdigit", containsDigit)
	validate.RegisterValidation("containsspecial", containsSpecial)
	validate.RegisterValidation("alpha", isAlpha)
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func containsUppercase(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func containsLowercase(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func containsDigit(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func containsSpecial(fl validator.FieldLevel) bool {
	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	return specialChars.MatchString(fl.Field().String())
}

func isAlpha(fl validator.FieldLevel) bool {
	alphaRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return alphaRegex.MatchString(fl.Field().String())
}
