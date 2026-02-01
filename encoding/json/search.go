package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"regexp"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/tuple"
	"github.com/bantling/micro/union"
)

var (
	// regexPathParts is a pre compiled regex for object key and array index path parts
	// the leading dot of an object key and square brackets around an array index are not returned, just the keys and indexes
	regexPathParts = regexp.MustCompile(`(?:\.([^.\[\]]+)|\[([0-9]+)\])`)

	errIllegalPathMsg = "The path %s is not a valid path, it must consist of a series of object keys and indexes, such as .addresses[3].city"
	errNoSuchPathMsg  = "The path %s cannot be found, as %s is not the correct type, or does not contain the index %v"
)

// parsePath parses a path string into a slice of tuple/union structures that represents what to search for.
// The structures contain {path, union of {key, index}}.
// The idea is that for .addresses[3].city, there would be the following tuple/union structures:
// - {t: ".addresses", u: {t: "addresses"}}
// - {t: "[3]", {u: 3}}
// - {t: ".city", {t: "city"}}
//
// An error is returned if the given path does not match the regex for a valid path
func parsePath(p string) (lookups []tuple.Two[string, union.Two[string, int]], err error) {
	// Ensure there are no extra characters in the string before or after the path parts
	// The only way to do this is to replace all matches with the empty string, and verify the result is an empty string
	if len(regexPathParts.ReplaceAllLiteralString(p, "")) > 0 {
		err = fmt.Errorf(errIllegalPathMsg, p)
		return
	}

	// Break up the path into individual object key and array index lookups
	// Result is a [][]string, where the inner []string for a given match is [match, key without ., index without brackets],
	// such as [".address", "addresses", ""] or ["[3]", "", "3"]
	parts := regexPathParts.FindAllStringSubmatch(p, -1)

	if parts == nil {
		// The string path is not recognizable by the regex
		err = fmt.Errorf(errIllegalPathMsg, p)
		return
	}

	// Convert each path part into a tuple of {full path so far, union{string key, int index}}.
	// The full path is for errors, to indicate at which point the failure occurs in a lookup.
	// The key or index is for performing the lookup.
	var (
		path, key, fullPath string
		index               int
	)

	for _, part := range parts {
		path, key = part[0], part[1]
		fullPath += path

		if len(key) > 0 {
			// Object key
			lookups = append(lookups, tuple.Of2(fullPath, union.Of2T[string, int](key)))
		} else {
			// Array index
			if err = conv.To(part[2], &index); err != nil {
				// Must be an index so large it can't be converted to an int
				err = fmt.Errorf(errIllegalPathMsg, path)
				return
			}

			lookups = append(lookups, tuple.Of2(fullPath, union.Of2U[string, int](index)))
		}
	}

	return
}

// lookupsToFunc converts the tuple/union structure of parsePath into a func(Value) (Value, error).
// The function applies the path to the Value to try and find the element.
//
// The function returns an error if:
// - The Value given is not an Object or Array
// - Any portion of the path cannot be found in the Value
func lookupsToFunc(lookups []tuple.Two[string, union.Two[string, int]]) func(Value) (Value, error) {
	return func(in Value) (out Value, err error) {
		var (
			curValue      = in
			keyIndex      union.Two[string, int]
			fullPath, key string
			index         int
			typ           Type
			slc           []Value
		)

		for _, lookup := range lookups {
			fullPath, keyIndex = lookup.Values()

			// Is this path part a key?
			if keyIndex.Which() == union.T {
				key = keyIndex.T()

				// Is this value an Object?
				if typ = curValue.Type(); typ != Object {
					// No, can't apply key to it
					err = fmt.Errorf(errNoSuchPathMsg, fullPath, typ, key)
					return
				}

				// Does the object contain the key?
				if v, hasIt := curValue.AsMap()[key]; hasIt {
					// value of key is next json Value
					curValue = v
				} else {
					// No, can't find key
					err = fmt.Errorf(errNoSuchPathMsg, fullPath, typ, key)
					return
				}
			} else {
				// Which must be union.U
				index = keyIndex.U()

				// Is this value an array?
				if typ = curValue.Type(); typ != Array {
					// No, can't apply it
					err = fmt.Errorf(errNoSuchPathMsg, fullPath, typ, index)
					return
				}

				// Does the array contain the index?
				if slc = curValue.AsSlice(); len(slc) > index {
					// value of index is next json Value
					curValue = slc[index]
				} else {
					// No, can't find index
					err = fmt.Errorf(errNoSuchPathMsg, fullPath, typ, index)
					return
				}
			}
		}

		out = curValue
		return
	}
}

// Convert a path into an Object or Array, such as .addresses[3].city, or [3].city, into a func(Value) Value that performs
// the lookup on an input Value.
// If any path part does not exist in the Value passed to the func, an Invalid Value is returned.
// If the given path is not valid, an error is returned.
func ParsePath(p string) (fn func(Value) (Value, error), err error) {
	// Convert path into lookups
	var lookups []tuple.Two[string, union.Two[string, int]]
	if lookups, err = parsePath(p); err != nil {
		return
	}

	// Convert lookups into a function
	fn = lookupsToFunc(lookups)
	return
}
