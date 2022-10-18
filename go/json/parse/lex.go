package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"unicode/utf16"

	"github.com/bantling/micro/go/iter"
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
	errIncompleteArray                = fmt.Errorf("An array must be terminated by a ]")
	errIncompleteObject               = fmt.Errorf("An object must be terminated by a }")
	errIncompleteStringMsg            = "Incomplete string %s: a string must be terminated by a \""
	errControlCharInStringMsg         = "The ascii control character %x is not valid in a string"
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
	typ tokenType
	val string
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
	tokEOF      = token{tEOF, ""}
)

// lexString lexes a string token, where the iter begins by returning the opening quote character.
// the returned token does not contain the quotes, only the runes between them.
func lexString(it *iter.Iter[rune]) token {
	var (
		str          = []rune{}
		readUTF16Hex = func(temp *[]rune) rune {
			var (
				hexVal rune
				r      rune
			)

			for i := 0; i < 4; i++ {
				if !it.Next() {
					panic(fmt.Errorf(errIllegalStringEscapeMsg, string(*temp)))
				}

				r = it.Value()
				*temp = append(*temp, r)
				if (r >= '0') && (r <= '9') {
					hexVal = hexVal*16 + r - '0'
				} else if (r >= 'A') && (r <= 'F') {
					hexVal = hexVal*16 + r - 'A' + 10
				} else if (r >= 'a') && (r <= 'f') {
					hexVal = hexVal*16 + r - 'a' + 10
				} else {
					panic(fmt.Errorf(errIllegalStringEscapeMsg, string(*temp)))
				}
			}

			return hexVal
		}
	)

	// Discard opening quote
	it.Must()

	// Loop until we find unescaped closing double quote
	for it.Next() {
		r := it.Value()

		// ASCII control characters cannot be in a string
		if r < ' ' {
			panic(fmt.Errorf(errControlCharInStringMsg, r))
		}

		// If we read a backslash, it must be followed by ", \, /, b, f, n, r, t, or u and 4 hex chars
		if r == '\\' {
			if !it.Next() {
				panic(fmt.Errorf(errIncompleteStringEscapeMsg, string(append(str, r))))
			}

			r = it.Value()
			switch r {
			case '"': // Escaped double quote
				// Special case - we cannot just set r to a ", that would cause later code to think it is the end of the string
				str = append(str, '"')
				continue
			case '\\': // Escaped backslash
				r = '\\'
			case '/': // Escaped forward slash
				r = '/'
			case 'b': // Escaped backspace
				r = '\b'
			case 'f': // Escaped form feed
				r = '\f'
			case 'n': // Escaped line feed
				r = '\n'
			case 'r': // Escaped carriage return
				r = '\r'
			case 't': // Escaped tab
				r = '\t'
			case 'u': // Escaped unicode char in UTF-16
				// Must be followed by 4 hex chars
				var (
					temp   = []rune{'\\', 'u'}
					hexVal = readUTF16Hex(&temp)
				)

				// Is this escape a UTF-16 part of a surrogate pair?
				if utf16.IsSurrogate(hexVal) {
					temp2 := []rune{}

					// Then we expect another \u sequence to immediately follow that is alao part of a surrogate pair
					if !it.Next() {
						panic(fmt.Errorf(errOneSurrogateEscapeMsg, string(temp)))
					}
					if r = it.Value(); r != '\\' {
						panic(fmt.Errorf(errOneSurrogateEscapeMsg, string(temp)))
					}
					temp2 = append(temp2, r)

					if !it.Next() {
						panic(fmt.Errorf(errOneSurrogateEscapeMsg, string(temp)))
					}
					if r = it.Value(); r != 'u' {
						panic(fmt.Errorf(errOneSurrogateEscapeMsg, string(temp)))
					}
					temp2 = append(temp2, r)

					hexVal2 := readUTF16Hex(&temp2)
					if !utf16.IsSurrogate(hexVal2) {
						panic(fmt.Errorf(errSurrogateNonSurrogateEscapeMsg, string(temp), string(temp2)))
					}

					// DecodeRune returns 0xFFFD if the pair is not a valid UTF-16 pair.
					// Note that UTF-16 can be big or little endian, and so can the processor.
					// Go decodes in the order presented in RFC8259 regardless of processor, where U+1D11E is encoded as \uD834\uDD1E.
					if hexVal = utf16.DecodeRune(hexVal, hexVal2); hexVal == 0xFFFD {
						panic(fmt.Errorf(errSurrogateDecodeEscapeMsg, string(temp)+string(temp2)))
					}
				}

				r = hexVal
			default:
				panic(fmt.Errorf(errIllegalStringEscapeMsg, string('\\')+string(r)))
			}
		}

		// r is the next chara to add, whether a plain or escaped char, or escaped surrogate pair
		// Discard closing quote
		if r == '"' {
			return token{tString, string(str)}
		}
		str = append(str, r)
	}

	panic(fmt.Errorf(errIncompleteStringMsg, `"`+string(str)))
}

// lexNumber lexes a number token, where the iter begins by returning the first rune, which may be a leading - or a digit.
// The returned token contains every character up to last digit read.
func lexNumber(it *iter.Iter[rune]) token {
	// There must be first char, it may be a - or digit, just add it. That way, loop can read required digits before dot.
	it.Next()
	var (
		str       []rune
		r         = it.Value()
		haveDigit = r != '-'
		die       = func() { panic(fmt.Errorf(errInvalidNumberMsg, string(str))) }
		tok       = func() token { it.Unread(r); return token{tNumber, string(str)} }
	)
	str = append(str, r)

	// Read digits until dot or e or non-dot non-e non-digit
	r = 0
	for it.Next() {
		if r = it.Value(); (r >= '0') && (r <= '9') {
			haveDigit = true
			str = append(str, r)
		} else if r == '.' {
			str = append(str, r)
			if haveDigit {
				break
			}

			die()
		} else if (r == 'e') || (r == 'E') {
			str = append(str, r)
			if haveDigit {
				break
			}

			die()
		} else {
			// Ok to have just optional - and some digits
			if haveDigit {
				return tok()
			}

			die()
		}

		r = 0
	}

	// EOF may have occured
	if r == 0 {
		if !haveDigit {
			die()
		}
		return tok()
	}

	// Last char may be dot or e or E. If it's a dot, read at least one digit until e or non-e non-digit
	if r == '.' {
		// Must have at least one digit after dot
		haveDigit = false
		if !it.Next() {
			die()
		}
		it.Unread(it.Value())

		r = 0
		for it.Next() {
			if r = it.Value(); (r >= '0') && (r <= '9') {
				haveDigit = true
				str = append(str, r)
			} else if (r == 'e') || (r == 'E') {
				str = append(str, r)
				if haveDigit {
					break
				}

				die()
			} else {
				if haveDigit {
					return tok()
				}

				die()
			}

			r = 0
		}
	}

	if r == 0 {
		return tok()
	}

	// Last char must be e or E. Read optional + or - and at least one digit until a non-digit
	haveDigit = false
	if !it.Next() {
		die()
	}

	if r = it.Value(); (r == '+') || (r == '-') {
		// Append sign
		str = append(str, r)
	} else {
		// Unread what is hopefully a digit
		it.Unread(r)
	}

	r = 0
	for it.Next() {
		if r = it.Value(); (r >= '0') && (r <= '9') {
			haveDigit = true
			str = append(str, r)
		} else {
			if haveDigit {
				// Have to break here so we can have return statement at top level of function
				break
			}

			die()
		}

		r = 0
	}

	if (r == 0) && (!haveDigit) {
		die()
	}

	return tok()
}

// lexBooleanNull lexes a boolean or null token, where the iter begins by returning the first rune.
func lexBooleanNull(it *iter.Iter[rune]) token {
	var (
		str []rune
		r   rune
	)

	// Just read chars until a non-lowercase letter
	for it.Next() {
		if r = it.Value(); (r >= 'a') && (r <= 'z') {
			str = append(str, r)
		} else {
			// First char for next token
			it.Unread(r)
			break
		}
	}

	cstr := string(str)
	if cstr == "true" {
		return tokTrue
	} else if cstr == "false" {
		return tokFalse
	} else if cstr == "null" {
		return tokNull
	}

	panic(fmt.Errorf(errInvalidBooleanNullMsg, string(str)))
}

// lex lexes the next token, which must be [, ], {, }, comma, :, string, number, boolean, null, or eof.
// Skip whitespace chars.
func lex(it *iter.Iter[rune]) token {
	// Handle eof
	if !it.Next() {
		return tokEOF
	}
	it.Unread(it.Value())

	// Skip ws
	var r rune
	for it.Next() {
		if r = it.Value(); !((r == ' ') || (r == '\n') || (r == '\r') || (r == '\t')) {
			break
		}
	}

	// First non-ws char
	switch {
	case r == '[':
		return tokOBracket
	case r == ']':
		return tokCBracket
	case r == '{':
		return tokOBrace
	case r == '}':
		return tokCBrace
	case r == ',':
		return tokComma
	case r == ':':
		return tokColon
	case r == '"':
		it.Unread(r)
		return lexString(it)
	case (r == '-') || ((r >= '0') && (r <= '9')):
		it.Unread(r)
		return lexNumber(it)
	case (r == 't') || (r == 'f') || (r == 'n'):
		it.Unread(r)
		return lexBooleanNull(it)
	}

	// Anything except the above is an illegal character
	panic(fmt.Errorf(errInvalidCharMsg, r))
}
