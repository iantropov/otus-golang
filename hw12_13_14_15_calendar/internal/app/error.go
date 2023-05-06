package app

import "fmt"

type InternalError struct {
	Err error
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal server error: %v", e.Err)
}
