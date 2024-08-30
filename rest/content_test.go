package rest

// SPDX-License-Identifier: Apache-2.0

import (
	"compress/gzip"
	goio "io"
	"net/http"
	"strings"
	"testing"

	"github.com/bantling/micro/encoding/json"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/iter"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestAcceptGzip_(t *testing.T) {
	// Has gzip content
	{
		var sb strings.Builder
		gzw := gzip.NewWriter(&sb)
		gzw.Write([]byte("foobar"))
		gzw.Close()

		r := funcs.MustValue(http.NewRequest("GET", "/foo", nil))
		r.Header.Set(acceptEncoding, gzipEncoding)
		r.Body = goio.NopCloser(strings.NewReader(sb.String()))

		gzr := AcceptGzip(r)
		bytes := funcs.MustValue(goio.ReadAll(gzr))
		assert.Nil(t, r.Body.Close())
		assert.NotEqual(t, gzr, r.Body)
		assert.Equal(t, "foobar", string(bytes))
	}

	// No gzip content
	{
		r := funcs.MustValue(http.NewRequest("GET", "/foo", nil))
		r.Body = goio.NopCloser(strings.NewReader("foobar"))

		gzr := AcceptGzip(r)
		bytes := funcs.MustValue(goio.ReadAll(gzr))
		assert.Nil(t, r.Body.Close())
		assert.Equal(t, gzr, r.Body)
		assert.Equal(t, "foobar", string(bytes))
	}
}

func TestNegotiateCSVContent_(t *testing.T) {
	// Has CSV content
	{
		str := `"FirstName","LastName"
"Jane","Doe"
`

		r := funcs.MustValue(http.NewRequest("GET", "/foo", nil))
		r.Header.Set(contentType, csvContent)
		r.Body = goio.NopCloser(strings.NewReader(str))

		it := NegotiateCSVContent(r)
		assert.NotNil(t, it)
		assert.Equal(t, union.OfResult([]string{"FirstName", "LastName"}), iter.Maybe(it))
		assert.Equal(t, union.OfResult([]string{"Jane", "Doe"}), iter.Maybe(it))
		assert.Equal(t, union.OfError[[]string](iter.EOI), iter.Maybe(it))
		assert.Nil(t, r.Body.Close())
	}

	// No csv content
	{
		r := funcs.MustValue(http.NewRequest("GET", "/foo", nil))
		r.Body = goio.NopCloser(strings.NewReader("foobar"))

		it := NegotiateCSVContent(r)
		assert.Nil(t, it)
		assert.Nil(t, r.Body.Close())
	}
}

func TestNegotiateJSONContent_(t *testing.T) {
	// Has JSON content
	{
		r := funcs.MustValue(http.NewRequest("GET", "/foo", nil))
		r.Header.Set(contentType, jsonContent)
		r.Body = goio.NopCloser(strings.NewReader(`{"FirstName": "Jane", "LastName": "Doe"}`))

		it := NegotiateJSONContent(r)
		assert.NotNil(t, it)
		assert.Equal(t, union.OfResult(json.MustToValue(map[string]any{"FirstName": "Jane", "LastName": "Doe"})), iter.Maybe(it))
		assert.Equal(t, union.OfError[json.Value](iter.EOI), iter.Maybe(it))
		assert.Nil(t, r.Body.Close())
	}

	// No json content
	{
		r := funcs.MustValue(http.NewRequest("GET", "/foo", nil))
		r.Body = goio.NopCloser(strings.NewReader("foobar"))

		it := NegotiateJSONContent(r)
		assert.Nil(t, it)
		assert.Nil(t, r.Body.Close())
	}
}
