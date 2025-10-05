package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			input: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Age:    25,
				Email:  "test@mail.com",
				Role:   "admin",
				Phones: []string{"89001234567"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid user - multiple errors",
			input: User{
				ID:     "short-id",
				Age:    15,
				Email:  "bad_email",
				Role:   "guest",
				Phones: []string{"123", "456"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrLen},
				{Field: "Age", Err: ErrMin},
				{Field: "Email", Err: ErrRegexp},
				{Field: "Role", Err: ErrIn},
				{Field: "Phones[0]", Err: ErrLen},
				{Field: "Phones[1]", Err: ErrLen},
			},
		},
		{
			name:        "valid app",
			input:       App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			name:  "invalid response code",
			input: Response{Code: 123},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrIn},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // copy loop variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			checkValidation(t, tt.input, tt.expectedErr)
		})
	}
}

func checkValidation(t *testing.T, input interface{}, expected error) {
	t.Helper()
	err := Validate(input)

	if expected == nil {
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		return
	}

	var gotErrs ValidationErrors
	if !errors.As(err, &gotErrs) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}

	var wantErrs ValidationErrors
	if !errors.As(expected, &wantErrs) {
		t.Fatalf("expected expectedErr to be ValidationErrors, got %T", expected)
	}

	if len(gotErrs) != len(wantErrs) {
		t.Fatalf("expected %d errors, got %d", len(wantErrs), len(gotErrs))
	}

	for i := range wantErrs {
		if gotErrs[i].Field != wantErrs[i].Field {
			t.Errorf("expected field %q, got %q", wantErrs[i].Field, gotErrs[i].Field)
		}
		if !errors.Is(gotErrs[i].Err, wantErrs[i].Err) {
			t.Errorf("field %q: expected error %v, got %v", gotErrs[i].Field, wantErrs[i].Err, gotErrs[i].Err)
		}
	}
}
