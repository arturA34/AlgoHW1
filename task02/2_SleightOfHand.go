package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	k, _ := strconv.Atoi(scanner.Text())

	keyboardCharacters := make(map[int]int)

	for i := 0; i < 4; i++ {
		scanner.Scan()
		keyboardLine := scanner.Text()

		for j := 0; j < 4; j++ {
			n, err := strconv.Atoi(keyboardLine[j : j+1])
			if err == nil {
				keyboardCharacters[n]++
			}
		}
	}

	fmt.Println(countScore(k, keyboardCharacters))
}

func countScore(k int, c map[int]int) int {
	var res int
	for i := 1; i < 10; i++ {
		n, ok := c[i]
		if ok && n <= k*2 {
			res++
		}
	}
	return res
}
