package memorystorage

import (
	"fmt"
	"sync"
)

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func New() *Storage {
	fmt.Println("Started in-memory storage!")
	return &Storage{}
}

// TODO
