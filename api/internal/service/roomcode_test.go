package service

import "testing"

func TestHasPairedConsecutive(t *testing.T) {
	cases := map[string]bool{
		"1123": true,
		"2255": true,
		"1234": false,
		"5571": true, // 设计示例：含相邻 55，属于优选的两两连号
		"1212": false,
		"1110": true,
		"0011": true,
	}
	for in, want := range cases {
		if got := hasPairedConsecutive(in); got != want {
			t.Errorf("hasPairedConsecutive(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestGen4DigitPreferredUnique(t *testing.T) {
	used := map[string]bool{}
	// claim 原子占用：未占用则占用并返回 true，已占用返回 false。
	claim := func(code string) bool {
		if used[code] {
			return false
		}
		used[code] = true
		return true
	}
	for i := 0; i < 50; i++ {
		code := gen4DigitPreferred(claim)
		if len(code) != 4 {
			t.Fatalf("expected 4-digit code, got %q", code)
		}
	}
	if len(used) != 50 {
		t.Fatalf("expected 50 unique codes, got %d", len(used))
	}
}

func TestGen5DigitFallback(t *testing.T) {
	claim := func(code string) bool { return true }
	code := gen5Digit(claim)
	if len(code) != 5 {
		t.Fatalf("expected 5-digit code, got %q", code)
	}
}
