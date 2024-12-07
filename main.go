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

// Stolen from Rust :)
type Result[T any] struct {
	value T
	err   error
}

func (r Result[T]) Ok() bool {
	return r.err == nil
}

func main() {
	//day_1()
	day_2()
}

func day_2() {
	data, err := os.ReadFile("./2/input.txt")
	check(err)
	lines := strings.Split(string(data), "\n")

	output_channel := make(chan Result[bool], len(lines))

	for _, line := range lines {
		report_numbers := parse_report(line)
		is_report_safe_with_combinations(report_numbers, output_channel)
	}

	num_safe := 0
	for range len(lines) {
		result := <-output_channel
		if result.Ok() && result.value {
			num_safe += 1
		}
	}

	fmt.Println(num_safe)
}

func parse_report(report string) []int {
	var parts = strings.Split(strings.TrimSpace(report), " ")
	var numbers = make([]int, 0)

	for i := range parts {
		num, err := strconv.Atoi(parts[i])

		// Bad data.
		if err != nil {
			continue
		}

		numbers = append(numbers, num)
	}
	return numbers
}

func generate_report_combinations(report []int) [][]int {
	// Make combinations of this report, each with 1 item missing.
	combinations := make([][]int, 0)

	for idx := range report {
		// Append a combination but removing the given index.
		combo := slices.Delete(slices.Clone(report), idx, idx+1)
		combinations = append(combinations, combo)
	}
	return combinations
}

func is_report_safe_with_combinations(report []int, output chan Result[bool]) {
	combinations := generate_report_combinations(report)
	results := make([]bool, 0)
	for _, combination := range combinations {
		result := is_safe(combination)

		if result.Ok() {
			results = append(results, result.value)
		}
	}

	// Output is true if any were true
	// i.e. any combination returned true.
	output <- Result[bool]{
		err:   nil,
		value: slices.Contains(results, true),
	}
}

func is_safe(report []int) Result[bool] {
	differences := make([]int, 0)

	for idx := range len(report) - 1 {
		differences = append(differences, report[idx+1]-report[idx])
	}

	// Ensure no zeros (i.e. not ascending or descending)
	if slices.Contains(differences, 0) {
		return Result[bool]{
			err:   nil,
			value: false,
		}
	}

	// Ensure all differences in the expected range.
	for _, difference := range differences {
		if Abs(difference) < 1 || Abs(difference) > 3 {
			return Result[bool]{
				err:   nil,
				value: false,
			}
		}
	}

	// Ensure ascendeness and descendedness is ok
	for idx := range len(differences) - 1 {
		if sign(differences[idx]) != sign(differences[idx+1]) {
			return Result[bool]{
				err:   nil,
				value: false,
			}
		}
	}

	// All checks passed, good!
	return Result[bool]{
		err:   nil,
		value: true,
	}
}

func sign(v int) int {
	if v < 0 {
		return -1
	} else if v > 0 {
		return 1
	} else {
		return 0
	}
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
