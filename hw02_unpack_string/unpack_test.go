package hw02unpackstring_test

import (
	"errors"
	"testing"

	hw02unpackstring "github.com/iantropov/otus-golang/hw02_unpack_string"
	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "漢字3", expected: "漢字字字"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := hw02unpackstring.Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "4", "45", "aaa10b", `\`, `qwe\`, `qwe\455`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := hw02unpackstring.Unpack(tc)
			require.Truef(t, errors.Is(err, hw02unpackstring.ErrInvalidString), "actual error %q", err)
		})
	}
}

func BenchmarkStaticStringUnpack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hw02unpackstring.Unpack("QA7w2tovAPA4wj9z3Ad9Z6c8TQSY3DTa2h8eFBЯHDiafXOBnBETyK3SoV1MtjCNy0a2h5E1v6Fk8VxzpLAf9yr7jxI6oCB3vdiHqLQA7w2tovAPA4wj9z3Ad9Z6c8TQSY3DTa2h8eFBЯHDiafXOBnBETyK3SoV1MtjCNy0a2h5E1v6Fk8VxzpLAf9yr7jxI6oCB3vdiHqLQA7w2tovAPA4wj9z3Ad9Z6c8TQSY3DTa2h8eFBЯHDiafXOBnBETyK3SoV1MtjCNy0a2h5E1v6Fk8VxzpLAf9yr7jxI6oCB3vdiHqLQA7w2tovAPA4wj9z3Ad9Z6c8TQSY3DTa2h8eFBЯHDiafXOBnBETyK3SoV1MtjCNy0a2h5E1v6Fk8VxzpLAf9yr7jxI6oCB3vdiHqL")
	}
}
