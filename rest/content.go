package rest

// SPDX-License-Identifier: Apache-2.0

import (
  "github.com/bantling/micro/iter"
)

var (
  acceptEncoding = "Accept-Encoding"
  gzipEncoding   = "gzip"
)

// AcceptGzip checks if the request Accept-Encoding header is gzip.
// If so, it returns an Iter[byte] which decompresses using gzip.
// Otherwise, it returns an Iter[byte] of the request body as is.
func AcceptGzip(r *http.Request) Iter[byte] {
  it := iter.OfReader(r.Body)
  if r.Header.Get(acceptEncoding) == gzipEncoding {
    it = iter.
}

// NegotiateCSVContent returns an Iter[[]string] if the Content-Type header is text/csv, otherwise it returns nil
// If the Accept-Encoding header is gzip, then the content is decompressed with gzip first
func NegotiateCSVContent(r *http.Request) Iter[[]string] {
  
} 
