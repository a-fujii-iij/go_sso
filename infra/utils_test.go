package infra

import (
	"strings"
	"testing"
)

func TestRandString(t *testing.T) {
	length := 10
	str := RandString(length)
	if len(str) != length {
		t.Errorf("Expected string of length %d, got %d", length, len(str))
	}

	// 文字列がランダムであることを確認するために、複数回関数を実行して結果が異なることを確認します。
	str2 := RandString(length)
	if str == str2 {
		t.Errorf("Expected different strings, got two identical: %s", str)
	}

	// 生成された文字列が指定された文字セットのみを含んでいることを確認します。
	for _, char := range str {
		if !strings.ContainsRune(rs3Letters, char) {
			t.Errorf("String contains character not in the allowed set: %v", char)
		}
	}
}
