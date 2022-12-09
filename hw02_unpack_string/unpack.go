package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type Mode int

const (
	StartMode Mode = iota
	ValueMode
	QuantityMode
	ErrorMode
)

type parser struct {
	mode         Mode
	previousRune rune
	result       strings.Builder
}

func Unpack(value string) (string, error) {
	parser := &parser{}
	for _, r := range value {
		parser.processRune(r)
		if parser.mode == ErrorMode {
			return "", ErrInvalidString
		}
	}
	return parser.getResult(), nil
}

func (parser *parser) processRune(r rune) {
	if unicode.IsDigit(r) {
		if parser.mode != ValueMode {
			parser.mode = ErrorMode
		} else {
			parser.expandResult(r)
			parser.mode = QuantityMode
		}
	} else {
		if parser.mode == ValueMode {
			parser.result.WriteRune(parser.previousRune)
		}
		parser.previousRune = r
		parser.mode = ValueMode
	}
}

func (parser *parser) expandResult(q rune) {
	quantity := int(q - '0')
	parser.result.WriteString(strings.Repeat(string(parser.previousRune), quantity))
}

func (parser *parser) getResult() string {
	if parser.mode == ValueMode {
		parser.result.WriteRune(parser.previousRune)
	}
	return parser.result.String()
}
