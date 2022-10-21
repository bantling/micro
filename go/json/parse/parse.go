package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"

	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/json"
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
// as only the caller knows what to do on EOF.
func parseValue(it *iter.Iter[token]) json.Value {
	tok := it.Must()

	switch tok.typ {
	case tOBrace:
		it.Unread(tok)
		return parseObject(it)
	case tOBracket:
		it.Unread(tok)
		// Collect all array elements into a slice, and return it
		return json.FromSliceOfValue(iter.ReduceToSlice(parseArray(it)).Must())
	case tString:
		return json.FromString(tok.value)
	case tNumber:
		return json.FromNumberString(json.NumberString(tok.value))
	case tBoolean:
		return json.FromBool(tok.value == "true")
	case tNull:
		return json.NullValue
	}

	return json.Value{}
}

// parseObject parses a JSON object, making potetially recursive calls to parseValue for each key value.
// The iter must provide the opening brace.
// Panics unless the correct lexical elements occur in the correct order.
// Panics if a duplicate key occurs.
func parseObject(it *iter.Iter[token]) json.Value {
	// Discard opening brace
	it.Must()

	var (
		haveIt       bool
		key          token
		colon        token
		value        json.Value
		invalidValue json.Value
		commaBrace   token
		object       map[string]json.Value = map[string]json.Value{}
	)

	// Read as many key/value pairs as are provided
	for {
		// Must have a closing brace or string key after opening brace
		if key, haveIt = it.NextValue(); !haveIt {
			panic(errObjectRequiresKeyOrBrace)
		}

		// If closing brace, valid empty object
		if key.typ == tCBrace {
			return json.FromMapOfValue(object)
		}

		// If not closing brace, must be string key
		if key.typ != tString {
			panic(errObjectRequiresKeyOrBrace)
		}

		// Panic if key is a duplicate
		if _, haveIt := object[key.value]; haveIt {
			panic(fmt.Errorf(errObjectDuplicateKeyMsg, key.value))
		}

		// Expect colon separator
		if colon, haveIt = it.NextValue(); (!haveIt) || (colon.typ != tColon) {
			panic(fmt.Errorf(errObjectKeyRequiresColonMsg, key.value))
		}

		// Expect value for key, and map it
		if value = parseValue(it); value == invalidValue {
			panic(fmt.Errorf(errObjectKeyRequiresValueMsg, key.value))
		}
		object[key.value] = value

		// Expect a comma or closing brace
		if commaBrace, haveIt = it.NextValue(); !haveIt {
			panic(fmt.Errorf(errObjectKeyValueRequiresCommaOrBraceMsg, key.value))
		}

		if commaBrace.typ == tComma {
			// Read next key/value pair
			continue
		}

		if commaBrace.typ == tCBrace {
			// Return completed object
			break
		}

		panic(fmt.Errorf(errObjectKeyValueRequiresCommaOrBraceMsg, key.value))
	}

	// Return value built from map
	return json.FromMapOfValue(object)
}

// parseArray parses a JSON array, making potetially recursive calls to parseValue for each element value.
// The iter must provide the opening brace.
//
// Returns an iter of each array element, in case this array is the document, so that the caller can process each element
// as they are parsed. If the array is empty, the first call to the returned iter.Next will be false.
//
// Panics unless the correct lexical elements occur in the correct order, which could happen after returning some correctly
// formed elements.
func parseArray(it *iter.Iter[token]) *iter.Iter[json.Value] {
	var (
		first        = true
		tok          token
		haveIt       bool
		value        json.Value
		invalidValue json.Value
		commaBracket token
	)

	// Discard opening bracket
	it.Must()

	return iter.NewIter(func() (json.Value, bool) {
		if first {
			// Must have a closing bracket or value after opening bracket
			if tok, haveIt = it.NextValue(); !haveIt {
				panic(errArrayRequiresValueOrBracket)
			}

			// If closing bracket, valid empty array
			if tok.typ == tCBracket {
				return invalidValue, false
			}

			// If not closing bracket, must be a value
			it.Unread(tok)
			if value = parseValue(it); value == invalidValue {
				// If it's not a closing bracket or value, then it isn't valid
				panic(errArrayRequiresValueOrBracket)
			}

			// Return first value
			first = false
			return value, true
		}

		// At least one element has already been returned, expect a comma or closing bracket
		if commaBracket, haveIt = it.NextValue(); !haveIt {
			panic(errArrayRequiresCommaOrBracket)
		}

		if commaBracket.typ == tComma {
			// Expect value for element, and return it
			if value = parseValue(it); value == invalidValue {
				panic(errArrayRequiresValue)
			}

			return value, true
		}

		if commaBracket.typ == tCBracket {
			// Indicate end of array
			return invalidValue, false
		}

		panic(errArrayRequiresCommaOrBracket)
	})
}

// Iterate parses a JSON document into an iter[json.Value].
// If the document is an object, the object is parsed completely before returning it as an iter of one element.
// If the document is an array, the array top level elements are parsed as the returnd iter is iterated.
//
// Useful for cases like writing JSON data to a database, where the JSON input could contain a large number of records,
// and it is preferable to store each record one at a time, or perhaps in batches of some fixed maximum size.
//
// Panics if the input is not an object or array (including empty/whitespace only input), or if lexical elements do not
// occur in the correct order (eg unbalanced brackets).
func Iterate(src io.Reader) *iter.Iter[json.Value] {
	// First lexical element must be a { or [
	var (
		// Reader > iter[rune] > iter[token]
		it               = lexer(iter.OfReaderAsRunes(src))
		firstTok, haveIt = it.NextValue()
	)

	// Die if empty
	if !haveIt {
		panic(errEmptyDocument)
	}

	// If object, return an iter of one element that is parsed right now
	if firstTok.typ == tOBrace {
		it.Unread(firstTok)
		return iter.OfOne(parseObject(it))
	}

	// If array, return iter of array elements, which are parsed later as the iter is iterated
	if firstTok.typ == tOBracket {
		it.Unread(firstTok)
		return parseArray(it)
	}

	// Die if some other token exists that is not a brace or bracket
	panic(errObjectOrArrayRequired)
}

// Parse parses a JSON document fully before returning it, unlike Iterate which provides an iter.
// Useful for cases like a configuration file, where you need the whole document, and it is easier to not have to iterate.
// The top level object or array is provided as a Value.
func Parse(src io.Reader) json.Value {
	// First lexical element must be a { or [
	var (
		// Reader > iter[rune] > iter[token]
		it               = lexer(iter.OfReaderAsRunes(src))
		firstTok, haveIt = it.NextValue()
	)

	// Die if empty
	if !haveIt {
		panic(errEmptyDocument)
	}

	// If object or array, return the fully parsed object as a Value
	if (firstTok.typ == tOBrace) || (firstTok.typ == tOBracket) {
		it.Unread(firstTok)
		return parseValue(it)
	}

	// Die if some other token exists that is not a brace or bracket
	panic(errObjectOrArrayRequired)
}
