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

func parseUint64(s string) uint64 {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Panicf("Could not parse uint %v", s)
	}
	return num
}

type Race struct {
	Time   uint64
	Record uint64
}

func findWaysToRace(r *Race) uint64 {
	winningLengths := uint64(0)
	for t := uint64(0); t <= r.Time; t++ {
		distance := t * (r.Time - t)
		if distance > r.Record {
			winningLengths += 1
		}
	}
	return winningLengths
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	numberRegexp := regexp.MustCompile(`([0-9]+)`)
	races := make([]Race, 0)

	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}

		line := scanner.Text()

		if strings.Contains(line, "Time:") {
			times := numberRegexp.FindAllStringSubmatch(strings.Split(line, ":")[1], len(line))
			for _, t := range times {
				races = append(races, Race{parseUint64(t[1]), 0})
			}
		}
		if strings.Contains(line, "Distance:") {
			distances := numberRegexp.FindAllStringSubmatch(strings.Split(line, ":")[1], len(line))
			for i, d := range distances {
				races[i] = Race{races[i].Time, parseUint64(d[1])}
			}
		}
	}

	result := uint64(1)
	for _, r := range races {
		result = result * findWaysToRace(&r)
	}

	log.Printf("%d\n", result)

}
