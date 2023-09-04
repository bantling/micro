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
	errNilHandler   = fmt.Errorf("The handler cannot be nil")
)

const (
	NotFoundMsg         = "Not Found"
	MethodNotAllowedMsg = "Method Not Allowed"
)

const (
	UUIDGroup = "([0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12})"
)

// Handler is a REST handler
// The only difference compared to http.Handler is en extra argument of regex matches for url parts,
// so that the handler doesn't to re-evaluate the same regex a second time to get those values.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, urlParts []string)
}

// HandlerFunc is an adapter func for Handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request, urlParts []string)

// ServeHTTP is HandlerFunc adapter method
func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, urlParts []string) {
	hf(w, r, urlParts)
}

// A ServeMux handles REST requests by performing pattern matching that considers the method and URL, rather than just
// the URL alone like the http.ServeMux implementation.
// The zero value is ready to use.
type ServeMux struct {
	// [](method, regexp, handler)
	// A flat structure if used to disambiguate a 404 not found error from a 405 method not allowed error
	handlers []tuple.Three[string, *regexp.Regexp, Handler]
}

// Handle maps the given method and regex string to the given handler.
// The regex allows capturing path parts that are variable, like a UUID.
// The Handle can be called anytime, even while HTTP requests are being served.
func (rsm *ServeMux) Handle(method, pattern string, handler Handler) error {
	// The method cannot be empty
	if method == "" {
		return errEmptyMethod
	}

	// The pattern cannot be empty
	if pattern == "" {
		return errEmptyPattern
	}

	// The handler cannot be nil
	if handler == nil {
		return errNilHandler
	}

	// The pattern must be a valid regex
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// Add method, pattern, and handler triple
	funcs.SliceAdd(&rsm.handlers, tuple.Of3(method, regex, handler))

	return nil
}

// MustHandle is a must version of Handle
func (rsm *ServeMux) MustHandle(method, pattern string, handler Handler) {
	funcs.Must(rsm.Handle(method, pattern, handler))
}

// ServeHTTP is http.Handler interface method
func (rsm ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the method and url
	method, url := r.Method, (*r.URL).Path

	// Iterate the method, regexp, and handler triples to find the first method match and regexp url matches, if any
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
		http.Error(w, MethodNotAllowedMsg, http.StatusMethodNotAllowed)
	} else {
		// URL not found
		http.Error(w, NotFoundMsg, http.StatusNotFound)
	}
}
