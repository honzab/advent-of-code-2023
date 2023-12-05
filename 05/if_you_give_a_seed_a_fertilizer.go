package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseInt64(s string) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Panicf("Could not parse int %v", s)
	}
	return num
}

func parseNextLines(scanner *bufio.Scanner, mapping map[int64][]Range, index int64) {
	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}
		line := scanner.Text()
		if len(line) == 0 {
			return
		}

		numberRegexp := regexp.MustCompile(`([0-9]+)`)
		numbers := numberRegexp.FindAllStringSubmatch(line, len(line))
		if len(numbers) != 3 {
			log.Panicf("Weird numbers found: %v", numbers)
		}
		from := parseInt64(numbers[0][1])
		to := parseInt64(numbers[1][1])
		length := parseInt64(numbers[2][1])

		mapping[index] = append(mapping[index], Range{from, to, length})
	}

}

type Range struct {
	From   int64
	To     int64
	Length int64
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	mapping := map[int64][]Range{}

	numberRegexp := regexp.MustCompile(`([0-9]+)`)

	seeds := make([]int64, 0)

	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}

		line := scanner.Text()

		if strings.Contains(line, "seeds") {
			for _, num := range numberRegexp.FindAllStringSubmatch(line, len(line)) {
				seeds = append(seeds, parseInt64(num[0]))
			}
			continue
		}

		switch line {
		case "seed-to-soil map:":
			parseNextLines(scanner, mapping, 1)
		case "soil-to-fertilizer map:":
			parseNextLines(scanner, mapping, 2)
		case "fertilizer-to-water map:":
			parseNextLines(scanner, mapping, 3)
		case "water-to-light map:":
			parseNextLines(scanner, mapping, 4)
		case "light-to-temperature map:":
			parseNextLines(scanner, mapping, 5)
		case "temperature-to-humidity map:":
			parseNextLines(scanner, mapping, 6)
		case "humidity-to-location map:":
			parseNextLines(scanner, mapping, 7)
		}
	}

	shortestDistance := int64(0)
	for i, seedId := range seeds {
		e := walkTroughForSeed(seedId, mapping)
		if i == 0 || e < shortestDistance {
			shortestDistance = e
		}
	}
	log.Printf("%d", shortestDistance)

	shortestDistance = int64(-1)
	for q := int64(0); q < int64(len(seeds)); q += 2 {
		seedId := seeds[q]
		length := seeds[q+1]
		for j := seedId; j < seedId+length; j++ {
			e := walkTroughForSeed(j, mapping)
			if shortestDistance == -1 || e < shortestDistance {
				shortestDistance = e
			}
		}
	}
	log.Printf("%d", shortestDistance)
}

func walkTroughForSeed(seedId int64, mapping map[int64][]Range) int64 {
	curentNodeId := seedId
	for i := int64(1); i < 8; i++ {
		for _, mappingRange := range mapping[i] {
			if mappingRange.To <= curentNodeId && curentNodeId < mappingRange.To+mappingRange.Length {
				diff := mappingRange.From - mappingRange.To
				curentNodeId += diff
				break
			}
		}
	}
	return curentNodeId
}
