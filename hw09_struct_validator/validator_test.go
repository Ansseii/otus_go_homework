package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage `validate:"len:12"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          1,
			expectedErr: ErrorNotStruct,
		},
		{
			in: User{
				ID:     "12345",
				Name:   "Admin",
				Age:    14,
				Email:  "wrong_gmail.com",
				Role:   "unknown",
				Phones: []string{"1234567890", "12345678900"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrorInvalidLength},
				ValidationError{Field: "Age", Err: ErrorNumberIsLessThenMinimum},
				ValidationError{Field: "Email", Err: ErrorInvalidStringPattern},
				ValidationError{Field: "Role", Err: ErrorElementIsNotInSet},
				ValidationError{Field: "Phones", Err: ErrorInvalidLength},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Admin",
				Age:    26,
				Email:  "test@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678900"},
				meta:   nil,
			},
		},
		{
			in: App{
				Version: "1.1",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrorInvalidLength},
			},
		},
		{in: Token{
			Header:    []byte("test"),
			Payload:   []byte("test"),
			Signature: []byte("test"),
		}},
		{
			in: Response{
				Code: 100,
				Body: "qwerty",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrorElementIsNotInSet},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
