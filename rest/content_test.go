package rest

// SPDX-License-Identifier: Apache-2.0

import (
  "compress/gzip"
  goio "io"
  "net/http"
  "strings"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestAcceptGzip_(t *testing.T) {
  r, err := http.NewRequest("GET", "/foo", nil)
  assert.Nil(t, err)
  r.Header.Set(acceptEncoding, gzipEncoding)
  
  var sb strings.Builder
  gzw := gzip.NewWriter(&sb)
  gzw.Write([]byte("foobar"))
  r.Body = goio.NopCloser(strings.NewReader(sb.String()))
  
  gz := AcceptGzip(r)
  assert.NotEqual(t, gz, r.Body)
  bytes, err := goio.ReadAll(gz)
  assert.Equal(t, "foobar", string(bytes))
  assert.Nil(t, err)
}
