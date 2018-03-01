package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	// 現在のHashのhash b を加えて、その結果のsliceを返す
	sum := h.Sum(nil)
	// hexadecimal(16進数)をエンコーディングする
	return hex.EncodeToString(sum)
}