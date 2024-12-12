package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Block struct {
	id     int
	length int
}

func (block Block) is_free_space() bool {
	return block.id == -1
}

const INPUT = "2333133121414131402"

func readFile(fileName string) string {
	contents, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(string(contents), "\n", "")
}

func parseBlocks(input string) []*Block {
	next_id := 0

	// All disk blocks
	blocks := []*Block{}

	for index, digit := range input {
		length, err := strconv.Atoi(string(digit))

		is_file := index%2 == 0

		if err != nil {
			panic("Failed to parse digit")
		}

		var id int

		if is_file {
			id = next_id
		} else {
			// Use -1 to denote whether a block is a file or empty space
			id = -1
		}

		// Record this disk block
		block := &Block{
			id:     id,
			length: length,
		}

		if is_file {
			// Prepare the next ID
			next_id += 1
		}

		for i := 0; i < length; i++ {
			// Add a pointer to this individual space for each length of the block
			blocks = append(blocks, block)
		}

	}

	return blocks
}

func printDisk(blocks []*Block) {
	for _, block := range blocks {
		if block.is_free_space() {
			fmt.Print(".")
		} else {
			fmt.Printf("%d", block.id)
		}
	}

	fmt.Printf("\n")
}

func main() {
	input := readFile("input.txt")
	// input = INPUT

	blocks := parseBlocks(input)
	length := len(blocks)

	for outer := length - 1; outer >= 0; outer-- {
		outer_block := blocks[outer]

		is_free_space := outer_block.is_free_space()

		if is_free_space {
			// Don't swap free spaces
			continue
		}

		for inner := 0; inner < length; inner++ {
			// Move each element from the end to the first free space on the left
			inner_block := blocks[inner]

			if inner_block.is_free_space() {
				// Swap this free space with a pointer to a real block
				blocks[inner] = outer_block
				blocks[outer] = inner_block
				break
			}
		}
	}

	// Rotate the array if it starts with a free space because I am dumb
	for i := 1; i < length; i++ {
		previous := blocks[i-1]
		current := blocks[i]

		should_swap := previous.is_free_space() && !current.is_free_space()

		if !should_swap {
			// Can't swap any more
			break
		}

		blocks[i] = previous
		blocks[i-1] = current
	}

	// Calculate the checksum
	checksum := 0

	// Assume there are no free spaces, so expect every character to be a digit
	for index, block := range blocks {
		if block.is_free_space() {
			// We've reached free spaces, stop increasing checksum
			break
		}

		checksum += block.id * index
	}

	printDisk(blocks)
	fmt.Printf("Checksum: %d\n", checksum)
}
