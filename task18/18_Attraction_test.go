package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestSortPermutation тестирует локальную функцию сортировки перестановки
func TestSortPermutation(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		opType byte
	}{
		{
			name:   "Длина 1 (уже на месте)",
			arr:    []int{1},
			opType: 'R',
		},
		{
			name:   "Длина 2 (инверсия)",
			arr:    []int{2, 1},
			opType: 'R',
		},
		{
			name:   "Циклическая перестановка [3, 1, 2]",
			arr:    []int{3, 1, 2},
			opType: 'C',
		},
		{
			name:   "Полностью обратная [4, 3, 2, 1]",
			arr:    []int{4, 3, 2, 1},
			opType: 'R',
		},
		{
			name:   "Уже отсортированная [1, 2, 3, 4, 5]",
			arr:    []int{1, 2, 3, 4, 5},
			opType: 'C',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyArr := make([]int, len(tt.arr))
			copy(copyArr, tt.arr)
			var ops []Operation

			sortPermutation(copyArr, tt.opType, &ops)

			// Проверяем, что массив стал 1, 2, ..., N
			for i, v := range copyArr {
				if v != i+1 {
					t.Fatalf("массив после сортировки не отсортирован: %v (позиция %d: %d)", copyArr, i, v)
				}
			}

			// Проверяем, что число операций не превышает len(arr) - 1
			if len(ops) > len(tt.arr) {
				t.Errorf("слишком много операций: %d при длине %d", len(ops), len(tt.arr))
			}

			// Проверяем корректность каждой операции
			for _, op := range ops {
				if op.cmd != tt.opType {
					t.Errorf("неверный тип команды: got %c, want %c", op.cmd, tt.opType)
				}
				if op.x < 1 || op.x > len(tt.arr) || op.y < 1 || op.y > len(tt.arr) {
					t.Errorf("индекс операции вне диапазона: x=%d, y=%d (len=%d)", op.x, op.y, len(tt.arr))
				}
			}
		})
	}
}

// TestAttractionMatrixSemanticChecker выполняет сквозное тестирование:
// запускает main() с матрицей на входе, забирает операции и применяет их к исходной матрице.
func TestAttractionMatrixSemanticChecker(t *testing.T) {
	tests := []struct {
		name string
		n, m int
		grid [][]int
	}{
		{
			name: "Пример 1 из условия (2x2)",
			n:    2,
			m:    2,
			grid: [][]int{
				{4, 3},
				{2, 1},
			},
		},
		{
			name: "Граничный: 1x1 (минимальный размер)",
			n:    1,
			m:    1,
			grid: [][]int{
				{1},
			},
		},
		{
			name: "Граничный: 1x4 (одна строка, только перестановка колонок)",
			n:    1,
			m:    4,
			grid: [][]int{
				{4, 2, 3, 1},
			},
		},
		{
			name: "Граничный: 4x1 (один столбец, только перестановка строк)",
			n:    4,
			m:    1,
			grid: [][]int{
				{4},
				{1},
				{3},
				{2},
			},
		},
		{
			name: "Граничный: 3x3 уже в исходном положении",
			n:    3,
			m:    3,
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			name: "Сложный 3x3 с переставленными строками и столбцами",
			n:    3,
			m:    3,
			grid: [][]int{
				{9, 7, 8},
				{3, 1, 2},
				{6, 4, 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Собираем stdin
			var inBuf bytes.Buffer
			inBuf.WriteString(fmt.Sprintf("%d %d\n", tt.n, tt.m))
			for r := 0; r < tt.n; r++ {
				for c := 0; c < tt.m; c++ {
					inBuf.WriteString(fmt.Sprintf("%d ", tt.grid[r][c]))
				}
				inBuf.WriteString("\n")
			}

			// Подменяем stdin и stdout
			oldStdin := os.Stdin
			rIn, wIn, _ := os.Pipe()
			_, _ = wIn.Write(inBuf.Bytes())
			wIn.Close()
			os.Stdin = rIn

			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			main()

			wOut.Close()
			os.Stdin = oldStdin
			os.Stdout = oldStdout

			var outBuf bytes.Buffer
			_, _ = io.Copy(&outBuf, rOut)

			// Парсим операции из вывода
			lines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
			if len(lines) == 0 || len(lines[0]) == 0 {
				t.Fatalf("Пустой вывод программы")
			}

			k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
			if err != nil {
				t.Fatalf("Первая строка должна содержать число K, получено: %s", lines[0])
			}

			if k > 1000 {
				t.Fatalf("Число операций %d превышает лимит 1000 по условию!", k)
			}

			if len(lines)-1 != k {
				t.Fatalf("Заявлено K=%d операций, но строк с операциями: %d", k, len(lines)-1)
			}

			// Копируем матрицу и применяем операции
			curGrid := make([][]int, tt.n)
			for r := 0; r < tt.n; r++ {
				curGrid[r] = make([]int, tt.m)
				copy(curGrid[r], tt.grid[r])
			}

			for opIdx, line := range lines[1:] {
				parts := strings.Fields(line)
				if len(parts) != 3 {
					t.Fatalf("Операция #%d имеет неверный формат: %s", opIdx+1, line)
				}
				cmd := parts[0]
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])

				xIdx := x - 1
				yIdx := y - 1

				switch cmd {
				case "R":
					if xIdx < 0 || xIdx >= tt.n || yIdx < 0 || yIdx >= tt.n {
						t.Fatalf("Операция #%d: строки вне диапазона [1, %d]: %d, %d", opIdx+1, tt.n, x, y)
					}
					curGrid[xIdx], curGrid[yIdx] = curGrid[yIdx], curGrid[xIdx]
				case "C":
					if xIdx < 0 || xIdx >= tt.m || yIdx < 0 || yIdx >= tt.m {
						t.Fatalf("Операция #%d: колонки вне диапазона [1, %d]: %d, %d", opIdx+1, tt.m, x, y)
					}
					for r := 0; r < tt.n; r++ {
						curGrid[r][xIdx], curGrid[r][yIdx] = curGrid[r][yIdx], curGrid[r][xIdx]
					}
				default:
					t.Fatalf("Неизвестная команда '%s' в операции #%d", cmd, opIdx+1)
				}
			}

			// Проверяем, что матрица вернулась в исходное состояние (r*M + c + 1)
			for r := 0; r < tt.n; r++ {
				for c := 0; c < tt.m; c++ {
					expectedVal := r*tt.m + c + 1
					if curGrid[r][c] != expectedVal {
						t.Fatalf("Кресло в позиции (%d, %d) имеет номер %d, ожидалось %d",
							r+1, c+1, curGrid[r][c], expectedVal)
					}
				}
			}
		})
	}
}
