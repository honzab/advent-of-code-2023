package main

import "bufio"
import "unicode"
import "fmt"
import "os"
import "io"
import "strconv"

func main() {
	WORDS := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	var sum uint64
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}
		value := scanner.Text()

		digits := make([]uint64, 0)
		for i, v := range value {
			if unicode.IsDigit(v) {
				value, err := strconv.ParseUint(string(v), 10, 64)
				if err != nil {
					return
				}
				digits = append(digits, value)
			} else {
				for j, w := range WORDS {
					if i+len(w) > len(value) {
						continue
					}
					if value[i:i+len(w)] == w {
						digits = append(digits, uint64(j+1))
					}
				}
			}
		}

		if len(digits) == 0 {
			continue
		} else if len(digits) == 1 {
			sum += 10*digits[0] + digits[0]
		} else {
			sum += 10*digits[0] + digits[len(digits)-1]
		}
		fmt.Printf("%s: %v\n", value, digits)
	}
	fmt.Printf("%d\n", sum)
}
