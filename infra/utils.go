package infra

import "math/rand"

const rs3Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs3Letters[int(rand.Int63()%int64(len(rs3Letters)))]
	}
	return string(b)
}
