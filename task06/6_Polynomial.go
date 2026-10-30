package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	s.Scan()
	data := strings.Split(s.Text(), " ")
	n, _ := strconv.Atoi(data[0])
	m, _ := strconv.Atoi(data[1])
	mod, _ := strconv.Atoi(data[2])

	coefficients := make([]int, 0, n+1)

	for i := 0; i < n+1; i++ {
		s.Scan()
		coef, _ := strconv.Atoi(s.Text())
		coefficients = append(coefficients, coef)
	}

	args := make([]int, 0, m+1)

	for j := 0; j < m; j++ {
		s.Scan()
		arg, _ := strconv.Atoi(s.Text())
		args = append(args, arg)
	}

	calculatePolynomial(coefficients, args, mod)
}

func calculatePolynomial(coefficients []int, args []int, mod int) {
	for _, arg := range args {
		var res int
		// Логика подсчета
		// Чтобы не возводить аргумент в большую степень, разбиваем многочлен на скобки формата ((ax + b)x +c)x + ...
		// Остаток от деления суммы равен сумме деления остатков
		for _, val := range coefficients {
			res = (res*arg + val) % mod
		}
		fmt.Println(res)
	}
}
