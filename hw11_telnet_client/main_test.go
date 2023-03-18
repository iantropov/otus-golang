package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConnectionString(t *testing.T) {
	// t.Run("empty args", func(t *testing.T) {
	// 	str, err := parseConnectionString([]string{})
	// 	require.Equal(t, "", str)
	// 	require.ErrorIs(t, ErrInvalidArgs, err)
	// })

	// t.Run("missed host", func(t *testing.T) {
	// 	str, err := parseConnectionString([]string{"", "--timeout 5s", "4242"})
	// 	require.Equal(t, "", str)
	// 	require.ErrorIs(t, ErrInvalidArgs, err)
	// })

	// t.Run("missed port", func(t *testing.T) {
	// 	str, err := parseConnectionString([]string{"", "localhost", "--timeout 5s"})
	// 	require.Equal(t, "", str)
	// 	require.ErrorIs(t, ErrInvalidArgs, err)
	// })

	// t.Run("valid args with timeout", func(t *testing.T) {
	// 	str, err := parseConnectionString([]string{"", "localhost", "4242", "--timeout 5s"})
	// 	require.Equal(t, "localhost:4242", str)
	// 	require.NoError(t, err)
	// })

	// t.Run("valid args without timeout", func(t *testing.T) {
	// 	str, err := parseConnectionString([]string{"", "localhost", "4242"})
	// 	require.Equal(t, "localhost:4242", str)
	// 	require.NoError(t, err)
	})
}
