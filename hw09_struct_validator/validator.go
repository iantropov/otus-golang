package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

var (
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrUnsupportedValidation = errors.New("unsupported validation")
	ErrIncorrectValidation   = errors.New("incorrect validation")
)

func Validate(v interface{}) (err error) {
	validationErrors := make(ValidationErrors, 0)
	fields := reflect.VisibleFields(reflect.TypeOf(v))
	for _, field := range fields {
		validationErrors, err = validateField(validationErrors, field, v)
		if err != nil {
			return err
		}
	}

	return validationErrors
}

func validateField(validationErrors ValidationErrors, field reflect.StructField, v interface{}) (_ ValidationErrors, err error) {
	validations, ok := field.Tag.Lookup("validate")
	if !ok {
		return validationErrors, nil
	}

	val := reflect.ValueOf(v)
	fieldValue := val.FieldByName(field.Name)
	if fieldValue.Kind() == reflect.Slice {
		for i := 0; i < fieldValue.Len(); i++ {
			sliceValue := fieldValue.Index(i)
			validationErrors, err = validateValueByValidations(validationErrors, sliceValue, validations)
			if err != nil {
				return validationErrors, err
			}
		}
	} else {
		validationErrors, err = validateValueByValidations(validationErrors, fieldValue, validations)
		if err != nil {
			return validationErrors, err
		}
	}

	return validationErrors, nil
}

func validateValueByValidations(validationErrors ValidationErrors, value reflect.Value, rawValidations string) (ValidationErrors, error) {
	validations := strings.Split(rawValidations, "|")
	for _, validation := range validations {
		validationError, ok, err := validateValueByValidation(value, validation)
		if err != nil {
			return validationErrors, err
		}
		if !ok {
			validationErrors = append(validationErrors, validationError)
		}
	}
	return validationErrors, nil
}

func validateValueByValidation(value reflect.Value, validation string) (_ ValidationError, res bool, err error) {
	var validationError ValidationError
	switch {
	case value.CanInt():
		validationError, res, err = validateInt(value.Int(), validation)
	case value.Kind() == reflect.String:
		validationError, res, err = validateString(value.String(), validation)
	default:
		return validationError, false, ErrUnsupportedType
	}

	return
}

func validateInt(val int64, validation string) (ValidationError, bool, error) {
	var validationError ValidationError
	switch {
	case strings.HasPrefix(validation, "min"):
		minValue, err := parseIntValidation1(validation)
		if err != nil {
			return validationError, true, err
		}
		if val >= minValue {
			return validationError, true, nil
		}
		validationError =
	}
	return validationError, true, nil
}

func validateString(val string, validation string) (ValidationError, bool, error) {
	var validationError ValidationError
	return validationError, true, nil
}

func parseIntValidation1(validation string) (int64, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return 0, ErrIncorrectValidation
	}
	value, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parseIntValidation1: %w", err)
	}
	return value, nil
}

func parseIntValidation2(validation string) (int64, int64, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return 0, 0, ErrIncorrectValidation
	}
	intParts := strings.SplitN(parts[1], 2)
	if len(intParts) != 2 {
		return 0, 0, ErrIncorrectValidation
	}
	value1, err := strconv.ParseInt(intParts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parseIntValidation2: %w", err)
	}
	value2, err := strconv.ParseInt(intParts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parseIntValidation2: %w", err)
	}
	return value1, value2, nil
}
