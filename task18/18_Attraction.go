package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Operation struct {
	cmd  byte
	x, y int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	m, _ := strconv.Atoi(scanner.Text())

	rows := make([]int, n)
	cols := make([]int, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())

			// По первой строке определяем исходные колонки
			if i == 0 {
				cols[j] = (val-1)%m + 1
			}
			// По первому столбцу определяем исходные ряды
			if j == 0 {
				rows[i] = (val-1)/m + 1
			}
		}
	}

	var ops []Operation

	sortPermutation(rows, 'R', &ops)
	sortPermutation(cols, 'C', &ops)

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fprintln(writer, len(ops))
	for _, op := range ops {
		fmt.Fprintf(writer, "%c %d %d\n", op.cmd, op.x, op.y)
	}
}

func sortPermutation(arr []int, opType byte, ops *[]Operation) {
	for i := 0; i < len(arr); i++ {
		for arr[i] != i+1 {
			target := arr[i]

			*ops = append(*ops, Operation{cmd: opType, x: i + 1, y: target})

			arr[i], arr[target-1] = arr[target-1], arr[i]
		}
	}
}
