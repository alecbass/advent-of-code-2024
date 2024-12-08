package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const MUL_REGEX = `mul\((?<first>\d+),(?<second>\d+)\)`
const DO_REGEX = `do\(\)`
const DO_OR_DONT_REGEX = `(do\(\))|(don't\(\))`

// A mul() call in the file, plus the leftmost index of where it appears
type Multiplication struct {
	// The first multiplication value
	first int
	// The second multiplication value
	second int
}

func readFile(fileName string) string {
	contents, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return string(contents)
}

func parseSubMatch(sub_match []int, contents string) Multiplication {
	// Where this multiplication appears in the file
	first_index_start := sub_match[2]
	first_index_end := sub_match[3]
	second_index_start := sub_match[4]
	second_index_end := sub_match[5]

	first_mul := contents[first_index_start:first_index_end]
	second_mul := contents[second_index_start:second_index_end]

	// The subMatch is guaranteed to return a valid integer due to the regex
	first, err := strconv.Atoi(first_mul)

	// Terrible bad logic, I am shocking at Go
	if err != nil {
		panic("Invalid first number")
	}

	second, err := strconv.Atoi(second_mul)

	if err != nil {
		panic("Invalid second number")
	}

	return Multiplication{
		first:  first,
		second: second,
	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("Usage: <executable> <input_file_name>")
	}

	mul_regex := regexp.MustCompile(MUL_REGEX)
	do_regex := regexp.MustCompile(DO_REGEX)

	file_name := args[1]
	contents := readFile(file_name)

	do_matches := do_regex.FindAllStringIndex(contents, -1)

	// Find all index ranges where do() has been called, and dont() has not
	valid_ranges := [][]int{}

	terminator_regex := regexp.MustCompile(DO_OR_DONT_REGEX)

	for index, do_match := range do_matches {
		do_start_index := do_match[0]

		is_first := index == 0

		if is_first {
			// Since the list starts enabled, prepend with a valid range from the start do the first do() or don't() call
			valid_ranges = append(valid_ranges, []int{0, do_start_index})
		}

		// Check this substring, ignoring the beginning "do()"
		contents_slice := contents[do_start_index:]
		terminator := terminator_regex.FindStringIndex(contents_slice[4:])

		// Read to the end of the file by default
		terminator_index := len(contents_slice) - 1

		if len(terminator) > 1 {
			// Not at the end of the file, i.e. we found a do() or don't() call
			terminator_index = terminator[1]
		}

		valid_ranges = append(valid_ranges, []int{do_start_index + 4, do_start_index + terminator_index})
	}

	// The running multiplication total
	total := 0

	for _, r := range valid_ranges {
		start := r[0]
		end := r[1]

		content_slice := contents[start:end]
		range_mul_matches := mul_regex.FindAllStringSubmatchIndex(content_slice, -1)

		for _, mul_match := range range_mul_matches {
			// Find all mul() calls within this given range
			multiplication := parseSubMatch(mul_match, content_slice)
			total += multiplication.first * multiplication.second
		}
	}

	fmt.Printf("Total: %d\n", total)
}
