package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
)
import "os"
import "io"

func parseUint64(s string) uint64 {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Panicf("Could not parse uint %v", s)
	}
	return num
}

var SCHEMATICS []string

func checkStringForSymbols(line string) bool {
	symbol := regexp.MustCompile(`^[0-9\.]+$`)
	matching := symbol.MatchString(line)
	//log.Printf("Checking line %s: %v\n", line, !matching)
	return !matching
}

func isAdjacentToSymbol(lineNumber, index, length int) uint64 {
	n := parseUint64(SCHEMATICS[lineNumber][index:length])

	leftBound := max(0, index-1)
	rightBound := min(len(SCHEMATICS[lineNumber]), length+1)
	if lineNumber-1 >= 0 {
		if checkStringForSymbols(SCHEMATICS[lineNumber-1][leftBound:rightBound]) {
			return n
		}
	}
	if index > 0 {
		if checkStringForSymbols(SCHEMATICS[lineNumber][index-1 : index]) {
			return n
		}
	}
	if length+1 <= len(SCHEMATICS[lineNumber]) {
		if checkStringForSymbols(SCHEMATICS[lineNumber][length : length+1]) {
			return n
		}
	}
	if lineNumber+1 < len(SCHEMATICS) {
		if checkStringForSymbols(SCHEMATICS[lineNumber+1][leftBound:rightBound]) {
			return n
		}
	}
	return 0
}

var ToMultiply map[string][]uint64

func checkStringForStars(line string, lineIndex, lb, rb int, n uint64) bool {
	symbol := regexp.MustCompile(`[*]`)
	matching := symbol.FindAllStringIndex(line, len(line))

	for _, m := range matching {
		key := fmt.Sprintf("%d-%v", lineIndex, m[0]+lb)
		if ToMultiply[key] == nil {
			ToMultiply[key] = []uint64{n}
		} else {
			ToMultiply[key] = append(ToMultiply[key], n)
		}
	}
	return len(matching) > 0
}

func isAdjacentToAStar(lineNumber, index, length int) uint64 {
	n := parseUint64(SCHEMATICS[lineNumber][index:length])

	leftBound := max(0, index-1)
	rightBound := min(len(SCHEMATICS[lineNumber]), length+1)
	if lineNumber-1 >= 0 {
		if checkStringForStars(SCHEMATICS[lineNumber-1][leftBound:rightBound], lineNumber-1, leftBound, rightBound, n) {
			return n
		}
	}
	if index > 0 {
		if checkStringForStars(SCHEMATICS[lineNumber][index-1:index], lineNumber, index-1, index, n) {
			return n
		}
	}
	if length+1 <= len(SCHEMATICS[lineNumber]) {
		if checkStringForStars(SCHEMATICS[lineNumber][length:length+1], lineNumber, length, length+1, n) {
			return n
		}
	}
	if lineNumber+1 < len(SCHEMATICS) {
		if checkStringForStars(SCHEMATICS[lineNumber+1][leftBound:rightBound], lineNumber+1, leftBound, rightBound, n) {
			return n
		}
	}
	return 0
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	ToMultiply = make(map[string][]uint64)

	numberRegexp := regexp.MustCompile(`[0-9]+`)

	// Parse into memory
	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}
		SCHEMATICS = append(SCHEMATICS, scanner.Text())
	}

	sum := uint64(0)
	for i, line := range SCHEMATICS {
		numbers := numberRegexp.FindAllStringIndex(line, len(line))
		for _, n := range numbers {
			sum += isAdjacentToSymbol(i, n[0], n[1])
		}
	}
	log.Printf("%d\n", sum)

	for i, line := range SCHEMATICS {
		numbers := numberRegexp.FindAllStringIndex(line, len(line))
		for _, n := range numbers {
			isAdjacentToAStar(i, n[0], n[1])
		}
	}

	gearSum := uint64(0)
	for _, v := range ToMultiply {
		if len(v) > 1 {
			product := uint64(1)
			for _, i := range v {
				product = product * i
			}
			gearSum += product
		}
	}
	log.Printf("%d\n", gearSum)
}
