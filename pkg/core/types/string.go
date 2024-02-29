package types

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"go.uber.org/zap"
	"math/big"
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

		randInt, randError := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))

		if randError != nil {
			zap.L().Error("failed to generate random integer", zap.Error(randError))
			continue
		}

		b.WriteRune(chars[randInt.Int64()])

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
