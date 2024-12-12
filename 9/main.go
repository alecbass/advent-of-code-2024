package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Block struct {
	id      int
	length  int
	is_file bool
}

const INPUT = "2333133121414131402"

func readFile(fileName string) string {
	contents, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	return strings.ReplaceAll(string(contents), "\n", "")
}

func main() {
	input := readFile("input.txt")
	input = INPUT

	next_id := 0

	// All disk blocks
	blocks := []Block{}

	// How many blocks are files? Used to make the swapping search faster
	file_block_count := 0

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
		block := Block{
			id:      id,
			length:  length,
			is_file: is_file,
		}

		if is_file {
			// Prepare the next ID
			next_id += 1

			// Increase how many file blocks we have stored
			file_block_count += length
		}

		blocks = append(blocks, block)

	}

	var builder strings.Builder

	for _, block := range blocks {
		is_file := block.id != -1

		if is_file {
			for i := 0; i < block.length; i++ {
				builder.WriteString(fmt.Sprintf("%d", block.id))
			}
		} else {
			for i := 0; i < block.length; i++ {
				builder.WriteString(".")
			}
		}
	}

	// Move each element from the end to the first free space on the left

	disk := []rune(builder.String())

	for outer := len(disk) - 1; outer >= 0; outer-- {
		is_free_space := disk[outer] == '.'

		if is_free_space {
			// Don't swap free spaces
			continue
		}

		// TODO: A rogue '.' is being swapped into the start of the end result
		// existing := strings.Clone(string(disk))

		// Find the first free space and replace it with this file block
		for inner := 0; inner < file_block_count; inner++ {
			is_file := disk[inner] != '.'

			if !is_file {
				// Space is free, copy this file block to earlier in the disk
				disk[inner] = disk[outer]
				break
			}

			// Swap the two characters
			// temp := disk[inner]
			// disk[inner] = disk[outer]
			// disk[outer] = temp
		}

		// fmt.Println(existing)
		// fmt.Println(string(disk))
	}

	fmt.Println(len(disk), file_block_count)

	// Since all file blocks have replaced and early free spaces, simply cut off the now-redundant file blocks
	disk = disk[0:file_block_count]
	fmt.Println(len(disk))

	// Forcefully rotate the array leftwards if the first element is a free space
	// Assume that the last two spaces are free
	// I am dumb
	// is_start_free_space := disk[0] == '.'
	//
	// if is_start_free_space {
	// 	for i := 1; i < length; i++ {
	// 		previous := disk[i-1]
	// 		current := disk[i]
	//
	// 		// Swap the elements
	// 		disk[i] = previous
	// 		disk[i-1] = current
	// 	}
	// }

	// fmt.Println(string(disk))
	// fmt.Println("0099811188827773336446555566")

	// Calculate the checksum
	checksum := 0

	// Assume there are no free spaces, so expect every character to be a digit
	for i := 0; i < len(disk); i++ {
		value, err := strconv.Atoi(string(disk[i]))

		if err != nil {
			panic("Failed to parse character when calculating the checksum")
		}

		checksum += value * i
	}

	fmt.Printf("Checksum: %d\n", checksum)
}
