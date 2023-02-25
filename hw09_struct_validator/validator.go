package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type ValidationResult string

const (
	Valid              ValidationResult = "ok"
	Unknown            ValidationResult = "unknown"
	InvalidWithMinimum ValidationResult = "less than the allowed minimum"
	InvalidWithMaximum ValidationResult = "more than the allowed maximum"
	InvalidWithSet     ValidationResult = "isn't element of the set"
	InvalidWithRegex   ValidationResult = "doesn't match the regex"
	InvalidWithLen     ValidationResult = "invalid length"
)

func (validationErrors ValidationErrors) Error() string {
	builder := strings.Builder{}
	for i, err := range validationErrors {
		builder.WriteString(err.Field)
		builder.WriteString(": ")
		builder.WriteString(err.Err.Error())
		if i != len(validationErrors)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}

var (
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrUnsupportedValidation = errors.New("unsupported validation")
	ErrIncorrectValidation   = errors.New("incorrect validation")

	ErrInvalidValue = errors.New("invalid value")
)

func Validate(v interface{}) (err error) {
	validationErrors := make(ValidationErrors, 0)
	fields := reflect.VisibleFields(reflect.TypeOf(v))
	for _, field := range fields {
		res, err := validateField(field, v)
		if err != nil {
			return err
		}
		if res != Valid {
			validationErrors = append(validationErrors, ValidationError{
				Field: field.Name,
				Err:   fmt.Errorf(string(res)),
			})
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func validateField(field reflect.StructField, v interface{}) (ValidationResult, error) {
	validations, ok := field.Tag.Lookup("validate")
	if !ok {
		return Valid, nil
	}

	val := reflect.ValueOf(v)
	fieldValue := val.FieldByName(field.Name)
	if fieldValue.Kind() == reflect.Slice {
		for i := 0; i < fieldValue.Len(); i++ {
			sliceValue := fieldValue.Index(i)
			res, err := validateValueByValidations(sliceValue, validations)
			if err != nil || res != Valid {
				return res, err
			}
		}
	} else {
		res, err := validateValueByValidations(fieldValue, validations)
		return res, err
	}

	return Valid, nil
}

func validateValueByValidations(value reflect.Value, rawValidations string) (ValidationResult, error) {
	validations := strings.Split(rawValidations, "|")
	for _, validation := range validations {
		res, err := validateValueByValidation(value, validation)
		if err != nil {
			return Unknown, err
		}
		if res != Valid {
			return res, nil
		}
	}
	return Valid, nil
}

func validateValueByValidation(value reflect.Value, validation string) (ValidationResult, error) {
	switch {
	case value.CanInt():
		return validateInt(value.Int(), validation)
	case value.Kind() == reflect.String:
		return validateString(value.String(), validation)
	default:
		return Valid, ErrUnsupportedType
	}
}

func validateInt(val int64, validation string) (ValidationResult, error) {
	switch {
	case strings.HasPrefix(validation, "min:"):
		minValue, err := parseIntValidation(validation)
		if err != nil {
			return Unknown, err
		}
		if val >= minValue {
			return Valid, nil
		}
		return InvalidWithMinimum, nil
	case strings.HasPrefix(validation, "max:"):
		maxValue, err := parseIntValidation(validation)
		if err != nil {
			return Unknown, err
		}
		if val <= maxValue {
			return Valid, nil
		}
		return InvalidWithMaximum, nil
	case strings.HasPrefix(validation, "in:"):
		intValues, err := parseIntValidations(validation)
		if err != nil {
			return Unknown, err
		}
		for _, intValue := range intValues {
			if intValue == val {
				return Valid, nil
			}
		}
		return InvalidWithSet, nil
	default:
		return Unknown, ErrUnsupportedValidation
	}
}

func validateString(val string, validation string) (ValidationResult, error) {
	switch {
	case strings.HasPrefix(validation, "len:"):
		lenValue, err := parseIntValidation(validation)
		if err != nil {
			return Unknown, err
		}
		if len(val) == int(lenValue) {
			return Valid, nil
		}
		return InvalidWithLen, nil
	case strings.HasPrefix(validation, "regexp:"):
		regexpValue, err := parseStringValidation(validation)
		if err != nil {
			return Unknown, err
		}
		matched, err := regexp.MatchString(regexpValue, val)
		if err != nil {
			return Unknown, ErrIncorrectValidation
		}
		if matched {
			return Valid, nil
		}
		return InvalidWithRegex, nil
	case strings.HasPrefix(validation, "in:"):
		stringValues, err := parseStringValidations(validation)
		if err != nil {
			return Unknown, err
		}
		for _, stringValue := range stringValues {
			if stringValue == val {
				return Valid, nil
			}
		}
		return InvalidWithSet, nil
	default:
		return Unknown, ErrUnsupportedValidation
	}
}

func parseIntValidation(validation string) (int64, error) {
	stringValue, err := parseStringValidation(validation)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return 0, ErrIncorrectValidation
	}
	return value, nil
}

func parseIntValidations(validation string) ([]int64, error) {
	stringValues, err := parseStringValidations(validation)
	if err != nil {
		return nil, err
	}
	intValues := make([]int64, 0, len(stringValues))
	for _, stringValue := range stringValues {
		intValue, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return nil, ErrIncorrectValidation
		}
		intValues = append(intValues, intValue)
	}
	return intValues, nil
}

func parseStringValidation(validation string) (string, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return "", ErrIncorrectValidation
	}
	return parts[1], nil
}

func parseStringValidations(validation string) ([]string, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return nil, ErrIncorrectValidation
	}
	stringParts := strings.Split(parts[1], ",")
	return stringParts, nil
}
