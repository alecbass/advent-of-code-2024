package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Block struct {
	id     int64
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
			id:     int64(id),
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

// Swaps two sub-slices within a slice
//
// Constraints: start_end_index <= end_start_index
func swapBlockSlices(blocks []*Block, end_start_index int, end_end_index int, start_start_index int, start_end_index int) []*Block {
	// fmt.Println(start_start_index, start_end_index, end_start_index, end_end_index)

	// Elements before the first slice
	start := blocks[:start_start_index]

	// The first slice
	first_slice := blocks[start_start_index:start_end_index]

	// Elements between the first and second slices
	middle := blocks[start_end_index:end_start_index]

	// The second slice
	second_slice := blocks[end_start_index:end_end_index]

	// Elements after the second slice
	end := blocks[end_end_index:]

	// Swap the first and second slices
	return slices.Concat(start, second_slice, middle, first_slice, end)
}

func main() {
	input := readFile("input.txt")
	input = INPUT

	blocks := parseBlocks(input)
	length := len(blocks)

	for outer_index := length; outer_index > 0; {
		// fmt.Printf("Outer index: %d/%d\n", outer_index, length)
		outer_block := blocks[outer_index-1]

		if outer_block.is_free_space() {
			// Don't swap free spaces, search for the next block
			outer_index -= 1
			continue
		}

		// How long this block sequence is
		block_sequence_length := outer_block.length

		// At what point in the disk does this block sequence ")start"?
		file_start_index := outer_index - block_sequence_length

		// At what point in the disk does this block sequence "finish"?
		file_end_index := outer_index

		if file_start_index < 0 {
			// Can't shuffle this block back into the disk any more
			break
		}

		// PART TWO: Find out if this entire block can be shifted leftwards to avoid fragmentation
		// The left-most known position of where to start searching for free spaces on the disk
		free_space_start_index := 0

		for ; free_space_start_index < file_start_index; free_space_start_index++ {
			// Search left-to-right through the disk and see if we can swap this entire block sequence into the first sequqnce of free blocks that are large enough
			free_space_end_index := free_space_start_index + block_sequence_length

			// Would moving this block of free space exceed the lenght of the disk?
			is_too_large := free_space_end_index >= file_start_index

			if is_too_large {
				// No more space left to swap
				break
			}

			// Test the next n elements and swap this file if they're all free
			potential_swap_slice := blocks[free_space_start_index:free_space_end_index]
			are_all_spaces_free := true

			for _, swap_block := range potential_swap_slice {
				if !swap_block.is_free_space() {
					are_all_spaces_free = false
					break
				}
			}

			if !are_all_spaces_free {
				// This tested space isn't all free, so we don't have enough space to swap the files into here
				continue
			}

			// We can swap this block sequence slice to earlier in the disk
			blocks = swapBlockSlices(
				blocks,
				file_start_index,
				file_end_index,
				free_space_start_index,
				free_space_end_index,
			)

			break
		}

		could_not_swap := free_space_start_index >= file_start_index

		if could_not_swap {
			fmt.Printf("Could not swap block %d back\n", outer_block.id)
		}

		// Either we've swapped the file, or we can't swap it, move the file index back
		outer_index = file_start_index

		printDisk(blocks[0:42])
	}

	printDisk(blocks)
	fmt.Println("00992111777.44.333....5555.6666.....8888..", "EXPECTED")

	// Calculate the checksum
	var checksum int64 = 0

	// Assume there are no free spaces, so expect every character to be a digit
	for index, block := range blocks {
		if block.is_free_space() {
			// Skip checksum of free spaces
			continue
		}

		checksum += block.id * int64(index)
	}

	printDisk(blocks)
	fmt.Printf("Checksum: %d\n", checksum)
}
