package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTag    = "validate"
	rulesSeparator = "|"
	rulesDelimiter = ":"
)

var (
	ErrorNotStruct                 = errors.New("not a struct received")
	ErrorUnsupportedType           = errors.New("unsupported type")
	ErrorUnsupportedRule           = errors.New("unsupported rule")
	ErrorInvalidLength             = errors.New("length is invalid")
	ErrorElementIsNotInSet         = errors.New("element is not in the set")
	ErrorInvalidStringPattern      = errors.New("string pattern is invalid")
	ErrorNumberIsLessThenMinimum   = errors.New("number is less then specified minimum")
	ErrorNumberIsBiggerThenMaximum = errors.New("number is bigger then specified minimum")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}

	for _, e := range v {
		builder.WriteString(fmt.Sprintf(`Field "%s": %s`, e.Field, e.Err.Error()))
	}

	return builder.String()
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Struct {
		return validateStruct(v)
	}

	return ErrorNotStruct
}

func validateStruct(s interface{}) error {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	var validationErrors ValidationErrors
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		validationErrors = append(validationErrors, validateByType(field, value)...)
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func validateByType(field reflect.StructField, value reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors
	rules := parseRules(field)
	if len(rules) == 0 {
		return validationErrors
	}
	switch value.Type().Kind() { //nolint:exhaustive
	case reflect.String:
		validationErrors = append(validationErrors, validateString(field, value, rules)...)
	case reflect.Int:
		validationErrors = append(validationErrors, validateInt(field, value, rules)...)
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			validationErrors = append(validationErrors, validateByType(field, value.Index(i))...)
		}
	default:
		validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorUnsupportedType})
	}

	return validationErrors
}

func validateInt(field reflect.StructField, value reflect.Value, rules map[string]string) ValidationErrors {
	var validationErrors ValidationErrors
	curVal := int(value.Int())

	for rule, v := range rules {
		switch rule {
		case "min":
			if min, err := strconv.Atoi(v); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			} else if curVal < min {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorNumberIsLessThenMinimum})
			}
		case "max":
			if max, err := strconv.Atoi(v); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			} else if curVal > max {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorNumberIsBiggerThenMaximum})
			}
		case "in":
			in := strings.Split(v, ",")
			if found, err := containsInt(in, curVal); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			} else if !found {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorElementIsNotInSet})
			}
		default:
			validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorUnsupportedRule})
		}
	}

	return validationErrors
}

func validateString(field reflect.StructField, value reflect.Value, rules map[string]string) ValidationErrors {
	var validationErrors ValidationErrors
	curVal := value.String()
	for rule, v := range rules {
		switch rule {
		case "len":
			l, err := strconv.Atoi(v)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			} else if len(curVal) < l {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorInvalidLength})
			}
		case "regexp":
			matched, err := regexp.MatchString(v, curVal)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			} else if !matched {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorInvalidStringPattern})
			}
		case "in":
			in := strings.Split(v, ",")
			if !containsStr(in, curVal) {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorElementIsNotInSet})
			}
		default:
			validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: ErrorUnsupportedRule})
		}
	}

	return validationErrors
}

func parseRules(field reflect.StructField) map[string]string {
	rules := make(map[string]string)
	if rulesStr, ok := field.Tag.Lookup(validateTag); ok {
		sRules := strings.Split(rulesStr, rulesSeparator)
		for _, el := range sRules {
			rule := strings.Split(el, rulesDelimiter)
			name, value := rule[0], rule[1]
			rules[name] = value
		}
	}

	return rules
}

func containsStr(arr []string, candidate string) bool {
	for _, el := range arr {
		if el == candidate {
			return true
		}
	}
	return false
}

func containsInt(arr []string, candidate int) (bool, error) {
	for _, el := range arr {
		if n, err := strconv.Atoi(el); err != nil {
			return false, err
		} else if n == candidate {
			return true, nil
		}
	}
	return false, nil
}
