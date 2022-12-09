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
	EscapeMode
	QuantityMode
	ErrorMode
)

type parser struct {
	mode        Mode
	currentRune rune
	result      strings.Builder
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
		if parser.mode == EscapeMode {
			parser.processRuneAsValue(r)
		} else if parser.mode != ValueMode {
			parser.mode = ErrorMode
		} else {
			parser.processRuneAsQuantity(r)
		}
	} else if r == '/' {
		parser.mode = EscapeMode
	} else {
		parser.processRuneAsValue(r)
	}
}

func (parser *parser) processRuneAsValue(r rune) {
	if parser.mode == ValueMode {
		parser.result.WriteRune(parser.currentRune)
	}
	parser.currentRune = r
	parser.mode = ValueMode
}

func (parser *parser) processRuneAsQuantity(r rune) {
	parser.expandResult(r)
	parser.mode = QuantityMode
}

func (parser *parser) expandResult(q rune) {
	quantity := int(q - '0')
	parser.result.WriteString(strings.Repeat(string(parser.currentRune), quantity))
}

func (parser *parser) getResult() string {
	if parser.mode == ValueMode {
		parser.result.WriteRune(parser.currentRune)
	}
	return parser.result.String()
}
