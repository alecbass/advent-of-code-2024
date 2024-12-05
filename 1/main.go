package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const INPUT_NAME = "input.txt"

func sortInts(a int, b int) int {
	return cmp.Compare(a, b)
}

func main() {
	file_bytes, err := os.ReadFile(INPUT_NAME)

	if err != nil {
		panic(err)
	}

	file_contents := string(file_bytes)
	lines := strings.Split(file_contents, "\n")

	left_column := []int{}
	right_column := []int{}

	// Mapping of right column values to how often they appear
	appearance_counts := make(map[int]int)

	for _, line := range lines {
		fmt.Println(line)
		values := strings.Split(string(line), "   ")

		if len(values) != 2 {
			continue
		}

		left, err := strconv.Atoi(values[0])

		if err != nil {
			continue
		}

		right, err := strconv.Atoi(values[1])

		if err != nil {
			continue
		}

		left_column = append(left_column, left)
		right_column = append(right_column, right)

		// Mark this value as having appeared
		appearance_counts[right] = appearance_counts[right] + 1
	}

	// Sort each list
	slices.SortFunc(left_column, sortInts)
	slices.SortFunc(right_column, sortInts)

	if len(left_column) != len(right_column) {
		panic("Column length mistmatch")
	}

	// Total distance between each value
	total_distance := 0

	// Mapping of left column values and how often they appear
	similarity_scores := make(map[int]int)

	for i := 0; i < len(left_column); i++ {
		left := left_column[i]
		right := right_column[i]

		var distance int

		if left > right {
			distance = left - right
		} else {
			distance = right - left
		}

		total_distance += distance

		// Record how often the left value appears in the right list, but only do it the first time this value has been seen
		appearance_count := appearance_counts[left]

		_, has_been_recorded := similarity_scores[left]

		if !has_been_recorded {
			similarity_scores[left] = left * appearance_count
		} else {
			fmt.Printf("Similarity score: %d\n", similarity_scores[left])
		}
	}

	total_similarity := 0

	for _, similarity_score := range similarity_scores {
		total_similarity += similarity_score
	}

	fmt.Printf("Total distance: %d\n", total_distance)
	fmt.Printf("Total similarity: %d\n", total_similarity)
}
