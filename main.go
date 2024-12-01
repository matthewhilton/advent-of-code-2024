package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// https://stackoverflow.com/a/76865734
func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	day_1()
}

/*
Two lists of numbers, pairs.
Sort each list by number, ascending
For each pair, get the distance between the two.
Then add all the distances
*/
func day_1() {
	data, err := os.ReadFile("./1/input.txt")
	check(err)

	// Split into pairs and remove any garbage.
	data_as_string := string(data)
	lines := strings.Split(data_as_string, "\n")

	var lefts []int
	var rights []int

	for _, line := range lines {
		parts := strings.Split(line, "   ")

		if len(parts) != 2 {
			continue
		}

		left, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		right, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

		if err1 != nil || err2 != nil {
			continue
		}

		lefts = append(lefts, left)
		rights = append(rights, right)
	}

	// Sort each side.
	slices.Sort(lefts)
	slices.Sort(rights)

	// Get the total
	total := 0
	for idx := range len(lefts) {
		total += Abs(lefts[idx] - rights[idx])
	}
	fmt.Println(total)

	// Part 2  - get similarity score.
	similarity := 0
	for _, left := range lefts {
		var occurrences = 0

		for _, right := range rights {
			if left == right {
				occurrences += 1
			}
		}

		similarity += occurrences * left
	}
	fmt.Println(similarity)
}
