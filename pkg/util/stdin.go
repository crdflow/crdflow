package util

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

// YesNo reads from stdin looking for one of "y", "yes", "n", "no" and returns
// true for "y" and false for "n"
func YesNo(reader *bufio.Reader) bool {
	for {
		text := readstdin(reader)
		switch text {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Printf("invalid input %q, should be [y/n]", text)
		}
	}
}

// Readstdin reads a line from stdin trimming spaces, and returns the value.
// log.Fatal's if there is an error.
func readstdin(reader *bufio.Reader) string {
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error when reading input: %v", err)
	}
	return strings.TrimSpace(text)
}
