package iter

// SPDX-License-Identifier: Apache-2.0

// IterReader adapts an iter[byte] into an io.Reader
type IterReader struct {
  it Iter[byte]
}

// Next is Iter.Next interface method
func (ir IterReader) Next() (byte, error) {
  return ir.it.Next()
}

// Unread is Iter.Unread interface method
func (ir IterReader) Unread(b byte) {
  ir.it.Unread(b)
}

// ToReader converts an Iter[byte] into an io.Reader
func ToReader(it Iter[byte]) IterReader {
  return IterReader{it}
}

// Read is the io.Reader interface
func (ir *IterReader) Read(p []byte) (n int, err error) {
  var (
    l = len(p)
    val byte
  )

  // If p is zero length, then return 0, nil
  if l == 0 {
    return 0, nil
  }
  
  // Read up to len(p) bytes, possibly encountering an error along the way
  for n = 0; n < l; n++ {
    val, err = ir.it.Next()
    if err != nil {
      return
    }
    
    p[n] = val
  }
  
  return
}

// OfGzip constructs an Iter[byte] from an existing iter[byte] that decompresses the input bytes using gzip
/*func OfGzip(it Iter[byte]) Iter[byte] {
  var done bool
   
  return OfIter(func() (byte, error) {
    if done {
      var zv U
    val, err := it.Next()
    if err == nil {
      return mapper(val), nil
    }

    var zv U
    return zv, err
  })
}*/
