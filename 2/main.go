package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isSafe(report []int) bool {
	if len(report) == 1 {
		return true
	}

	is_ascending := report[0] < report[1]
	is_dampened := false

	for i := 1; i < len(report); i++ {
		/** The previous level */
		var previous int

		if is_dampened {
			previous = report[i-2]
		} else {
			previous = report[i-1]
		}

		// Re-record the ascending/descending status if the second level was dampened
		if is_dampened && i == 2 {
			is_ascending = report[0] < report[2]
		}

		/** The current level */
		current := report[i]

		var difference int

		if is_ascending {
			difference = current - previous
		} else {
			difference = previous - current
		}

		if difference < 1 || difference > 3 {
			// Unsafe reading
			if !is_dampened {
				// Mitigate this instance
				is_dampened = true
				continue
			}

			return false
		}
	}

	return true
}

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("Usage: <program> <input_name>")
	}

	input_name := args[1]

	file_bytes, err := os.ReadFile(input_name)

	if err != nil {
		panic(err)
	}

	file_contents := string(file_bytes)
	lines := strings.Split(file_contents, "\n")

	safe_report_count := 0

	for _, line := range lines {
		if line == "" {
			continue
		}

		level_strs := strings.Split(line, " ")

		if len(level_strs) == 1 {
			// Only having one reading will always be safe
			safe_report_count++
			continue
		}

		/** All level readings for this report */
		report := []int{}

		for _, level_str := range level_strs {
			// The difference between the previous and current levels is less than one, or greater than three
			level, err := strconv.Atoi(level_str)

			if err != nil {
				continue
			}

			report = append(report, level)
		}

		/** Whether this report is safe or not */
		is_safe := isSafe(report)

		if is_safe {
			safe_report_count++
		}

		if !is_safe {
			fmt.Println(line, is_safe)
		}
	}

	fmt.Printf("Safe report count: %d\n", safe_report_count)
}
