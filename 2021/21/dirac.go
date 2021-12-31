package d21

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type player struct {
	pos   int
	score int
}

type game struct {
	player1, player2 player
	die              func() int
	rollCount        int
}

type quantumGame struct {
	// these map a player struct to the # of universes where that struct could exist
	player1 map[player]uint
	player2 map[player]uint

	player1Wins uint
	player2Wins uint
}

const boardSize = 10

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(21, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	game := parseInput(r)
	game.die = createDeterministicDie(100)

	_, loser := game.run(1000)

	result := loser.score * game.rollCount

	return strconv.Itoa(result)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	game := parseInput(r)

	qgame := quantumGame{
		player1:     make(map[player]uint),
		player2:     make(map[player]uint),
		player1Wins: 0,
		player2Wins: 0,
	}

	// players start the game on their starting position w/ 0 score in 1 universe
	qgame.player1[game.player1] = 1
	qgame.player2[game.player2] = 1

	qgame.run()

	if qgame.player1Wins > qgame.player2Wins {
		return strconv.FormatUint(uint64(qgame.player1Wins), 10)
	} else if qgame.player2Wins > qgame.player1Wins {
		return strconv.FormatUint(uint64(qgame.player2Wins), 10)
	} else {
		panic("they tied?")
	}
}

////////////////////////////////////////////////////////////////////////////////
// quantum game

func (q *quantumGame) run() {
	const maxScore = 21

	// maps move distance (key) to the # of universes in which the player would get that distance
	moves := make(map[int]uint)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				moves[i+j+k]++
			}
		}
	}

	for {

		q.player1Wins += calculateQuantumWins(moves, q.player1)

		q.player2Wins += calculateQuantumWins(moves, q.player2)

		fmt.Printf("%d vs %d\n", q.player1Wins, q.player2Wins)

		q.player1 = quantumStep(moves, q.player1)

		q.player2 = quantumStep(moves, q.player2)

		var totalUniverses uint
		for _, universes := range q.player1 {
			totalUniverses += universes
		}

		fmt.Printf("total: %d\n", totalUniverses)

		if q.player1Wins+q.player2Wins > totalUniverses {
			break
		}
	}
}

// returns the number of universes in which the given moves would move the
// player's score >= 21
func calculateQuantumWins(moves map[int]uint, playerStates map[player]uint) uint {

	var wins uint

	for moveDistance, moveUniverses := range moves {

		for pos := 1; pos <= 10; pos++ {

			scoreBoost := pos + moveDistance
			if scoreBoost > 10 {
				scoreBoost = scoreBoost % 10
			}

			// find all player states such that state.pos == pos && state.score + scoreBoost > 21
			for p, playerUniverses := range playerStates {

				if p.pos != pos {
					continue
				}

				if p.score >= 21 {
					continue
				}

				if p.score+scoreBoost < 21 {
					continue
				}

				wins += playerUniverses * moveUniverses
			}
		}
	}

	return wins
}

func quantumStep(moves map[int]uint, playerStates map[player]uint) map[player]uint {

	next := make(map[player]uint)

	for distance, moveUniverses := range moves {
		for p, playerUniverses := range playerStates {
			newPos := p.pos + distance
			if newPos > 10 {
				newPos = newPos % 10
			}
			nextPlayer := player{
				pos:   newPos,
				score: p.score + newPos,
			}
			next[nextPlayer] += (playerUniverses * moveUniverses)
		}
	}

	return next
}

////////////////////////////////////////////////////////////////////////////////
// classical game

func (p *player) move(steps int, boardSize int) {
	p.pos += steps
	for p.pos > boardSize {
		p.pos -= boardSize
	}
	p.score += p.pos
}

func (g *game) roll(times int) int {
	var result int
	for i := 0; i < times; i++ {
		g.rollCount++
		result += g.die()
	}
	return result
}

func (g *game) run(winningScore int) (*player, *player) {

	for {

		move := g.roll(3)
		g.player1.move(move, boardSize)
		if g.player1.score >= winningScore {
			return &g.player1, &g.player2
		}

		move = g.roll(3)
		g.player2.move(move, boardSize)
		if g.player2.score >= winningScore {
			return &g.player2, &g.player1
		}
	}
}

func createDeterministicDie(sides int) func() int {
	last := 0
	return func() int {
		if last == sides {
			last = 0
		}
		last++
		return last
	}

}

func parseInput(r io.Reader) game {
	s := bufio.NewScanner(r)

	rx := regexp.MustCompile("Player (\\d+) starting position: (\\d+)")

	positions := [2]int{0, 0}

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if len(line) == 0 {
			continue
		}

		m := rx.FindStringSubmatch(line)

		if m == nil {
			continue
		}

		player, err := strconv.ParseInt(m[1], 10, 16)
		if err != nil {
			continue
		}

		pos, err := strconv.ParseInt(m[2], 10, 16)
		if err != nil {
			continue
		}

		positions[player-1] = int(pos)
	}

	return game{
		player1: player{pos: positions[0]},
		player2: player{pos: positions[1]},
	}
}
