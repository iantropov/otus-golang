package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)
	t.Run("invalid cmd", func(t *testing.T) {
		_, err := RunCmd([]string{"asd"}, env)
		require.NotEmpty(t, err)
	})

	t.Run("cmd with non-zero return code", func(t *testing.T) {
		code, err := RunCmd([]string{"./testdata/return.sh", "1"}, env)
		require.Empty(t, err)
		require.Equal(t, 1, code)
	})

	t.Run("md with zero return code", func(t *testing.T) {
		code, err := RunCmd([]string{"./testdata/return.sh", "0"}, env)
		require.Empty(t, err)
		require.Equal(t, 0, code)
	})
}
