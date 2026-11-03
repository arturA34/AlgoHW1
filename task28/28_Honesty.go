package main

import (
	"bufio"
	"fmt"
	"os"
)

type Query struct {
	l      int
	r      int
	parity string
}

// find ищет главный элемент (корень) группы, в которой лежит x
func find(parent map[int]int, x int) int {
	if parent[x] == 0 {
		return x
	}
	parent[x] = find(parent, parent[x])
	return parent[x]
}

// union объединяет группы элементов a и b
func union(parent map[int]int, a, b int) {
	rootA := find(parent, a)
	rootB := find(parent, b)
	if rootA != rootB {
		parent[rootA] = rootB
	}
}

// solve принимает длину и список всех вопросов, возвращает номер первой лжи
func solve(length int, queries []Query) int {
	parent := make(map[int]int)

	for i, q := range queries {
		if q.l < 1 || q.r > length || q.l > q.r {
			return i
		}

		u := q.l
		v := q.r + 1

		if q.parity == "even" {
			// Четно значит u и v одной четности.
			// если u в группе с -v, значит ложь
			if find(parent, u) == find(parent, -v) {
				return i
			}
			union(parent, u, v)
			union(parent, -u, -v)
		} else {
			// Нечетно то есть u и v разной четности
			if find(parent, u) == find(parent, v) {
				return i
			}
			union(parent, u, -v)
			union(parent, -u, v)
		}
	}

	// если цикл закончился, значит все ответы честные
	return len(queries)
}

func main() {
	in := bufio.NewReader(os.Stdin)

	for {
		var length int
		if _, err := fmt.Fscan(in, &length); err != nil || length == -1 {
			break
		}

		var qCount int
		fmt.Fscan(in, &qCount)

		queries := make([]Query, qCount)
		for i := 0; i < qCount; i++ {
			fmt.Fscan(in, &queries[i].l, &queries[i].r, &queries[i].parity)
		}

		ans := solve(length, queries)
		fmt.Println(ans)
	}
}
