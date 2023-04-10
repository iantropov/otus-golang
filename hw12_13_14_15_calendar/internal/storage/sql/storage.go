package sqlstorage

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pgxPool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{}
}
