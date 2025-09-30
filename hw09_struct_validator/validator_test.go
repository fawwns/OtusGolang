package hw09structvalidator

import (
	"errors"
	"strings"
	"testing"
)

type UserRole string

type User struct {
	ID     string `validate:"len:36"`
	Name   string
	Age    int      `validate:"min:18|max:50"`
	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole `validate:"in:admin,stuff"`
	Phones []string `validate:"len:11"`
}

type App struct {
	Version string `validate:"len:5"`
}

type Response struct {
	Code int `validate:"in:200,404,500"`
	Body string
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		wantErr bool
		checks  []func(ValidationErrors, *testing.T)
	}{
		{
			name: "valid user",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Alice",
				Age:    25,
				Email:  "alice@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			wantErr: false,
		},
		{
			name: "invalid user - multiple errors",
			in: User{
				ID:     "shortid",
				Name:   "Bob",
				Age:    15,
				Email:  "invalid-email",
				Role:   "guest",
				Phones: []string{"123", "456"},
			},
			wantErr: true,
			checks: []func(ValidationErrors, *testing.T){
				func(errs ValidationErrors, t *testing.T) {
					found := false
					for _, e := range errs {
						if e.Field == "ID" && errors.Is(e.Err, ErrLen) {
							if !strings.Contains(e.Err.Error(), "36") {
								t.Errorf("ID error should mention 36, got %v", e.Err)
							}
							found = true
						}
					}
					if !found {
						t.Error("expected ErrLen for ID")
					}
				},
				func(errs ValidationErrors, t *testing.T) {
					found := false
					for _, e := range errs {
						if e.Field == "Age" && errors.Is(e.Err, ErrMin) {
							found = true
						}
					}
					if !found {
						t.Error("expected ErrMin for Age")
					}
				},
				func(errs ValidationErrors, t *testing.T) {
					found := false
					for _, e := range errs {
						if e.Field == "Email" && errors.Is(e.Err, ErrRegexp) {
							found = true
						}
					}
					if !found {
						t.Error("expected ErrRegexp for Email")
					}
				},
				func(errs ValidationErrors, t *testing.T) {
					found := false
					for _, e := range errs {
						if e.Field == "Role" && errors.Is(e.Err, ErrIn) {
							found = true
						}
					}
					if !found {
						t.Error("expected ErrIn for Role")
					}
				},
				func(errs ValidationErrors, t *testing.T) {
					count := 0
					for _, e := range errs {
						if strings.HasPrefix(e.Field, "Phones[") && errors.Is(e.Err, ErrLen) {
							count++
						}
					}
					if count != 2 {
						t.Errorf("expected 2 ErrLen errors for Phones, got %d", count)
					}
				},
			},
		},
		{
			name: "valid app",
			in: App{
				Version: "v1.00",
			},
			wantErr: false,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 302,
				Body: "Found",
			},
			wantErr: true,
			checks: []func(ValidationErrors, *testing.T){
				func(errs ValidationErrors, t *testing.T) {
					found := false
					for _, e := range errs {
						if e.Field == "Code" && errors.Is(e.Err, ErrIn) {
							if !strings.Contains(e.Err.Error(), "200") {
								t.Errorf("Code error should mention allowed values, got %v", e.Err)
							}
							found = true
						}
					}
					if !found {
						t.Error("expected ErrIn for Code")
					}
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var vErrs ValidationErrors
				if !errors.As(err, &vErrs) {
					t.Fatalf("expected ValidationErrors, got %T", err)
				}
				for _, check := range tt.checks {
					check(vErrs, t)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}
