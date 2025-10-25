package utils

import (
	"fmt"
	"strings"
)

func ToUSD(cents int) string {
	centsStr := fmt.Sprintf("%d", cents%100)
	if len(centsStr) < 2 {
		centsStr = "0" + centsStr
	}
	dollars := cents / 100

	if dollars >= 1000 {
		thousands := dollars / 1000
		hundreds := dollars % 1000
		return fmt.Sprintf("$%d,%03d.%s", thousands, hundreds, centsStr)
	}

	return fmt.Sprintf("$%d.%s", dollars, centsStr)
}

// EscapePDFString escapes special characters in PDF strings
func EscapePDFString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}

// GetStringWidth estimates the width of a string in points
// Based on standard Helvetica character widths
func GetStringWidth(s string, fontSize int) float64 {
	// Standard Helvetica character widths (relative to 1000 units)
	widths := map[rune]int{
		' ': 278, '!': 278, '"': 355, '#': 556, '$': 556, '%': 889, '&': 667, '\'': 191,
		'(': 333, ')': 333, '*': 389, '+': 584, ',': 278, '-': 333, '.': 278, '/': 278,
		'0': 556, '1': 556, '2': 556, '3': 556, '4': 556, '5': 556, '6': 556, '7': 556,
		'8': 556, '9': 556, ':': 278, ';': 278, '<': 584, '=': 584, '>': 584, '?': 556,
		'@': 1015, 'A': 667, 'B': 667, 'C': 722, 'D': 722, 'E': 667, 'F': 611, 'G': 778,
		'H': 722, 'I': 278, 'J': 500, 'K': 667, 'L': 556, 'M': 833, 'N': 722, 'O': 778,
		'P': 667, 'Q': 778, 'R': 722, 'S': 667, 'T': 611, 'U': 722, 'V': 667, 'W': 944,
		'X': 667, 'Y': 667, 'Z': 611, '[': 278, '\\': 278, ']': 278, '^': 469, '_': 556,
		'`': 333, 'a': 556, 'b': 556, 'c': 500, 'd': 556, 'e': 556, 'f': 278, 'g': 556,
		'h': 556, 'i': 222, 'j': 222, 'k': 500, 'l': 222, 'm': 833, 'n': 556, 'o': 556,
		'p': 556, 'q': 556, 'r': 333, 's': 500, 't': 278, 'u': 556, 'v': 500, 'w': 722,
		'x': 500, 'y': 500, 'z': 500, '{': 334, '|': 260, '}': 334, '~': 584,
	}

	totalWidth := 0
	for _, ch := range s {
		if w, ok := widths[ch]; ok {
			totalWidth += w
		} else {
			totalWidth += 556 // Default width for unknown characters
		}
	}

	// Convert from 1000-unit scale to actual points
	return float64(totalWidth) * float64(fontSize) / 1000.0
}
