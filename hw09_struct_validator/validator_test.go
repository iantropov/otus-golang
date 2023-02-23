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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

func TestValidateWithInvalidInput(t *testing.T) {
	tests := []struct {
		in         interface{}
		errMessage string
	}{
		{
			in:         Response{Code: 123, Body: "asd"},
			errMessage: "Code: invalid value",
		},
		{
			in:         App{Version: "123"},
			errMessage: "Version: invalid value",
		},
		{
			in: User{ID: "123", Name: "name", Age: 10, Email: "asd", Role: "asd", Phones: []string{"asd", "zxc"}},
			errMessage: "ID: invalid value, " +
				"Age: invalid value, " +
				"Email: invalid value, " +
				"Role: invalid value, " +
				"Phones: invalid value",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.EqualError(t, err, tt.errMessage)
		})
	}
}

func TestValidateWithValidInput(t *testing.T) {
	tests := []interface{}{
		Response{Code: 200, Body: "asd"},
		Token{Header: nil, Payload: nil, Signature: nil},
		App{Version: "12345"},
		User{
			ID:     "123456789x123456789x123456789x123456",
			Name:   "name",
			Age:    20,
			Email:  "asd@asd.asd",
			Role:   "admin",
			Phones: []string{"123456789x1", "123456789x2"},
			meta:   nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)
			require.Equal(t, err, nil)
		})
	}
}
