package main

import "unicode"

// The javascript reserved words cannot be used as unquoted keys
var reservedWords = map[string]bool{
	"break":      true,
	"case":       true,
	"catch":      true,
	"class":      true,
	"const":      true,
	"continue":   true,
	"debugger":   true,
	"default":    true,
	"delete":     true,
	"do":         true,
	"else":       true,
	"export":     true,
	"extends":    true,
	"false":      true,
	"finally":    true,
	"for":        true,
	"function":   true,
	"if":         true,
	"import":     true,
	"in":         true,
	"instanceof": true,
	"new":        true,
	"null":       true,
	"return":     true,
	"super":      true,
	"switch":     true,
	"this":       true,
	"throw":      true,
	"true":       true,
	"try":        true,
	"typeof":     true,
	"var":        true,
	"void":       true,
	"while":      true,
	"with":       true,
	"yield":      true,
}

// validIdentifier checks to see if a string is a valid
// JavaScript identifier
// E.g:
//
//	justLettersAndNumbers1 -> true
//	a key with spaces      -> false
//	1startsWithANumber	  -> false
func validIdentifier(s string) bool {
	for i, r := range s {
		valid := validSecondaryRune(r)
		if i == 0 {
			valid = validFirstRune(r)
		}
		if !valid {
			return false
		}
	}

	return s != "" && !reservedWords[s]
}

// validFirstRune returns true for runes that are valid
// as the first rune in an identifier.
// E.g:
//
//	'r' -> true
//	'7' -> false
func validFirstRune(r rune) bool {
	return unicode.In(r,
		unicode.Lu,
		unicode.Ll,
		unicode.Lm,
		unicode.Lo,
		unicode.Nl,
	) || r == '$' || r == '_'
}

// validSecondaryRune returns true for runes that are valid
// as anything other than the first rune in an identifier.
func validSecondaryRune(r rune) bool {
	if unicode.In(r, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc) {
		return true
	}
	return validFirstRune(r)
}
