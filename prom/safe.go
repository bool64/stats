// +build appengine

// This file contains the safe implementations of otherwise unsafe-using code.

package prom

func byteSlice2String(bs []byte) string {
	return string(bs)
}
