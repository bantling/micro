package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"unicode/utf16"

	"github.com/bantling/micro/iter"
)

// tokenType is an enum of lexical token types
type tokenType uint

// tokenType enum constants
const (
	tOBracket tokenType = iota
	tCBracket
	tOBrace
	tCBrace
	tComma
	tColon
	tString
	tNumber
	tBoolean
	tNull
	tEOF
)

// error constants
var (
	errIncompleteStringMsg            = "Incomplete string %s: a string must be terminated by a \""
	errControlCharInStringMsg         = "The ascii control character 0x%02x is not valid in a string"
	errIncompleteStringEscapeMsg      = "Incomplete string escape in %s"
	errIllegalStringEscapeMsg         = "Illegal string escape %s"
	errOneSurrogateEscapeMsg          = "The surrogate string escape %s must be followed by another surrogate escape to form valid UTF-16"
	errSurrogateNonSurrogateEscapeMsg = "The surrogate string escape %s cannot be followed by the non-surrogate escape %s"
	errSurrogateDecodeEscapeMsg       = "The surrogate string escape pair %s is not a valid UTF-16 surrogate pair"
	errInvalidNumberMsg               = "Invalid number %s: a number must satisfy the regex -?[0-9]+([.][0-9]+)?([eE][0-9]+)?"
	errInvalidBooleanNullMsg          = "Invalid sequence %s: an array, object, string, number, boolean, or null was expected"
	errInvalidCharMsg                 = "Invalid character %s: an array, object, string, number, boolean, or null was expected"
)

// token is a single lexical token
type token struct {
	typ   tokenType
	value string
}

// token constants
var (
	tokOBracket = token{tOBracket, "["}
	tokCBracket = token{tCBracket, "]"}
	tokOBrace   = token{tOBrace, "{"}
	tokCBrace   = token{tCBrace, "}"}
	tokComma    = token{tComma, ","}
	tokColon    = token{tColon, ":"}
	tokTrue     = token{tBoolean, "true"}
	tokFalse    = token{tBoolean, "false"}
	tokNull     = token{tNull, "null"}
)

// lexString lexes a string token, where the iter begins by returning the opening quote character.
// the returned token does not contain the quotes, only the runes between them.
func lexString(it iter.Iter[rune]) (token, error) {
	var (
		str          = []rune{}
		val          rune
		zv           = token{}
		err          error
		readUTF16Hex = func(temp *[]rune) (rune, error) {
			var hexVal rune

			for i := 0; i < 4; i++ {
				if val, err = it.Next(); err != nil {
					if err == iter.EOI {
						return 0, fmt.Errorf(errIncompleteStringEscapeMsg, string(*temp))
					}
					return 0, err
				}

				*temp = append(*temp, val)
				if (val >= '0') && (val <= '9') {
					hexVal = hexVal*16 + val - '0'
				} else if (val >= 'A') && (val <= 'F') {
					hexVal = hexVal*16 + val - 'A' + 10
				} else if (val >= 'a') && (val <= 'f') {
					hexVal = hexVal*16 + val - 'a' + 10
				} else {
					return 0, fmt.Errorf(errIllegalStringEscapeMsg, string(*temp))
				}
			}

			return hexVal, nil
		}
	)

	// Discard opening quote
	it.Next()

	// Loop until we find unescaped closing double quote
	for {
		if val, err = it.Next(); err != nil {
			if err == iter.EOI {
				return zv, fmt.Errorf(errIncompleteStringMsg, "\""+string(str))
			}
			return zv, err
		}

		// ASCII control characters cannot be in a string
		if val < ' ' {
			return zv, fmt.Errorf(errControlCharInStringMsg, val)
		}

		// If we read a backslash, it must be followed by ", \, /, b, f, n, r, t, or u and 4 hex chars
		if val == '\\' {
			if val, err = it.Next(); err != nil {
				if err == iter.EOI {
					return zv, fmt.Errorf(errIncompleteStringEscapeMsg, string(append(str, '\\')))
				}
				return zv, err
			}

			switch val {
			case '"': // Escaped double quote
				// Special case - we cannot just set val to a ", that would cause later code to think it is the end of the string
				str = append(str, '"')
				continue
			case '\\': // Escaped backslash
				val = '\\'
			case '/': // Escaped forward slash
				val = '/'
			case 'b': // Escaped backspace
				val = '\b'
			case 'f': // Escaped form feed
				val = '\f'
			case 'n': // Escaped line feed
				val = '\n'
			case 'r': // Escaped carriage return
				val = '\r'
			case 't': // Escaped tab
				val = '\t'
			case 'u': // Escaped unicode char in UTF-16
				// Must be followed by 4 hex chars
				var (
					temp   = []rune{'\\', 'u'}
					hexVal rune
				)
				if hexVal, err = readUTF16Hex(&temp); err != nil {
					return zv, err
				}

				// Is this escape part of a UTF-16 surrogate pair?
				if utf16.IsSurrogate(hexVal) {
					temp2 := []rune{}

					// Then we expect another \u sequence to immediately follow that is alao part of a surrogate pair
					if val, err = it.Next(); err != nil {
						if err == iter.EOI {
							return zv, fmt.Errorf(errOneSurrogateEscapeMsg, string(temp))
						}
						return zv, err
					}

					if val != '\\' {
						return zv, fmt.Errorf(errOneSurrogateEscapeMsg, string(temp))
					}
					temp2 = append(temp2, val)

					if val, err = it.Next(); err != nil {
						if err == iter.EOI {
							return zv, fmt.Errorf(errOneSurrogateEscapeMsg, string(temp))
						}
						return zv, err
					}

					if val != 'u' {
						return zv, fmt.Errorf(errOneSurrogateEscapeMsg, string(temp))
					}
					temp2 = append(temp2, val)

					var hexVal2 rune
					if hexVal2, err = readUTF16Hex(&temp2); err != nil {
						return zv, err
					}
					if !utf16.IsSurrogate(hexVal2) {
						return zv, fmt.Errorf(errSurrogateNonSurrogateEscapeMsg, string(temp), string(temp2))
					}

					// DecodeRune returns 0xFFFD if the pair is not a valid UTF-16 pair.
					// Note that UTF-16 can be big or little endian, and so can the processor.
					// Go decodes in the order presented in RFC8259 regardless of processor, where U+1D11E is encoded as \uD834\uDD1E.
					if hexVal = utf16.DecodeRune(hexVal, hexVal2); hexVal == 0xFFFD {
						return zv, fmt.Errorf(errSurrogateDecodeEscapeMsg, string(temp)+string(temp2))
					}
				}

				val = hexVal
			default:
				return zv, fmt.Errorf(errIllegalStringEscapeMsg, string('\\')+string(val))
			}
		}

		// val is the next char to add, whether a plain or escaped char, or escaped surrogate pair
		// Discard closing quote
		if val == '"' {
			break
		}
		str = append(str, val)
	}

	return token{tString, string(str)}, nil
}

// lexNumber lexes a number token, where the iter begins by returning the first rune, which may be a leading - or a digit.
// The returned token contains every character up to last digit read.
func lexNumber(it iter.Iter[rune]) (token, error) {
	// There must be first char, it may be a - or digit
	val, err := it.Next()
	var (
		str       []rune
		zv        token
		haveDigit = val != '-'
		die       = func() (token, error) {
			if (err != nil) && (err != iter.EOI) {
				// A problem
				return zv, err
			}

			// EOI or other lexing problem
			return token{}, fmt.Errorf(errInvalidNumberMsg, string(str))
		}
		tok = func() (token, error) {
			if err == nil {
				it.Unread(val)
			}
			return token{tNumber, string(str)}, nil
		}
	)
	str = append(str, val)

	// Read digits until dot or e or non-dot non-e non-digit
	for {
		if val, err = it.Next(); err != nil {
			if err == iter.EOI {
				if !haveDigit {
					return die()
				}

				return tok()
			}
			// A problem
			return zv, err
		}

		if (val >= '0') && (val <= '9') {
			haveDigit = true
			str = append(str, val)
		} else if val == '.' {
			str = append(str, val)
			if haveDigit {
				break
			}

			return die()
		} else if (val == 'e') || (val == 'E') {
			str = append(str, val)
			if haveDigit {
				break
			}

			return die()
		} else {
			// Ok to have just optional - and some digits
			if haveDigit {
				return tok()
			}

			return die()
		}
	}

	// Last char may be dot or e or E. If it's a dot, read at least one digit until e or non-e non-digit
	if val == '.' {
		// Must have at least one digit after dot
		haveDigit = false
		for {
			if val, err = it.Next(); err != nil {
				if err == iter.EOI {
					if !haveDigit {
						return die()
					}

					// Ok to just have optional -, some digits, a dot, and some digits.
					return tok()
				}

				// A problem
				return zv, err
			}

			if (val >= '0') && (val <= '9') {
				haveDigit = true
				str = append(str, val)
			} else if (val == 'e') || (val == 'E') {
				str = append(str, val)
				if haveDigit {
					break
				}

				return die()
			} else {
				if haveDigit {
					return tok()
				}

				return die()
			}
		}
	}

	// Last char must be e or E. Read optional + or - and at least one digit until a non-digit
	haveDigit = false
	if val, err = it.Next(); err != nil {
		return die()
	}

	if (val == '+') || (val == '-') {
		// Append sign
		str = append(str, val)
	} else {
		// Unread what is hopefully a digit
		it.Unread(val)
	}

	for {
		if val, err = it.Next(); err != nil {
			if (err == iter.EOI) && haveDigit {
				// At least one digit after e or E and optional sign
				return tok()
			}
			// Either EOI with no digit after e or E, or a problem
			return die()
		}

		if (val >= '0') && (val <= '9') {
			haveDigit = true
			str = append(str, val)
		} else {
			if haveDigit {
				break
			}

			return die()
		}
	}

	return tok()
}

// lexBooleanNull lexes a boolean or null token, where the iter begins by returning the first rune.
func lexBooleanNull(it iter.Iter[rune]) (token, error) {
	var (
		str []rune
		val rune
		err error
		zv  token
	)

	// Just read chars until a non-lowercase letter
	for {
		if val, err = it.Next(); err != nil {
			if err == iter.EOI {
				break
			}
			// A problem
			return zv, err
		}

		if (val >= 'a') && (val <= 'z') {
			str = append(str, val)
		} else {
			// First char for next token
			it.Unread(val)
			break
		}
	}

	cstr := string(str)
	switch cstr {
	case "true":
		return tokTrue, nil
	case "false":
		return tokFalse, nil
	case "null":
		return tokNull, nil
	}

	return zv, fmt.Errorf(errInvalidBooleanNullMsg, string(str))
}

// lex lexes the next token, which must be [, ], {, }, comma, :, string, number, boolean, null, or eof.
// Skip whitespace chars.
func lex(it iter.Iter[rune]) (token, error) {
	// Handle eof
	var zv token
	val, err := it.Next()
	if err != nil {
		// EOI or problem, doesn't matter which
		return zv, err
	}

	it.Unread(val)

	// Skip ws
	for {
		if val, err = it.Next(); err != nil {
			// EOI or problem, doesn't matter which
			return zv, err
		}

		if !((val == ' ') || (val == '\n') || (val == '\r') || (val == '\t')) {
			break
		}
	}

	// First non-ws char
	switch {
	case val == '[':
		return tokOBracket, nil
	case val == ']':
		return tokCBracket, nil
	case val == '{':
		return tokOBrace, nil
	case val == '}':
		return tokCBrace, nil
	case val == ',':
		return tokComma, nil
	case val == ':':
		return tokColon, nil
	case val == '"':
		it.Unread(val)
		return lexString(it)
	case (val == '-') || ((val >= '0') && (val <= '9')):
		it.Unread(val)
		return lexNumber(it)
	case (val == 't') || (val == 'f') || (val == 'n'):
		it.Unread(val)
		return lexBooleanNull(it)
	}

	// Anything except the above is an illegal character
	return zv, fmt.Errorf(errInvalidCharMsg, string(val))
}

// lexer uses lex and converts an iter[rune] into an iter[token]
func lexer(it iter.Iter[rune]) iter.Iter[token] {
	return iter.NewIter(func() (token, error) {
		return lex(it)
	})
}
