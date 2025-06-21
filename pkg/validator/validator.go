package validator

import (
	"github.com/go-playground/validator/v10"
)

var ValidBrazilianStates = map[string]bool{
	"AC": true, "AL": true, "AP": true, "AM": true, "BA": true,
	"CE": true, "DF": true, "ES": true, "GO": true, "MA": true,
	"MT": true, "MS": true, "MG": true, "PA": true, "PB": true,
	"PR": true, "PE": true, "PI": true, "RJ": true, "RN": true,
	"RS": true, "RO": true, "RR": true, "SC": true, "SP": true,
	"SE": true, "TO": true,
}

func ValidateBrazilianState(fl validator.FieldLevel) bool {
	state := fl.Field().String()
	return ValidBrazilianStates[state]
}

func SetupCustomValidators(v *validator.Validate) {
	v.RegisterValidation("brazilian_state", ValidateBrazilianState)
}
