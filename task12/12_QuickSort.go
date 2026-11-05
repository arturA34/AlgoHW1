package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Participant struct {
	name string
	p    int
	f    int
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	s.Scan()
	n, _ := strconv.Atoi(s.Text())

	participants := make([]Participant, 0, n+1)

	for i := 0; i < n; i++ {
		s.Scan()
		participantInfo := strings.Split(s.Text(), " ")
		p, _ := strconv.Atoi(participantInfo[1])
		f, _ := strconv.Atoi(participantInfo[2])

		participants = append(participants, Participant{name: participantInfo[0], p: p, f: f})
	}

	quickSort(participants, 0, len(participants)-1)

	writer := bufio.NewWriter(os.Stdout)
	for _, p := range participants {
		fmt.Fprintln(writer, p.name)
	}
	writer.Flush()
}

func quickSort(participants []Participant, left int, right int) {
	pivot := participants[(left+right)/2]

	i, j := left, right

	for i <= j {
		for ; isBetter(participants[i], pivot); i++ {
		}
		for ; isBetter(pivot, participants[j]); j-- {
		}

		if i <= j {
			participants[i], participants[j] = participants[j], participants[i]
			i++
			j--
		}
	}

	if left < j {
		quickSort(participants, left, j)
	}
	if i < right {
		quickSort(participants, i, right)
	}
}

func isBetter(p1, p2 Participant) bool {
	if p1.p > p2.p {
		return true
	} else if p1.p == p2.p {
		if p1.f < p2.f {
			return true
		} else if p1.f == p2.f {
			return p1.name < p2.name
		}
		return false
	}
	return false
}
