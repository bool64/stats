// +build !appengine

// This file encapsulates usage of unsafe.
// safe.go contains the safe implementations.

package prom

import (
	"unsafe"
)

func byteSlice2String(bs []byte) string {
	// nolint:gosec // Copied from https://github.com/golang/go/blob/go1.14.7/src/strings/builder.go#L46-L49.
	return *(*string)(unsafe.Pointer(&bs))
}
