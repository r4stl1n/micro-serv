package types

import (
	"encoding/hex"
	"encoding/json"
	"crypto/rand"
	"strings"
)

// String fancy string structure to allow additional string functionality
type String struct{}

// Random generates a random string of letters
func (s String) Random(characterCount int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := characterCount
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

// Hex turns a given structure into a hex object
func (s String) Hex(in any) (string, error) {
	out, marshalError := json.Marshal(in)
	if marshalError != nil {
		return "", marshalError
	}

	return hex.EncodeToString(out), nil
}
