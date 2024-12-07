package validator

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"reflect"
)

func validateUUID(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Type() == reflect.TypeOf(uuid.UUID{}) {
		return field.Interface().(uuid.UUID) != uuid.Nil
	}
	return false
}
