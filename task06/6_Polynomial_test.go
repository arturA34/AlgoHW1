package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func capturePolynomialOutput(coefficients []int, args []int, mod int) []int {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	calculatePolynomial(coefficients, args, mod)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	outputStr := strings.TrimSpace(buf.String())
	if len(outputStr) == 0 {
		return nil
	}

	lines := strings.Split(outputStr, "\n")
	res := make([]int, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if val, err := strconv.Atoi(line); err == nil {
			res = append(res, val)
		}
	}
	return res
}

func TestCalculatePolynomial(t *testing.T) {
	tests := []struct {
		name         string
		coefficients []int
		args         []int
		mod          int
		expected     []int
	}{
		{
			name:         "Пример 1 из условия (x^2 + 5x + 4 mod 10)",
			coefficients: []int{1, 5, 4},
			args:         []int{0, 1, 2, 3, 4},
			mod:          10,
			expected:     []int{4, 0, 8, 8, 0},
		},
		{
			name:         "Граничный: x = 0 (результат равен свободному члену a0 mod MOD)",
			coefficients: []int{7, 0, 0, 42},
			args:         []int{0},
			mod:          100,
			expected:     []int{42},
		},
		{
			name:         "Граничный: x = 1 (сумма всех коэффициентов mod MOD)",
			coefficients: []int{10, 20, 30, 40},
			args:         []int{1},
			mod:          1000,
			expected:     []int{100},
		},
		{
			name:         "Граничный: Коэффициенты больше MOD (a_i > MOD)",
			coefficients: []int{105, 200, 33},
			args:         []int{2},
			mod:          10,
			// P(2) = 105*4 + 200*2 + 33 = 420 + 400 + 33 = 853 -> 853 % 10 = 3
			expected: []int{3},
		},
		{
			name:         "Граничный: Все коэффициенты равны 0 (нулевой полином)",
			coefficients: []int{0, 0, 0},
			args:         []int{5, 10},
			mod:          100,
			expected:     []int{0, 0},
		},
		{
			name:         "Граничный: Большой модуль 10^9 и большие числа",
			coefficients: []int{1000000000, 1000000000, 1000000000},
			args:         []int{1000000000, 0},
			mod:          1000000000,
			expected:     []int{0, 0},
		},
		{
			name: "Большая степень N = 1000 (проверка корректности схемы Горнера)",
			// 1*x^1000 + 0 + ... + 0 при x = 1
			coefficients: func() []int {
				c := make([]int, 1001)
				c[0] = 1
				return c
			}(),
			args:     []int{1, 2},
			mod:      7,
			expected: []int{1, 2}, // 1^1000 mod 7 = 1; 2^1000 mod 7 = (2^3)^333 * 2 = 1^333 * 2 = 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := capturePolynomialOutput(tt.coefficients, tt.args, tt.mod)
			if len(got) != len(tt.expected) {
				t.Fatalf("получено строк %d, ожидалось %d (%v)", len(got), len(tt.expected), got)
			}
			for i := range tt.expected {
				if got[i] != tt.expected[i] {
					t.Errorf("аргумент #%d (x=%d): получено %d, ожидалось %d", i, tt.args[i], got[i], tt.expected[i])
				}
			}
		})
	}
}
