package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const MUL_REGEX = `mul\((?<first>\d+),(?<second>\d+)\)`
const DO_REGEX = `do\(\)`
const DONT_REGEX = `dont'\(\)`

func readFile(fileName string) string {
	contents, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return string(contents)
}

// A mul() call in the file, plus the leftmost index of where it appears
type Multiplication struct {
	// The index of the contents in which this multiplication appears
	index int
	// The first multiplication value
	first int
	// The second multiplication value
	second int
}

func parseSubMatch(subMatch []int, contents string) Multiplication {
	first_index_start := subMatch[2]
	first_index_end := subMatch[3]
	second_index_start := subMatch[4]
	second_index_end := subMatch[5]

	first_mul := contents[first_index_start:first_index_end]
	second_mul := contents[second_index_start:second_index_end]
	fmt.Println(first_mul, second_mul)

	// The subMatch is guaranteed to return a valid integer due to the regex
	first, err := strconv.Atoi(first_mul)

	// Terrible bad logic, I am shocking at Go
	if err != nil {
		return Multiplication{first: 0, second: 0}
	}

	second, err := strconv.Atoi(second_mul)

	if err != nil {
		return Multiplication{first: 0, second: 0}
	}

	return Multiplication{
		index:  0,
		first:  first,
		second: second,
	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("Usage: <executable> <input_file_name>")
	}

	regex := regexp.MustCompile(MUL_REGEX)
	do_regex := regexp.MustCompile(DO_REGEX)
	dont_regex := regexp.MustCompile(DONT_REGEX)

	fileName := args[1]
	contents := readFile(fileName)

	matches := regex.FindAllStringSubmatchIndex(contents, -1)
	doMatches := do_regex.FindAllStringIndex(contents, -1)
	dontMatches := dont_regex.FindAllStringIndex(contents, -1)

	total := 0

	for _, subMatch := range matches {
		// Find all multiplication values
		fmt.Println(subMatch)
		multiplication := parseSubMatch(subMatch, contents)
		total += multiplication.first * multiplication.second
	}

	fmt.Println(total)
	fmt.Println(doMatches)
	fmt.Println(dontMatches)

}
