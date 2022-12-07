package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type Mode int

const (
	ValueMode Mode = iota
	QuantityMode
	ErrorMode
)

type State struct {
	mode        Mode
	currentRune rune
	result      strings.Builder
}

func Unpack(value string) (string, error) {
	state := emptyState()
	for _, r := range value {
		error := state.processRune(r)
		if error != nil {
			return "", error
		}
	}
	return state.getResult(), nil
}

func emptyState() *State {
	return &State{
		mode: ValueMode,
	}
}

func (state *State) processRune(r rune) error {
	if state.mode == ValueMode {
		if isValidValue(r) {
			state.currentRune = r
			state.mode = QuantityMode
		} else {
			return ErrInvalidString
		}
	} else {
		if isValidQuantity(r) {
			state.expandResult()
			state.mode = ValueMode
		} else {
			return ErrInvalidString
		}
	}
	return nil
}

func (state *State) expandResult() {
	state.result.WriteString(strings.Repeat(string(state.currentRune), 5))
}

func (state *State) getResult() string {
	return state.result.String()
}

func isValidValue(value rune) bool {
	return true
}

func isValidQuantity(value rune) bool {
	return true
}
