package utils

import "crypto/sha256"

func Hash(raw string) string {
	h := sha256.New()

	h.Write([]byte(raw))
	bs := h.Sum(nil)

	return string(bs[:])
}
