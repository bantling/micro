// Package rest provides functions for simple REST handling
//
// SPDX-License-Identifier: Apache-2.0
package rest

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
)

var (
	errEmptyMethod  = fmt.Errorf("The method cannot be empty")
	errEmptyPattern = fmt.Errorf("The pattern cannot be empty")
)

const (
	notFoundMsg         = "Not Found"
	methodNotAllowedMsg = "Method Not Allowed"
)

// Handler is a REST handler
// The only difference compared to http.Handler is en extra argument of regex matches for url parts,
// so that the handler doesn't to reevaluate the same regex a second time to get those values.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, urlParts []string)
}

// A RESTServeMux handles REST requests by performing pattern matching that considers the method and URL, rather than just
// the URL alone like the http.ServeMux implementation.
type RESTServeMux struct {
	// [](method, regexp, handler)
	// A flat structure if used to disambiguate a 404 not found error from a 405 method not allowed error
	handlers []tuple.Three[string, *regexp.Regexp, Handler]
}

// Handle maps the given method and regex string to the given handler.
// The regex allows capturing path parts that are variable, like a UUID.
func (rsm *RESTServeMux) Handle(method, pattern string, handler Handler) error {
	// The method and pattern cannot be empty
	if method == "" {
		return errEmptyMethod
	}

	if pattern == "" {
		return errEmptyPattern
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	funcs.SliceAdd(&rsm.handlers, tuple.Of3(method, regex, handler))

	return nil
}

// MustHandle is a must version of Handle
func (rsm *RESTServeMux) MustHandle(method, pattern string, handler Handler) {
	funcs.Must(rsm.Handle(method, pattern, handler))
}

// ServeHTTP is http.Handler interface method
func (rsm RESTServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the method and url
	method, url := r.Method, (*r.URL).Path

	// Iterate the method/regexp/handler triples to find the first method match and regexp that matches the url, if any
	var urlMatched bool

	for _, methodPatternHandler := range rsm.handlers {
		mMethod, mPattern, mHandler := methodPatternHandler.Values()
		parts := mPattern.FindStringSubmatch(url)

		if parts != nil {
			// There is at least one match for the URL, but maybe not for the method
			urlMatched = true

			if method == mMethod {
				// Found match, call it and stop searching
				// ServeHTTP is our Handler interface method that accepts an extra arg of matching url parts
				mHandler.ServeHTTP(w, r, parts)
				return
			}
		}
	}

	// If no match is found, it's url or method not found
	if urlMatched {
		// Method not found
		http.Error(w, methodNotAllowedMsg, http.StatusMethodNotAllowed)
	} else {
		http.Error(w, notFoundMsg, http.StatusNotFound)
	}
}
