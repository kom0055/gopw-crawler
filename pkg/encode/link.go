package encode

import (
	_ "unsafe"
)

//go:linkname Escape net/url.escape
func Escape(s string, mode int) string
