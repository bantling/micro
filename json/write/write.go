package write

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/json"
	"github.com/bantling/micro/writer"
)

// Write writes any value to a writer
func Write(jv json.Value, dst writer.Writer[rune]) error {
	switch jv.Type() {
	case json.Object:
		return WriteObject(jv, dst)
	case json.Array:
		return WriteArray(jv, dst)
	case json.String:
		return dst.Write(append(append([]rune{'"'}, []rune(jv.AsString())...), '"')...)
	case json.Number:
		fallthrough
	case json.Boolean:
		return dst.Write([]rune(jv.AsString())...)
	}

	return dst.Write('n', 'u', 'l', 'l')
}

// MustWrite is a must version of Write
func MustWrite(jv json.Value, dst writer.Writer[rune]) {
	funcs.Must(Write(jv, dst))
}

// WriteObject writes an object to a writer
func WriteObject(jv json.Value, dst writer.Writer[rune]) error {
	if err := dst.Write('{'); err != nil {
		return err
	}

	var i = 0
	for k, v := range jv.AsMap() {
		var data []rune

		if i > 0 {
			data = append(data, ',')
		}
		i++

		data = append(append(append(data, '"'), []rune(k)...), '"', ':')
		if err := dst.Write(data...); err != nil {
			return err
		}

		if err := Write(v, dst); err != nil {
			return err
		}
	}

	return dst.Write('}')
}

// MustWriteObject is a must version of WriteObject
func MustWriteObject(jv json.Value, dst writer.Writer[rune]) {
	funcs.Must(WriteObject(jv, dst))
}

// WriteArray writes an array to a writer
func WriteArray(jv json.Value, dst writer.Writer[rune]) error {
	if err := dst.Write('['); err != nil {
		return err
	}

	for i, v := range jv.AsSlice() {
		if i > 0 {
			if err := dst.Write(','); err != nil {
				return err
			}
		}

		if err := Write(v, dst); err != nil {
			return err
		}
	}

	return dst.Write(']')
}

// MustWriteArray is a must version of WriteArray
func MustWriteArray(jv json.Value, dst writer.Writer[rune]) {
	funcs.Must(WriteArray(jv, dst))
}
