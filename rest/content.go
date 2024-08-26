package rest

// SPDX-License-Identifier: Apache-2.0

import (
  "compress/gzip"
  goio "io"
  "net/http"
  
  "github.com/bantling/micro/funcs"
  "github.com/bantling/micro/iter"
)

var (
  acceptEncoding = "Accept-Encoding"
  gzipEncoding   = "gzip"
  
  contentType    = "Content-Type"
  csvContent     = "text/csv"
)

// AcceptGzip checks if the request Accept-Encoding header is gzip.
// If so, it returns a new io.Reader which decompresses the body io.Reader using gzip.
// Otherwise, it returns the body io.Reader as is.
func AcceptGzip(r *http.Request) goio.Reader {
  if r.Header.Get(acceptEncoding) == gzipEncoding {
    return funcs.MustValue(gzip.NewReader(r.Body))
  }
  
  return r.Body
}

// NegotiateCSVContent returns an Iter[[]string] if the Content-Type header is text/csv, otherwise it returns nil
func NegotiateCSVContent(r *http.Request) iter.Iter[[]string] {
  if r.Header.Get(contentType) == csvContent {
    return iter.OfCSV(r.Body)
  }
  
  return nil
}
