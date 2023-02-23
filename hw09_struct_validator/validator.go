package hw09structvalidator

import (
	"errors"
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
		fieldValidationErrors, err := validateField(field, v)
		if err != nil {
			return err
		}
		validationErrors = append(validationErrors, fieldValidationErrors...)
	}

	return validationErrors
}

func validateField(field reflect.StructField, v interface{}) (validationErrors ValidationErrors, err error) {
	validations, ok := field.Tag.Lookup("validate")
	if !ok {
		return validationErrors, nil
	}

	val := reflect.ValueOf(v)
	fieldValue := val.FieldByName(field.Name)
	errors := make([]error, 0)
	if fieldValue.Kind() == reflect.Slice {
		for i := 0; i < fieldValue.Len(); i++ {
			sliceValue := fieldValue.Index(i)
			nextErrors, err := validateValueByValidations(sliceValue, validations)
			if err != nil {
				return validationErrors, err
			}
			errors = append(errors, nextErrors...)
		}
	} else {
		errors, err = validateValueByValidations(fieldValue, validations)
		if err != nil {
			return validationErrors, err
		}
	}

	return convertToValidationErrors(errors, field), nil
}

func validateValueByValidations(value reflect.Value, rawValidations string) ([]error, error) {
	validationErrors := make([]error, 0)
	validations := strings.Split(rawValidations, "|")
	for _, validation := range validations {
		validationError, err := validateValueByValidation(value, validation)
		if err != nil {
			return nil, err
		}
		if validationError != nil {
			validationErrors = append(validationErrors, validationError)
		}
	}
	return validationErrors, nil
}

func validateValueByValidation(value reflect.Value, validation string) (error, error) {
	switch {
	case value.CanInt():
		return validateInt(value.Int(), validation)
	case value.Kind() == reflect.String:
		return validateString(value.String(), validation)
	default:
		return nil, ErrUnsupportedType
	}
}

func validateInt(val int64, validation string) (error, error) {
	switch {
	case strings.HasPrefix(validation, "min:"):
		minValue, err := parseIntValidation1(validation)
		if err != nil {
			return nil, err
		}
		if val >= minValue {
			return nil, nil
		}
		return ErrInvalidValue, nil
	case strings.HasPrefix(validation, "max:"):
		maxValue, err := parseIntValidation1(validation)
		if err != nil {
			return nil, err
		}
		if val <= maxValue {
			return nil, nil
		}
		return ErrInvalidValue, nil
	case strings.HasPrefix(validation, "in:"):
		intValues, err := parseIntValidationN(validation)
		if err != nil {
			return nil, err
		}
		for _, intValue := range intValues {
			if intValue == val {
				return nil, nil
			}
		}
		return ErrInvalidValue, nil
	default:
		return nil, ErrUnsupportedValidation
	}
}

func validateString(val string, validation string) (error, error) {
	switch {
	case strings.HasPrefix(validation, "len:"):
		lenValue, err := parseIntValidation1(validation)
		if err != nil {
			return nil, err
		}
		if len(val) == int(lenValue) {
			return nil, nil
		}
		return ErrInvalidValue, nil
	case strings.HasPrefix(validation, "regexp:"):
		regexpValue, err := parseStringValidation1(validation)
		if err != nil {
			return nil, err
		}
		matched, err := regexp.MatchString(regexpValue, val)
		if err != nil {
			return nil, ErrIncorrectValidation
		}
		if matched {
			return nil, nil
		}
		return ErrInvalidValue, nil
	case strings.HasPrefix(validation, "in:"):
		stringValues, err := parseStringValidationN(validation)
		if err != nil {
			return nil, err
		}
		for _, stringValue := range stringValues {
			if stringValue == val {
				return nil, nil
			}
		}
		return ErrInvalidValue, nil
	default:
		return nil, ErrUnsupportedValidation
	}
}

func parseIntValidation1(validation string) (int64, error) {
	stringValue, err := parseStringValidation1(validation)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return 0, ErrIncorrectValidation
	}
	return value, nil
}

func parseIntValidationN(validation string) ([]int64, error) {
	stringValues, err := parseStringValidationN(validation)
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

func parseStringValidation1(validation string) (string, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return "", ErrIncorrectValidation
	}
	return parts[1], nil
}

func parseStringValidationN(validation string) ([]string, error) {
	parts := strings.SplitN(validation, ":", 2)
	if len(parts) != 2 {
		return nil, ErrIncorrectValidation
	}
	stringParts := strings.Split(parts[1], ",")
	return stringParts, nil
}

func convertToValidationErrors(errors []error, field reflect.StructField) ValidationErrors {
	validationErrors := make(ValidationErrors, 0, len(errors))
	for _, err := range errors {
		validationErrors = append(validationErrors, ValidationError{
			Field: field.Name,
			Err:   err,
		})
	}
	return validationErrors
}
