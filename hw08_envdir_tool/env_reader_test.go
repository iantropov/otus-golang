package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		_, err := ReadDir("asd")
		require.NotEmpty(t, err)
	})

	t.Run("empty dir", func(t *testing.T) {
		err := os.Mkdir(".empty-dir", os.ModePerm)
		require.Empty(t, err)

		env, err := ReadDir(".empty-dir")
		require.Empty(t, err)
		require.Len(t, env, 0)

		err = os.Remove(".empty-dir")
		require.Empty(t, err)
	})

	t.Run("non-empty dir", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.Empty(t, err)
		require.Len(t, env, 5)
	})
}
