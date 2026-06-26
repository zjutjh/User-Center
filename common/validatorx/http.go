package validatorx

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type HTTPValidator struct{}

func (HTTPValidator) Validate(_ *http.Request, data any) error {
	return validateRequired(reflect.ValueOf(data))
}

func validateRequired(value reflect.Value) error {
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	valueType := value.Type()
	for i := range value.NumField() {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		if field.Anonymous {
			if err := validateRequired(fieldValue); err != nil {
				return err
			}
		}

		if !hasRequiredTag(field.Tag.Get("validate")) {
			continue
		}
		if isZero(fieldValue) {
			return fmt.Errorf("field %q is not set", jsonFieldName(field))
		}
	}

	return nil
}

func hasRequiredTag(tag string) bool {
	for _, part := range strings.Split(tag, ",") {
		if strings.TrimSpace(part) == "required" {
			return true
		}
	}
	return false
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	name, _, _ := strings.Cut(tag, ",")
	if name == "" || name == "-" {
		return field.Name
	}
	return name
}

func isZero(value reflect.Value) bool {
	if value.Kind() == reflect.String {
		return strings.TrimSpace(value.String()) == ""
	}
	return value.IsZero()
}
