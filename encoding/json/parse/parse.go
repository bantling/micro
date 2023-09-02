package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"

	"github.com/bantling/micro/encoding/json"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/iter"
	"github.com/bantling/micro/stream"
)

// Error constants
var (
	errObjectRequiresKeyOrBrace              = fmt.Errorf("A JSON object must have a string key or closing brace after the opening brace")
	errObjectDuplicateKeyMsg                 = "A JSON object cannot have duplicate key %q"
	errObjectKeyRequiresColonMsg             = "The JSON object key %q just be followed by a colon"
	errObjectKeyRequiresValueMsg             = "The JSON object key %q must be have a value that is an object, arrray, string, number, boolean, or null"
	errObjectKeyValueRequiresCommaOrBraceMsg = "The JSON key/value pair %q must be followed by a colon or closing brace"
	errArrayRequiresValue                    = fmt.Errorf("A JSON array element must be an object, array, string, number, boolean, or null")
	errArrayRequiresValueOrBracket           = fmt.Errorf("A JSON array must have an element or closing bracket after the opening bracket")
	errArrayRequiresCommaOrBracket           = fmt.Errorf("A JSON array element must be followed by a comma or closing bracket")
	errEmptyDocument                         = fmt.Errorf("A JSON document cannot be empty")
	errObjectOrArrayRequired                 = fmt.Errorf("A JSON document must begin with a brace or bracket")
)

// parseValue detects the type of value (object, array, string, number,  boolean, or null).
// Objects and arrays are dispatched to parseObject and parseArray, while scalars are returned as an appropriate Value.
// If the next token is not an opening brace, opening bracket, string, number, boolean, or null, a zero value is returned
// so the caller can panic with an appropriate error.
//
// It is assumed that there does exist at least one more token - it is up to the caller to test this before calling,
// as only the caller knows what to do on EOI.
func parseValue(it iter.Iter[token]) (json.Value, error) {
	// Get first token
	tok, err := it.Next()
	var zv json.Value

	switch tok.typ {
	case tOBrace:
		it.Unread(tok)
		return parseObject(it)
	case tOBracket:
		it.Unread(tok)
		// Assume this array is not top level document, collect all array elements into a slice and return it
		var slc []json.Value
		if slc, err = stream.ReduceToSlice(parseArray(it)).Next(); err != nil {
			return zv, err
		}
		return json.FromSliceOfValue(slc), nil
	case tString:
		return json.FromString(tok.value), nil
	case tNumber:
		return json.FromNumber(json.NumberString(tok.value)), nil
	case tBoolean:
		return json.FromBool(tok.value == "true"), nil
	case tNull:
		return json.NullValue, nil
	}

	// Only the caller knows what to do if an invalid token occurs
	return zv, nil
}

// parseObject parses a JSON object, making potetially recursive calls to parseValue for each key value.
// The iter must provide the opening brace.
func parseObject(it iter.Iter[token]) (json.Value, error) {
	// Discard opening brace
	it.Next()

	var (
		key        token
		err        error
		colon      token
		valueTok   token
		value      json.Value
		zv         json.Value
		commaBrace token
		object     = map[string]json.Value{}
	)

	// Read as many key/value pairs as are provided
	for {
		// Must have a closing brace or string key after opening brace
		if key, err = it.Next(); err != nil {
			if err == iter.EOI {
				return zv, errObjectRequiresKeyOrBrace // Case 1
			}
			// A problem
			return zv, err // Case 2
		}

		// If closing brace, valid empty object
		if key.typ == tCBrace {
			return json.FromMapOfValue(object), nil
		}

		// If not closing brace, must be string key
		if key.typ != tString {
			return zv, errObjectRequiresKeyOrBrace // Case 3
		}

		// Panic if key is a duplicate
		if _, haveIt := object[key.value]; haveIt {
			return zv, fmt.Errorf(errObjectDuplicateKeyMsg, key.value) // Case 4
		}

		// Expect colon separator
		if colon, err = it.Next(); err != nil {
			if err == iter.EOI {
				return zv, fmt.Errorf(errObjectKeyRequiresColonMsg, key.value) // Case 5
			}
			// A problem
			return zv, err // Case 6
		}
		if colon.typ != tColon {
			return zv, fmt.Errorf(errObjectKeyRequiresColonMsg, key.value) // Case 7
		}

		// parseValue expects caller to verify there is another token, and provide an appropriate error if not
		if valueTok, err = it.Next(); err != nil {
			if err == iter.EOI {
				return zv, fmt.Errorf(errObjectKeyRequiresValueMsg, key.value) // Case 8
			}
			// A problem
			return zv, err // Case 9
		}

		// Expect value for key, and map it
		it.Unread(valueTok)
		if value, err = parseValue(it); err != nil {
			// A problem
			return zv, err // Case 10
		}
		object[key.value] = value

		// Expect a comma or closing brace
		if commaBrace, err = it.Next(); err != nil {
			if err == iter.EOI {
				return zv, fmt.Errorf(errObjectKeyValueRequiresCommaOrBraceMsg, key.value) // Case 11
			}
			// A problem
			return zv, err // Case 12
		}

		if commaBrace.typ == tComma {
			// Read next key/value pair
			continue
		}

		if commaBrace.typ == tCBrace {
			// Return completed object
			break
		}

		return zv, fmt.Errorf(errObjectKeyValueRequiresCommaOrBraceMsg, key.value) // Case 13
	}

	// Return value built from map
	return json.FromMapOfValue(object), nil
}

// parseArray parses a JSON array, making potetially recursive calls to parseValue for each element value.
// The iter must provide the opening brace.
//
// Returns an iter of each array element, in case this array is the document, so that the caller can process each element
// as they are parsed. If the array is empty, the first call to the returned iter.Next will be false.
func parseArray(it iter.Iter[token]) iter.Iter[json.Value] {
	var (
		first        = true
		tok          token
		err          error
		value        json.Value
		zv           json.Value
		commaBracket token
	)

	return iter.OfIter(func() (json.Value, error) {
		if first {
			// Discard opening bracket
			it.Next()

			// Must have a closing bracket or value after opening bracket
			if tok, err = it.Next(); err != nil {
				if err == iter.EOI {
					return zv, errArrayRequiresValueOrBracket // Case 1
				}
				// A problem
				return zv, err // Case 2
			}

			// If closing bracket, valid empty array
			if tok.typ == tCBracket {
				return zv, iter.EOI // Case 3
			}

			// If not closing bracket, must be a value that begins with the token we just read
			it.Unread(tok)
			if value, err = parseValue(it); err != nil {
				// A problem - cannot be EOI
				return zv, err // Case 4
			}

			// Return first value
			first = false
			return value, nil
		}

		// At least one element has already been returned, expect a comma or closing bracket
		if commaBracket, err = it.Next(); err != nil {
			if err == iter.EOI {
				return zv, errArrayRequiresCommaOrBracket // Case 5
			}
			// A problem
			return zv, err // Case 6
		}

		if commaBracket.typ == tComma {
			// Ensure there is another token that can be read by parseValue
			if tok, err = it.Next(); err != nil {
				if err == iter.EOI {
					return zv, errArrayRequiresValue // Case 7
				}
				// A problem
				return zv, err // Case 8
			}

			// Expect value for element, and return it
			it.Unread(tok)
			if value, err = parseValue(it); err != nil {
				// A problem - cannot be EOI
				return zv, err // Case 9
			}

			return value, nil
		}

		if commaBracket.typ == tCBracket {
			// Indicate end of array
			return zv, iter.EOI // Case 10
		}

		return zv, errArrayRequiresCommaOrBracket // Case 11
	})
}

// Iterate parses a JSON document into an iter[json.Value].
// If the document is an object, the object is parsed completely before returning it as an iter of one element.
// If the document is an array, the array top level elements are parsed as the returnd iter is iterated.
//
// Useful for cases like writing JSON data to a database, where the JSON input could contain a large number of records,
// and it is preferable to store each record one at a time, or perhaps in batches of some fixed maximum size.
func Iterate(src io.Reader) iter.Iter[json.Value] {
	// First lexical element must be a { or [
	var (
		// Reader > iter[rune] > iter[token]
		it            = lexer(iter.OfReaderAsRunes(src))
		firstTok, err = it.Next()
	)

	// Die if empty
	if err != nil {
		if err == iter.EOI {
			return iter.SetError(iter.OfEmpty[json.Value](), errEmptyDocument) // Case 1
		}
		// A problem
		return iter.SetError(iter.OfEmpty[json.Value](), err) // Case 2
	}

	// If object, return an iter of one element that is parsed right now
	if firstTok.typ == tOBrace {
		it.Unread(firstTok)

		var val json.Value
		val, err = parseObject(it)
		if err != nil {
			// Can't be EOI, we know a token existed before the call
			return iter.SetError(iter.OfEmpty[json.Value](), err) // Case 3
		}

		return iter.OfOne(val)
	}

	// If array, return iter of array elements, which are parsed later as the iter is iterated
	if firstTok.typ == tOBracket {
		it.Unread(firstTok)
		return parseArray(it)
	}

	// Die if some other token exists that is not a brace or bracket
	return iter.SetError(iter.OfEmpty[json.Value](), errObjectOrArrayRequired) // Case 4
}

// Parse parses a JSON document fully before returning it, unlike Iterate which provides an iter.
// Useful for cases like a configuration file, where you need the whole document, and it is easier to not have to iterate.
// The top level object or array is provided as a Value.
// If the reader can be parsed into a valid json document the result is (Value, nil), else it is (invalid value, error).
func Parse(src io.Reader) (json.Value, error) {
	// First lexical element must be a { or [
	// Reader > iter[rune] > iter[token]
	var (
		it            = lexer(iter.OfReaderAsRunes(src))
		firstTok, err = it.Next()
		doc           json.Value
		zv            json.Value
	)

	// Die if empty
	if err != nil {
		if err == iter.EOI {
			return zv, errEmptyDocument // Case 1
		}
		// A problem
		return zv, err // Case 2
	}

	// If object or array, return the fully parsed object as a Value
	if (firstTok.typ == tOBrace) || (firstTok.typ == tOBracket) {
		it.Unread(firstTok)
		if doc, err = parseValue(it); err != nil {
			// Can't be EOI, already checked that
			return zv, err // Case 3
		}

		return doc, nil
	}

	// Die if some other token exists that is not a brace or bracket
	return zv, errObjectOrArrayRequired // Case 4
}

// MustParse is a must version of Parse
func MustParse(src io.Reader) json.Value {
	return funcs.MustValue(Parse(src))
}
