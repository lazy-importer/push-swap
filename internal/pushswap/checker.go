package pushswap

import (
	"bufio"
	"io"
	"strings"
)

// Check executes the instructions read from r on a copy of numbers and
// reports whether stack a is sorted and stack b is empty at the end.
//
// An unknown or badly formatted instruction makes Check return an error.
func Check(numbers []int, r io.Reader) (bool, error) {
	s := New(numbers)

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if err := s.Apply(line); err != nil {
			return false, err
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	return s.Done(), nil
}
