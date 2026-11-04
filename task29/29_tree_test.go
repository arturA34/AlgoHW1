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

// bfsDist вычисляет кратчайшее расстояние между вершинами u и v с помощью BFS
func bfsDist(n int, adj [][]int, start, target int) int {
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	dist[start] = 0
	queue := []int{start}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == target {
			return dist[curr]
		}

		for _, nxt := range adj[curr] {
			if dist[nxt] == -1 {
				dist[nxt] = dist[curr] + 1
				queue = append(queue, nxt)
			}
		}
	}
	return -1
}

// maxDistFromTree возвращает максимальное расстояние от вершины start до любой вершины дерева
func maxDistFromTree(n int, adj [][]int, start int) int {
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	dist[start] = 0
	queue := []int{start}
	mx := 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if dist[curr] > mx {
			mx = dist[curr]
		}

		for _, nxt := range adj[curr] {
			if dist[nxt] == -1 {
				dist[nxt] = dist[curr] + 1
				queue = append(queue, nxt)
			}
		}
	}
	return mx
}

func resetGlobals() {
	for i := 0; i < MAXN; i++ {
		g[i] = nil
		depthA[i] = 0
		depthB[i] = 0
		dist[i] = 0
		for k := 0; k < LOG; k++ {
			upA[i][k] = 0
			upB[i][k] = 0
		}
	}
	farthestNode = 0
	maxDist = 0
}

func TestTreeQueriesSemanticChecker(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		edges   [][2]int
		queries [][2]int
	}{
		{
			name:    "Дерево-линия (1-2-3-4)",
			n:       4,
			edges:   [][2]int{{1, 2}, {2, 3}, {3, 4}},
			queries: [][2]int{{1, 0}, {1, 2}, {1, 3}, {1, 4}, {2, 1}},
		},
		{
			name:    "Граничный: N = 1 (одна вершина)",
			n:       1,
			edges:   [][2]int{},
			queries: [][2]int{{1, 0}, {1, 1}},
		},
		{
			name:    "Дерево-звезда (центр 1, лучи 2, 3, 4, 5)",
			n:       5,
			edges:   [][2]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}},
			queries: [][2]int{{2, 2}, {1, 1}, {1, 2}, {3, 0}},
		},
		{
			name:    "Разветвленное бинарное дерево (7 вершин)",
			n:       7,
			edges:   [][2]int{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {3, 7}},
			queries: [][2]int{{4, 4}, {4, 2}, {1, 3}, {4, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobals()

			// Собираем входной поток
			var inBuf bytes.Buffer
			inBuf.WriteString(fmt.Sprintf("%d %d\n", tt.n, len(tt.queries)))
			for _, e := range tt.edges {
				inBuf.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
			}
			for _, q := range tt.queries {
				inBuf.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
			}

			// Строим список смежности для независимой проверки
			adj := make([][]int, tt.n+1)
			for _, e := range tt.edges {
				adj[e[0]] = append(adj[e[0]], e[1])
				adj[e[1]] = append(adj[e[1]], e[0])
			}

			// Перехватываем stdin и stdout
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

			lines := strings.Fields(outBuf.String())
			if len(lines) != len(tt.queries) {
				t.Fatalf("Ожидалось %d ответов, получено %d", len(tt.queries), len(lines))
			}

			for i, q := range tt.queries {
				v := q[0]
				d := q[1]
				ans, err := strconv.Atoi(lines[i])
				if err != nil {
					t.Fatalf("Запрос #%d: не удалось распарсить число '%s'", i+1, lines[i])
				}

				if ans == 0 {
					// Проверяем: действительно ли нет вершины на расстоянии d?
					maxD := maxDistFromTree(tt.n, adj, v)
					if maxD >= d {
						t.Errorf("Запрос #%d (v=%d, d=%d): программа вернула 0, но вершина на расстоянии %d существует (макс. расстояние %d)!",
							i+1, v, d, d, maxD)
					}
				} else {
					if ans < 1 || ans > tt.n {
						t.Errorf("Запрос #%d (v=%d, d=%d): номер вершины %d вне диапазона [1, %d]",
							i+1, v, d, ans, tt.n)
						continue
					}
					realD := bfsDist(tt.n, adj, v, ans)
					if realD != d {
						t.Errorf("Запрос #%d (v=%d, d=%d): программа вернула вершину %d, но реальное расстояние = %d!",
							i+1, v, d, ans, realD)
					}
				}
			}
		})
	}
}
