package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestMain(m *testing.M) {
// 	// call flag.Parse() here if TestMain uses flags
// 	os.Exit(m.Run())
// }

func TestTest(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		require.Equal(t, true, false)
	})
}
