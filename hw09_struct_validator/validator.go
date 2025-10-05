package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNotStruct     = errors.New("expected struct")
	ErrInvalidTag    = errors.New("invalid validate tag")
	ErrInvalidRegexp = errors.New("invalid regexp")

	ErrLen    = errors.New("invalid length")
	ErrMin    = errors.New("too small")
	ErrMax    = errors.New("too large")
	ErrIn     = errors.New("not in allowed set")
	ErrRegexp = errors.New("regexp mismatch")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, e := range v {
		sb.WriteString(fmt.Sprintf("%s: %s\n", e.Field, e.Err.Error()))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errs ValidationErrors
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fv := val.Field(i)
		if !fv.CanInterface() {
			continue
		}

		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			parts := strings.SplitN(rule, ":", 2)
			name := parts[0]
			arg := ""
			if len(parts) == 2 {
				arg = parts[1]
			}

			if err := applyRule(fv, field.Name, name, arg); err != nil {
				var vErrs ValidationErrors
				if errors.As(err, &vErrs) {
					errs = append(errs, vErrs...)
				} else {
					errs = append(errs, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func applyRule(fv reflect.Value, fieldName, name, arg string) error {
	switch fv.Kind() {
	case reflect.String:
		return validateString(fv.String(), name, arg)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return validateInt(int(fv.Int()), name, arg)

	case reflect.Slice:
		return validateSlice(fieldName, fv, name, arg)

	case reflect.Bool,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Struct,
		reflect.UnsafePointer,
		reflect.Invalid:
		return nil
	}

	return nil
}

func validateString(val string, name, arg string) error {
	switch name {
	case "len":
		n, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidTag, arg)
		}
		if len(val) != n {
			return fmt.Errorf("%w: expected %d, got %d", ErrLen, n, len(val))
		}

	case "regexp":
		re, err := regexp.Compile(arg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidRegexp, arg)
		}
		if !re.MatchString(val) {
			return fmt.Errorf("%w: %s", ErrRegexp, arg)
		}

	case "in":
		opts := strings.Split(arg, ",")
		for _, opt := range opts {
			if val == opt {
				return nil
			}
		}
		return fmt.Errorf("%w: must be one of %v", ErrIn, opts)
	}
	return nil
}

func validateInt(val int, name, arg string) error {
	switch name {
	case "min":
		n, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidTag, arg)
		}
		if val < n {
			return fmt.Errorf("%w: must be >= %d", ErrMin, n)
		}

	case "max":
		n, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidTag, arg)
		}
		if val > n {
			return fmt.Errorf("%w: must be <= %d", ErrMax, n)
		}

	case "in":
		opts := strings.Split(arg, ",")
		for _, opt := range opts {
			n, err := strconv.Atoi(opt)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidTag, opt)
			}
			if val == n {
				return nil
			}
		}
		return fmt.Errorf("%w: must be one of %v", ErrIn, opts)
	}
	return nil
}

func validateSlice(fieldName string, fv reflect.Value, name, arg string) error {
	var errs ValidationErrors
	for i := 0; i < fv.Len(); i++ {
		elem := fv.Index(i)
		if err := applyRule(elem, fmt.Sprintf("%s[%d]", fieldName, i), name, arg); err != nil {
			errs = append(errs, ValidationError{
				Field: fmt.Sprintf("%s[%d]", fieldName, i),
				Err:   err,
			})
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
