package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConnectionString(t *testing.T) {
	t.Run("invalid args", func(t *testing.T) {
		str, err := parseConnectionString([]string{})
		require.Equal(t, "", str)
		require.ErrorIs(t, ErrInvalidArgs, err)
	})
}
