package main

import (
	"bufio"
	"os"
	"strconv"
)

const MAXN = 20005 // с запасом
const LOG = 16     // 2^15 = 32768, хватит для 20000 вершин

var (
	g      [MAXN][]int
	depthA [MAXN]int
	depthB [MAXN]int
	upA    [MAXN][LOG]int
	upB    [MAXN][LOG]int

	// Вспомогательные переменные
	dist         [MAXN]int
	farthestNode int
	maxDist      int
)

// для поиска самого далёкого края
func findFarthest(v, p, d int) {
	dist[v] = d
	if d > maxDist {
		maxDist = d
		farthestNode = v
	}
	for _, to := range g[v] {
		if to != p {
			findFarthest(to, v, d+1)
		}
	}
}

func buildTree(v, p, d int, depth *[MAXN]int, up *[MAXN][LOG]int) {
	depth[v] = d
	up[v][0] = p // предок на расстоянии 1 (2^0)

	// Заполняем прыжки на 2, 4, 8, 16 шагов (степени двойки)
	for k := 1; k < LOG; k++ {
		ancestor := up[v][k-1]
		up[v][k] = up[ancestor][k-1]
	}

	for _, to := range g[v] {
		if to != p {
			buildTree(to, v, d+1, depth, up)
		}
	}
}

// Функция мгновенного прыжка на d шагов вверх
func jump(v, d int, up *[MAXN][LOG]int) int {
	for k := LOG - 1; k >= 0; k-- {
		// Проверяем, включен ли k-й бит в числе d
		if (d & (1 << k)) != 0 {
			v = up[v][k]
		}
	}
	return v
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	readInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	n := readInt()
	q := readInt()

	for i := 0; i < n-1; i++ {
		u := readInt()
		v := readInt()
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	maxDist = -1
	findFarthest(1, 0, 0)
	nodeA := farthestNode

	maxDist = -1
	findFarthest(nodeA, 0, 0)
	nodeB := farthestNode

	buildTree(nodeA, nodeA, 0, &depthA, &upA)

	buildTree(nodeB, nodeB, 0, &depthB, &upB)

	for i := 0; i < q; i++ {
		v := readInt()
		d := readInt()

		if depthA[v] >= d {
			// До края A хватает шагов, значит прыгаем к A
			ans := jump(v, d, &upA)
			writer.WriteString(strconv.Itoa(ans) + "\n")
		} else if depthB[v] >= d {
			// До края B хватает шагов, значит к B
			ans := jump(v, d, &upB)
			writer.WriteString(strconv.Itoa(ans) + "\n")
		} else {
			// Никуда не хватает
			writer.WriteString("0\n")
		}
	}
}
