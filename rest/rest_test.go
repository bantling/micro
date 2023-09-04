// Package rest provides functions for simple REST handling
//
// SPDX-License-Identifier: Apache-2.0
package rest

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	urlParts []string
}

func (t *testData) ServeHTTP(w http.ResponseWriter, r *http.Request, urlParts []string) {
	t.urlParts = urlParts
}

func TestHandle_(t *testing.T) {
	var (
		mux  ServeMux
		td   testData
		h    Handler = &td
		pat1         = "^/customer$"
		pat2         = "^/customer/" + UUIDGroup + "$"
	)

	// handle GET /customer
	assert.Nil(t, mux.Handle("GET", pat1, h))
	assert.Equal(
		t,
		funcs.SliceOf(tuple.Of3("GET", regexp.MustCompile(pat1), h)),
		mux.handlers,
	)

	// test GET /customer
	recorder := httptest.NewRecorder()
	td.urlParts = nil
	mux.ServeHTTP(recorder, httptest.NewRequest("GET", "/customer", nil))
	assert.Equal(t, []string{"/customer"}, td.urlParts)
	assert.Equal(t, "200 OK", recorder.Result().Status)
	assert.Equal(t, 200, recorder.Result().StatusCode)

	// handle GET /customer/UUID
	mux.MustHandle("GET", pat2, h)
	assert.Equal(
		t,
		funcs.SliceOf(
			tuple.Of3("GET", regexp.MustCompile(pat1), h),
			tuple.Of3("GET", regexp.MustCompile(pat2), h),
		),
		mux.handlers,
	)

	// test GET /customer/UUID
	recorder = httptest.NewRecorder()
	td.urlParts = nil
	mux.ServeHTTP(recorder, httptest.NewRequest("GET", "/customer/12345678-9ABC-DEF0-1234-567890abcdef", nil))
	assert.Equal(t, []string{"/customer/12345678-9ABC-DEF0-1234-567890abcdef", "12345678-9ABC-DEF0-1234-567890abcdef"}, td.urlParts)
	assert.Equal(t, "200 OK", recorder.Result().Status)
	assert.Equal(t, 200, recorder.Result().StatusCode)

	// handler PUT /customer/UUID
	mux.MustHandle("PUT", pat2, HandlerFunc(func(w http.ResponseWriter, r *http.Request, urlParts []string) {
		td.urlParts = urlParts
	}))

	// test PUT /customer/UUID
	recorder = httptest.NewRecorder()
	td.urlParts = nil
	mux.ServeHTTP(recorder, httptest.NewRequest("PUT", "/customer/12345678-9ABC-DEF0-1234-567890abcdef", nil))
	assert.Equal(t, []string{"/customer/12345678-9ABC-DEF0-1234-567890abcdef", "12345678-9ABC-DEF0-1234-567890abcdef"}, td.urlParts)
	assert.Equal(t, "200 OK", recorder.Result().Status)
	assert.Equal(t, 200, recorder.Result().StatusCode)

	// test method not found
	recorder = httptest.NewRecorder()
	td.urlParts = nil
	mux.ServeHTTP(recorder, httptest.NewRequest("GET", "/foo", nil))
	assert.Equal(t, []string(nil), td.urlParts)
	assert.Equal(t, "404 Not Found", recorder.Result().Status)
	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)

	// test url not found
	recorder = httptest.NewRecorder()
	td.urlParts = nil
	mux.ServeHTTP(recorder, httptest.NewRequest("PUT", "/customer", nil))
	assert.Equal(t, []string(nil), td.urlParts)
	assert.Equal(t, "405 Method Not Allowed", recorder.Result().Status)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Result().StatusCode)

	// test errors in handle method
	assert.Equal(t, errEmptyMethod, mux.Handle("", pat1, h))
	assert.Equal(t, errEmptyPattern, mux.Handle("GET", "", h))
	assert.Equal(t, errNilHandler, mux.Handle("GET", pat1, nil))
	assert.Equal(t, funcs.SecondValue2(regexp.Compile("(")), mux.Handle("GET", "(", h))
}
