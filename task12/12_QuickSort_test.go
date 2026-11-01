package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIsBetter(t *testing.T) {
	tests := []struct {
		name     string
		p1       Participant
		p2       Participant
		expected bool
	}{
		{
			name:     "Больше решенных задач (6 > 4)",
			p1:       Participant{name: "gena", p: 6, f: 1000},
			p2:       Participant{name: "timofey", p: 4, f: 80},
			expected: true,
		},
		{
			name:     "Меньше решенных задач (2 < 4)",
			p1:       Participant{name: "gosha", p: 2, f: 90},
			p2:       Participant{name: "alla", p: 4, f: 100},
			expected: false,
		},
		{
			name:     "Одинаково задач, меньший штраф лучше (80 < 100)",
			p1:       Participant{name: "timofey", p: 4, f: 80},
			p2:       Participant{name: "alla", p: 4, f: 100},
			expected: true,
		},
		{
			name:     "Одинаково задач и штрафа, имя раньше по алфавиту (gosha < rita)",
			p1:       Participant{name: "gosha", p: 2, f: 90},
			p2:       Participant{name: "rita", p: 2, f: 90},
			expected: true,
		},
		{
			name:     "Полностью одинаковые участники (не лучше)",
			p1:       Participant{name: "alla", p: 4, f: 100},
			p2:       Participant{name: "alla", p: 4, f: 100},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBetter(tt.p1, tt.p2)
			if got != tt.expected {
				t.Errorf("isBetter(%v, %v) = %v; ожидалось %v", tt.p1, tt.p2, got, tt.expected)
			}
		})
	}
}

func TestQuickSortLogic(t *testing.T) {
	tests := []struct {
		name     string
		input    []Participant
		expected []string // ожидаемый порядок логинов
	}{
		{
			name: "Пример 1 из условия",
			input: []Participant{
				{name: "alla", p: 4, f: 100},
				{name: "gena", p: 6, f: 1000},
				{name: "gosha", p: 2, f: 90},
				{name: "rita", p: 2, f: 90},
				{name: "timofey", p: 4, f: 80},
			},
			expected: []string{"gena", "timofey", "alla", "gosha", "rita"},
		},
		{
			name: "Пример 2 из условия (все 0 0, только алфавитная сортировка)",
			input: []Participant{
				{name: "alla", p: 0, f: 0},
				{name: "gena", p: 0, f: 0},
				{name: "gosha", p: 0, f: 0},
				{name: "rita", p: 0, f: 0},
				{name: "timofey", p: 0, f: 0},
			},
			expected: []string{"alla", "gena", "gosha", "rita", "timofey"},
		},
		{
			name: "Граничный: N = 1 (один участник)",
			input: []Participant{
				{name: "solo", p: 10, f: 50},
			},
			expected: []string{"solo"},
		},
		{
			name: "Равные задачи, разные штрафы",
			input: []Participant{
				{name: "bob", p: 5, f: 200},
				{name: "alice", p: 5, f: 100},
				{name: "charlie", p: 5, f: 150},
			},
			expected: []string{"alice", "charlie", "bob"},
		},
		{
			name: "Равные задачи и штрафы, алфавитный порядок",
			input: []Participant{
				{name: "zack", p: 5, f: 100},
				{name: "alex", p: 5, f: 100},
				{name: "mary", p: 5, f: 100},
			},
			expected: []string{"alex", "mary", "zack"},
		},
		{
			name: "Уже отсортированный массив",
			input: []Participant{
				{name: "first", p: 10, f: 0},
				{name: "second", p: 9, f: 0},
				{name: "third", p: 8, f: 0},
				{name: "fourth", p: 7, f: 0},
			},
			expected: []string{"first", "second", "third", "fourth"},
		},
		{
			name: "Обратно отсортированный массив",
			input: []Participant{
				{name: "fourth", p: 7, f: 0},
				{name: "third", p: 8, f: 0},
				{name: "second", p: 9, f: 0},
				{name: "first", p: 10, f: 0},
			},
			expected: []string{"first", "second", "third", "fourth"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]Participant, len(tt.input))
			copy(arr, tt.input)

			if len(arr) > 0 {
				quickSort(arr, 0, len(arr)-1)
			}

			gotNames := make([]string, len(arr))
			for i, p := range arr {
				gotNames[i] = p.name
			}

			if !reflect.DeepEqual(gotNames, tt.expected) {
				t.Errorf("quickSort дал неверный порядок.\nПолучено:  %v\nОжидалось: %v", gotNames, tt.expected)
			}
		})
	}
}

// TestOutputFormatMain проверяет формат вывода всей программы через stdin/stdout.
// В условии сказано: "выведите по порядку их логины по одному в строке".
func TestOutputFormatMain(t *testing.T) {
	inputData := "5\nalla 4 100\ngena 6 1000\ngosha 2 90\nrita 2 90\ntimofey 4 80\n"
	expectedLogins := []string{"gena", "timofey", "alla", "gosha", "rita"}

	// Подменяем stdin
	oldStdin := os.Stdin
	rIn, wIn, _ := os.Pipe()
	_, _ = wIn.WriteString(inputData)
	wIn.Close()
	os.Stdin = rIn

	// Перехватываем stdout
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	main()

	wOut.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, rOut)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	cleanLines := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) > 0 {
			cleanLines = append(cleanLines, l)
		}
	}

	if !reflect.DeepEqual(cleanLines, expectedLogins) {
		t.Errorf("ОШИБКА ФОРМАТА ВЫВОДА!\nОжидался вывод каждого логина на отдельной строке:\n%v\nФактический вывод:\n%s",
			expectedLogins, buf.String())
	}
}
