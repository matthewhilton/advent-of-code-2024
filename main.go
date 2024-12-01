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
	// Using goroutines, for fun.
	data_as_string := string(data)
	lines := strings.Split(data_as_string, "\n")

	var channels []chan [2]int
	for _, line := range lines {
		channels = append(channels, get_pair(line))
	}

	var lefts []int
	var rights []int

	for _, c := range channels {
		v, ok := <-c

		if ok {
			lefts = append(lefts, v[0])
			rights = append(rights, v[1])
		}
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
	// With channels, just to be interesting.
	similarity := 0
	similarity_channel := make(chan int)
	for _, left := range lefts {
		go calculate_similarity(left, rights, similarity_channel)
	}

	for range len(lefts) {
		similarity += <-similarity_channel
	}
	close(similarity_channel)

	fmt.Println(similarity)
}

func get_pair(value string) chan [2]int {
	c := make(chan [2]int)

	go func() {
		parts := strings.Split(value, "   ")
		if len(parts) != 2 {
			close(c)
			return
		}

		left, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		right, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

		if err1 != nil || err2 != nil {
			close(c)
			return
		}

		c <- [2]int{left, right}
	}()

	return c
}

/*
Finds the number of ocurrences of value in data, and sends back the number to the given channel.
*/
func calculate_similarity(value int, data []int, c chan int) {
	occurrences := 0
	for _, data_value := range data {
		if value == data_value {
			occurrences += 1
		}
	}
	c <- occurrences * value
}
