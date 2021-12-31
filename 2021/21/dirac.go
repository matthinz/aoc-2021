package d21

import (
	"bufio"
	_ "embed"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

// TODO: Implement part 1 using quantum games

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
	// this maps game states to the number of universes in which they exist
	states map[quantumGameState]uint
}

type quantumGameState struct {
	player1Pos, player2Pos     int
	player1Score, player2Score int
}

type quantumGameResult struct {
	// universes in which each player wins
	player1Wins, player2Wins uint
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

	result := runQuantumGame(
		game.player1.pos,
		game.player2.pos,
	)

	if result.player1Wins > result.player2Wins {
		return strconv.FormatUint(uint64(result.player1Wins), 10)
	} else if result.player2Wins > result.player1Wins {
		return strconv.FormatUint(uint64(result.player2Wins), 10)
	} else {
		return "tie"
	}
}

////////////////////////////////////////////////////////////////////////////////
// quantum game

func runQuantumGame(player1Pos, player2Pos int) quantumGameResult {

	const rollsPerTurn = 3
	const maxScore = 21

	// states maps a game state to the number of universes in which it exists
	states := make(map[quantumGameState]uint)

	// initially we have our starting game state, which exists in 1 universe
	initialState := quantumGameState{
		player1Pos: player1Pos,
		player2Pos: player2Pos,
	}
	states[initialState] = 1

	for {

		// Move the first player
		nextStates := make(map[quantumGameState]uint)
		moves := rollQuantumDie(3, 3)
		anyStillRunning := false

		for state, universes := range states {

			if state.player1Score >= maxScore || state.player2Score >= maxScore {
				nextStates[state] += universes
				continue
			}

			anyStillRunning = true

			for move, moveUniverses := range moves {
				nextPos := getNextPosition(state.player1Pos, move)
				nextState := quantumGameState{
					player1Pos:   getNextPosition(state.player1Pos, move),
					player1Score: state.player1Score + nextPos,
					player2Score: state.player2Score,
					player2Pos:   state.player2Pos,
				}
				nextStates[nextState] += universes * moveUniverses
			}
		}

		states = nextStates
		nextStates = make(map[quantumGameState]uint)

		if !anyStillRunning {
			break
		}

		nextStates = make(map[quantumGameState]uint)
		moves = rollQuantumDie(3, 3)
		anyStillRunning = false

		for state, universes := range states {

			if state.player1Score >= maxScore || state.player2Score >= maxScore {
				nextStates[state] += universes
				continue
			}

			anyStillRunning = true

			for move, moveUniverses := range moves {
				nextPos := getNextPosition(state.player2Pos, move)
				nextState := quantumGameState{
					player1Pos:   state.player1Pos,
					player1Score: state.player1Score,
					player2Pos:   nextPos,
					player2Score: state.player2Score + nextPos,
				}
				nextStates[nextState] += universes * moveUniverses
			}
		}

		states = nextStates

		if !anyStillRunning {
			break
		}
	}

	result := quantumGameResult{}

	for state, universes := range states {
		if state.player1Score >= maxScore {
			result.player1Wins += universes
		} else if state.player2Score >= maxScore {
			result.player2Wins += universes
		}
	}

	return result
}

func getNextPosition(startingPos, distance int) int {
	result := startingPos + distance
	for result > boardSize {
		result -= boardSize
	}
	return result
}

func rollQuantumDie(sides, rolls int) map[int]uint {
	result := make(map[int]uint)

	// The first roll makes each side spawn a single universe
	for i := 1; i <= sides; i++ {
		result[i] = 1
	}

	for roll := 1; roll < rolls; roll++ {
		nextResult := make(map[int]uint)
		for prevValue, prevUniverses := range result {
			for side := 1; side <= sides; side++ {
				nextResult[prevValue+side] += prevUniverses
			}
		}
		result = nextResult
	}

	return result

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
