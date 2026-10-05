package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode"
)

// A token is a chunk of text from a statement with a type
type token struct {
	text string
	typ  tokenTyp
}

// A tokenTyp identifies what kind of token something is
type tokenTyp int

const (
	// A bare word is a unquoted key; like 'foo' in json.foo = 1;
	typBare tokenTyp = iota

	// Numeric key; like '2' in json[2] = "foo";
	typNumericKey

	// A quoted key; like 'foo bar' in json["foo bar"] = 2;
	typQuotedKey

	// Punctuation types
	typDot    // .
	typLBrace // [
	typRBrace // ]
	typEquals // =
	typSemi   // ;
	typComma  // ,

	// Value types
	typString      // "foo"
	typNumber      // 4
	typTrue        // true
	typFalse       // false
	typNull        // null
	typEmptyArray  // []
	typEmptyObject // {}

	// Ignored token
	typIgnored

	// Error token
	typError
)

// a sprintFn adds color to its input
type sprintFn func(...interface{}) string

// mapping of token types to the appropriate color sprintFn
var sprintFns = map[tokenTyp]sprintFn{
	typBare:        bareColor.SprintFunc(),
	typNumericKey:  numColor.SprintFunc(),
	typQuotedKey:   strColor.SprintFunc(),
	typLBrace:      braceColor.SprintFunc(),
	typRBrace:      braceColor.SprintFunc(),
	typString:      strColor.SprintFunc(),
	typNumber:      numColor.SprintFunc(),
	typTrue:        boolColor.SprintFunc(),
	typFalse:       boolColor.SprintFunc(),
	typNull:        boolColor.SprintFunc(),
	typEmptyArray:  braceColor.SprintFunc(),
	typEmptyObject: braceColor.SprintFunc(),
}

// isValue returns true if the token is a valid value type
func (t token) isValue() bool {
	offset := uint(t.typ) - uint(typString)
	return offset <= uint(typEmptyObject-typString)
}

// isPunct returns true if the token is a punctuation type
func (t token) isPunct() bool {
	switch t.typ {
	case typDot, typLBrace, typRBrace, typEquals, typSemi, typComma:
		return true
	default:
		return false
	}
}

// format returns the formatted version of the token text
func (t token) format() string {
	if t.typ == typEquals {
		return " " + t.text + " "
	}
	return t.text
}

// formatColor returns the colored formatted version of the token text
func (t token) formatColor() string {
	text := t.format()

	// apply color after formatting so spacing stays identical
	if fn, ok := sprintFns[t.typ]; ok {
		return fn(text)
	}
	return text
}

// valueTokenFromInterface takes any valid value and
// returns a value token to represent it
func valueTokenFromInterface(v interface{}) token {
	switch vv := v.(type) {

	case map[string]interface{}:
		return token{"{}", typEmptyObject}
	case []interface{}:
		return token{"[]", typEmptyArray}
	case json.Number:
		return token{vv.String(), typNumber}
	case string:
		return token{quoteString(vv), typString}
	case bool:
		if vv {
			return token{"true", typTrue}
		}
		return token{"false", typFalse}
	case nil:
		return token{"null", typNull}
	default:
		return token{"", typError}
	}
}

// shortEscapes contains the exact spellings used in gron output
var shortEscapes = map[rune]string{
	'\\':     `\\`,
	'"':      `\"`,
	'\b':     `\b`,
	'\f':     `\f`,
	'\n':     `\n`,
	'\r':     `\r`,
	'\t':     `\t`,
	'\u2028': `\u2028`,
	'\u2029': `\u2029`,
}

// writeQuotedRune writes one rune without the surrounding quotes
func writeQuotedRune(out *bytes.Buffer, r rune) {
	if escaped, ok := shortEscapes[r]; ok {
		_, _ = out.WriteString(escaped)
		return
	}

	// Preserve uppercase hexadecimal escapes for other control runes.
	if unicode.IsControl(r) {
		_, _ = fmt.Fprintf(out, `\u%04X`, r)
		return
	}

	_, _ = out.WriteRune(r)
}

// quoteString takes a string and returns a quoted and
// escaped string valid for use in gron output.
func quoteString(s string) string {
	out := &bytes.Buffer{}

	// Buffer writes cannot fail; ignored errors match the original code.
	_ = out.WriteByte('"')

	for _, r := range s {
		writeQuotedRune(out, r)
	}

	_ = out.WriteByte('"')
	return out.String()
}
