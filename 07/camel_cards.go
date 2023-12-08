package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseUint64(s string) uint64 {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Panicf("Could not parse uint %v", s)
	}
	return num
}

type HandType int

const (
	FIVE_OF_A_KIND HandType = iota
	FOUR_OF_A_KIND
	FULL_HOUSE
	THREE_OF_A_KIND
	TWO_PAIR
	ONE_PAIR
	HIGH_CARD
)

const VALUES = "AKQJT98765432"

const VALUES_PART2 = "AKQT98765432J"

type Hand struct {
	Type HandType
	Repr string
	Bid  uint64
}

type HandsByValue []Hand

func (a HandsByValue) Len() int { return len(a) }
func (a HandsByValue) Less(i, j int) bool {
	if a[i].Type < a[j].Type {
		return true
	}
	for x := 0; x < len(a[i].Repr); x++ {
		if a[i].Repr[x] == a[j].Repr[x] {
			continue
		}
		return strings.Index(VALUES, string(a[i].Repr[x])) > strings.Index(VALUES, string(a[j].Repr[x]))
	}
	return false
}
func (a HandsByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type HandsByJokerValue []Hand

func (a HandsByJokerValue) Len() int { return len(a) }
func (a HandsByJokerValue) Less(i, j int) bool {
	if a[i].Type < a[j].Type {
		return true
	}
	for x := 0; x < len(a[i].Repr); x++ {
		if a[i].Repr[x] == a[j].Repr[x] {
			continue
		}
		return strings.Index(VALUES_PART2, string(a[i].Repr[x])) > strings.Index(VALUES_PART2, string(a[j].Repr[x]))
	}
	return false
}
func (a HandsByJokerValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func cardType(repr string) HandType {
	cardIncidence := make(map[string]uint64, 0)
	for i := 0; i < len(repr); i++ {
		cardIncidence[string(repr[i])] += 1
	}

	maxIncidence := uint64(0)
	for _, v := range cardIncidence {
		if v > maxIncidence {
			maxIncidence = v
		}
	}

	if maxIncidence == 5 {
		return FIVE_OF_A_KIND
	} else if maxIncidence == 4 {
		return FOUR_OF_A_KIND
	} else if maxIncidence == 3 && len(cardIncidence) == 2 {
		return FULL_HOUSE
	} else if maxIncidence == 3 {
		return THREE_OF_A_KIND
	} else if maxIncidence == 2 && len(cardIncidence) == 3 {
		return TWO_PAIR
	} else if maxIncidence == 2 {
		return ONE_PAIR
	}
	return HIGH_CARD
}

func jokerCardType(repr string) HandType {
	cardIncidence := make(map[string]uint64, 0)
	for i := 0; i < len(repr); i++ {
		cardIncidence[string(repr[i])] += 1
	}

	originalType := cardType(repr)

	jokers := cardIncidence["J"]
	if jokers == 0 {
		return originalType
	}

	// Cherry-pick "of a kinds"
	if originalType == THREE_OF_A_KIND && jokers == 1 {
		return FOUR_OF_A_KIND
	}
	if originalType == ONE_PAIR && jokers == 1 {
		return THREE_OF_A_KIND
	}
	if originalType == TWO_PAIR && jokers == 2 {
		return FOUR_OF_A_KIND
	}
	if originalType == TWO_PAIR && jokers == 1 {
		return FULL_HOUSE
	}

	return max(HandType(0), originalType-HandType(jokers))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	handTypes := make(map[HandType][]Hand)
	jokerHandTypes := make(map[HandType][]Hand)

	lines := 0
	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}

		line := scanner.Text()
		lines += 1

		splitLine := strings.Split(line, " ")

		ct := cardType(splitLine[0])
		if handTypes[ct] == nil {
			handTypes[ct] = make([]Hand, 0)
		}

		handTypes[ct] = append(handTypes[ct], Hand{ct, splitLine[0], parseUint64(splitLine[1])})

		jct := jokerCardType(splitLine[0])
		if jokerHandTypes[jct] == nil {
			jokerHandTypes[jct] = make([]Hand, 0)
		}
		jokerHandTypes[jct] = append(jokerHandTypes[jct], Hand{jct, splitLine[0], parseUint64(splitLine[1])})
	}

	rank := uint64(1)
	jokerRank := uint64(1)
	winnings := uint64(0)
	jokerWinnings := uint64(0)

	for v := 6; v >= 0; v-- {
		sort.Sort(HandsByValue(handTypes[HandType(v)]))
		for _, h := range handTypes[HandType(v)] {
			winnings += h.Bid * rank
			rank += 1
		}
		sort.Sort(HandsByJokerValue(jokerHandTypes[HandType(v)]))
		for _, h := range jokerHandTypes[HandType(v)] {
			jokerWinnings += h.Bid * jokerRank
			jokerRank += 1
		}
	}
	log.Printf("%v", winnings)
	log.Printf("%v", jokerWinnings)

}
