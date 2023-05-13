package functest

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestMain(m *testing.M) {
// 	// call flag.Parse() here if TestMain uses flags
// 	os.Exit(m.Run())
// }

var amqpDSN = os.Getenv("TESTS_AMQP_DSN")
var postgresDSN = os.Getenv("TESTS_POSTGRES_DSN")

func init() {
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@localhost:5672/"
	}
	if postgresDSN == "" {
		postgresDSN = "host=localhost port=5432 user=calendar password=password dbname=calendar sslmode=disable"
	}
}

func TestTest(t *testing.T) {
	fmt.Println("Hello, AMQP", amqpDSN)
	fmt.Println("Hello, POSTGRES", postgresDSN)
	t.Run("hello world", func(t *testing.T) {
		require.Equal(t, true, false)
	})
}
