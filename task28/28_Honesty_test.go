package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestSolveParity(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		queries  []Query
		expected int
	}{
		{
			name:   "Прямое противоречие на 2-м вопросе (1 2 even vs 1 2 odd)",
			length: 10,
			queries: []Query{
				{l: 1, r: 2, parity: "even"},
				{l: 1, r: 2, parity: "odd"},
			},
			expected: 1, // первые 1 условий непротиворечивы, 2-е ложно
		},
		{
			name:   "Транзитивное противоречие на 3-м вопросе",
			length: 10,
			queries: []Query{
				{l: 1, r: 2, parity: "even"}, // S2 == S0
				{l: 3, r: 4, parity: "even"}, // S4 == S2 => S4 == S0
				{l: 1, r: 4, parity: "odd"},  // S4 != S0 -> противоречие!
				{l: 1, r: 4, parity: "even"},
			},
			expected: 2, // первые 2 условия непротиворечивы
		},
		{
			name:   "Все ответы друга правдивы",
			length: 100,
			queries: []Query{
				{l: 1, r: 10, parity: "even"},
				{l: 11, r: 20, parity: "odd"},
				{l: 1, r: 20, parity: "odd"},
			},
			expected: 3, // все 3 условия выполняются
		},
		{
			name:     "Граничный: 0 вопросов",
			length:   50,
			queries:  []Query{},
			expected: 0,
		},
		{
			name:   "Граничный: Максимальная длина 10^9",
			length: 1000000000,
			queries: []Query{
				{l: 1, r: 1000000000, parity: "even"},
				{l: 1, r: 1000000000, parity: "odd"},
			},
			expected: 1,
		},
		{
			name:   "Граничный: Отрезок из 1 элемента [3, 3]",
			length: 5,
			queries: []Query{
				{l: 3, r: 3, parity: "even"},
				{l: 3, r: 3, parity: "odd"},
			},
			expected: 1,
		},
		{
			name:   "Пересекающиеся отрезки: [1, 5] odd и [2, 5] even => [1, 1] odd",
			length: 10,
			queries: []Query{
				{l: 1, r: 5, parity: "odd"},
				{l: 2, r: 5, parity: "even"},
				// Теперь известно, что 1-й бит нечетный (1). Если друг говорит [1, 1] even - это ложь
				{l: 1, r: 1, parity: "even"},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solve(tt.length, tt.queries)
			if got != tt.expected {
				t.Errorf("solve(%d, queries) = %d; ожидалось %d", tt.length, got, tt.expected)
			}
		})
	}
}

// TestHonestyMainMultipleTests тестирует обработку нескольких тестов во входном потоке до -1
func TestHonestyMainMultipleTests(t *testing.T) {
	inputData := "10\n2\n1 2 even\n1 2 odd\n5\n1\n1 5 odd\n-1\n"
	expectedAnswers := []int{1, 1}

	oldStdin := os.Stdin
	rIn, wIn, _ := os.Pipe()
	_, _ = wIn.WriteString(inputData)
	wIn.Close()
	os.Stdin = rIn

	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	main()

	wOut.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, rOut)

	lines := strings.Fields(buf.String())
	if len(lines) != len(expectedAnswers) {
		t.Fatalf("Ожидалось %d ответов, получено %d: %v", len(expectedAnswers), len(lines), lines)
	}

	for i, exp := range expectedAnswers {
		val, err := strconv.Atoi(lines[i])
		if err != nil || val != exp {
			t.Errorf("Тест #%d: получено %s, ожидалось %d", i+1, lines[i], exp)
		}
	}
}
