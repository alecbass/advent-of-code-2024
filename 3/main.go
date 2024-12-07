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

type Multiplication struct {
	first  int
	second int
}

func parseSubMatch(subMatch []string) Multiplication {
	// The subMatch is guaranteed to return a valid integer due to the regex
	first, err := strconv.Atoi(subMatch[1])

	// Terrible bad logic, I am shocking at Go
	if err != nil {
		return Multiplication{first: 0, second: 0}
	}

	second, err := strconv.Atoi(subMatch[2])

	if err != nil {
		return Multiplication{first: 0, second: 0}
	}

	return Multiplication{
		first,
		second,
	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("Usage: <executable> <input_file_name>")
	}

	regex := regexp.MustCompile(MUL_REGEX)

	fileName := args[1]
	contents := readFile(fileName)

	matches := regex.FindAllStringSubmatch(contents, -1)

	total := 0

	for _, subMatch := range matches {
		multiplication := parseSubMatch(subMatch)
		fmt.Println(multiplication)

		total += multiplication.first * multiplication.second
	}

	fmt.Println(total)

}
