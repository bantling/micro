package rest

// SPDX-License-Identifier: Apache-2.0

import (
  "compress/gzip"
  goio "io"
  "net/http"
  "strings"
  "testing"

  "github.com/bantling/micro/funcs"
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
