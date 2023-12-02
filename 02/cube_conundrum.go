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

type Game struct {
	Id    uint64
	Draws []Draw
}

type Draw struct {
	Blue  uint64
	Green uint64
	Red   uint64
}

func parseUint64(s string) uint64 {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Panicf("Could not parse uint %v", s)
	}
	return num
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	games := make([]Game, 0)

	gameIdRegexp := regexp.MustCompile(`Game\ ([0-9]+): (.*)`)
	gameDrawRegexp := regexp.MustCompile(`([0-9]+)\ (red|green|blue)`)

	// Parse into memory
	for scanner.Scan() {
		err := scanner.Err()
		if err == io.EOF {
			return
		}
		inputLine := scanner.Text()

		gameIdMatches := gameIdRegexp.FindStringSubmatch(inputLine)
		if len(gameIdMatches) != 3 {
			log.Panicf("%v did not match", inputLine)
		}
		gameId := parseUint64(gameIdMatches[1])

		gameDrawInputs := strings.Split(gameIdMatches[2], ";")
		game := Game{gameId, make([]Draw, 0)}

		for _, v := range gameDrawInputs {
			draws := gameDrawRegexp.FindAllStringSubmatch(v, len(v))
			draw := Draw{0, 0, 0}
			for _, b := range draws {
				if string(b[2]) == "green" {
					draw.Green = parseUint64(string(b[1]))
				} else if string(b[2]) == "red" {
					draw.Red = parseUint64(string(b[1]))
				} else if string(b[2]) == "blue" {
					draw.Blue = parseUint64(string(b[1]))
				} else {
					log.Panicf("Could not understand color %v", string(b[2]))
				}
			}
			game.Draws = append(game.Draws, draw)
		}
		games = append(games, game)
	}

	// Find possible games
	haveReds := uint64(12)
	haveGreens := uint64(13)
	haveBlues := uint64(14)
	sumIds := uint64(0)
	sumPowers := uint64(0)
	for _, g := range games {
		possible := true

		minReds := uint64(0)
		minGreens := uint64(0)
		minBlues := uint64(0)

		for _, d := range g.Draws {
			if d.Green > haveGreens || d.Red > haveReds || d.Blue > haveBlues {
				possible = false
			}
			minReds = max(minReds, d.Red)
			minGreens = max(minGreens, d.Green)
			minBlues = max(minBlues, d.Blue)
		}
		if possible {
			sumIds += g.Id
		}
		power := minBlues * minReds * minGreens
		sumPowers += power
		log.Printf("Checking game %d: possible %v, power: %d", g.Id, possible, power)
	}
	log.Printf("%d\n", sumIds)
	log.Printf("%d\n", sumPowers)
}
