package iff

import (
	"fmt"
	"strings"
)

// helpers for output...

// RenderBytes is a helper for debug output
func RenderBytes(bs []byte, perLine int) string {
	lines := make([]string, 0, (len(bs)+perLine-1)/perLine)
	hex := make([]string, perLine)
	// ch := make([]string, perLine)

	for i := 0; i < len(bs); i += perLine {
		c := len(bs) - i
		if c > perLine {
			c = perLine
		}
		partial := bs[i : i+c]

		for j, b := range partial {
			hex[j] = fmt.Sprintf("%02x", b)
			// ch[j] = fmt.Sprintf("%c", b)
		}

		for j := c; j < perLine; j++ {
			hex[j] = "  "
			// ch[j] = " "
		}

		// lines = append(lines, fmt.Sprintf("        [%04x]: %s  -  %s", i, strings.Join(hex, " "), strings.Join(ch, " ")))
		lines = append(lines, fmt.Sprintf("        [%04x]: %s", i, strings.Join(hex, " ")))
	}

	return strings.Join(lines, "\n")
}
