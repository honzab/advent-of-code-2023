package main

import (
	"bufio"
	"log"
	"regexp"
	"strconv"
	"strings"
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

type Card struct {
	Id          uint64
	Scratched   []uint64
	NumbersHave []uint64
}

func inSlice(needle uint64, haystack []uint64) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cards := make([]Card, 0)

	idRegex := regexp.MustCompile(`Card[ ]+([0-9]+):`)
	numberRegexp := regexp.MustCompile(`([0-9]+)`)

	// Parse into memory
	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}

		line := scanner.Text()

		id := idRegex.FindAllStringSubmatch(line, 1)

		numbers := numberRegexp.FindAllStringSubmatch(strings.Split(line, "|")[0], len(line))
		scratched := make([]uint64, 0)
		for _, n := range numbers {
			scratched = append(scratched, parseUint64(n[1]))
		}

		numbers = numberRegexp.FindAllStringSubmatch(strings.Split(line, "|")[1], len(line))
		have := make([]uint64, 0)
		for _, n := range numbers {
			have = append(have, parseUint64(n[1]))
		}

		cards = append(cards, Card{Id: parseUint64(id[0][1]), Scratched: scratched[1:], NumbersHave: have})
	}

	sum := uint64(0)
	for _, card := range cards {
		cardValue := 0
		for _, number := range card.Scratched {
			if inSlice(number, card.NumbersHave) {
				if cardValue == 0 {
					cardValue = 1
				} else {
					cardValue = cardValue * 2
				}
			}
		}
		sum += uint64(cardValue)
	}

	log.Printf("%d\n", sum)

}
