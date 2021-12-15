package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type pair [2]byte

type game struct {
	initialPairs       []pair
	pairInsertionRules []pairInsertionRule
	startTime          time.Time
}

type pairInsertionRule struct {
	pair   pair
	insert byte
}

func main() {

	game := parseInput(os.Stdin)

	m := run(game, 10)
	mostCommonChar, leastCommonChar := findMostAndLeastCommon(m)
	fmt.Println(m[mostCommonChar] - m[leastCommonChar])

	m = run(game, 40)
	mostCommonChar, leastCommonChar = findMostAndLeastCommon(m)
	fmt.Println(m[mostCommonChar] - m[leastCommonChar])

}

func run(g game, stepCount int) map[rune]int {

	g.startTime = time.Now()

	cache := make(map[pair][]pair)
	cacheDepth := 0
	primeCache(&g, &cache, cacheDepth)

	characterCounts := make(map[rune]int)
	lengths := make(map[int]int)

	estimatedLength := len(g.initialPairs) + 1
	for i := 0; i < stepCount; i++ {
		estimatedLength *= 2
	}
	fmt.Fprintf(os.Stderr, "Estimated final length after %d steps: %d\n", stepCount, estimatedLength)

	for i := range g.initialPairs {
		p := &g.initialPairs[i]

		runRules(p, &g, stepCount, &characterCounts, &lengths, &cache, cacheDepth)

		isRightMost := i == len(g.initialPairs)-1
		if isRightMost {
			right := rune((*p)[1])
			characterCounts[right]++
		}

		fmt.Fprintf(os.Stderr, "Pair %s complete\n", p.String())
	}

	fmt.Println(lengths)

	return characterCounts
}

func primeCache(g *game, cache *map[pair][]pair, depth int) {

	// first, derive a set of unique letters used
	uniqueLetters := make(map[rune]bool)
	for _, p := range g.initialPairs {
		uniqueLetters[rune(p[0])] = true
		uniqueLetters[rune(p[1])] = true
	}
	for _, rule := range g.pairInsertionRules {
		uniqueLetters[rune(rule.insert)] = true
		uniqueLetters[rune(rule.pair[0])] = true
		uniqueLetters[rune(rule.pair[1])] = true
	}

	// now come up with the full set of combinations
	for left := range uniqueLetters {
		for right := range uniqueLetters {
			p := pair{byte(left), byte(right)}
			(*cache)[p] = solvePairToDepth(&p, g, depth)
		}
	}

	fmt.Fprintf(os.Stderr, "Cached primed with %d pairs to depth %d\n", len(*cache), depth)

}

func countChars(m map[rune]int) int {
	var result int
	for _, ct := range m {
		result += ct
	}
	return result
}

func solvePairToDepth(p *pair, g *game, depth int) []pair {
	result := []pair{*p}

	for i := 0; i < depth; i++ {
		nextLayer := make([]pair, 0, len(result)*2)

		for _, p := range result {

			anyRuleMatched := false

			for _, rule := range g.pairInsertionRules {

				if rule.pair != p {
					continue
				}

				anyRuleMatched = true

				left := pair{p[0], rule.insert}
				right := pair{rule.insert, p[1]}

				nextLayer = append(
					nextLayer,
					left,
					right,
				)

				break

			}

			if !anyRuleMatched {
				nextLayer = append(nextLayer, p)
			}

		}

		result = nextLayer
	}

	return result
}

func runRules(p *pair, g *game, stepCount int, characterCounts *map[rune]int, lengths *map[int]int, cache *map[pair][]pair, cacheDepth int) {

	if rand.Float64() < 0.0000001 {
		length := countChars(*characterCounts)
		duration := time.Now().Sub(g.startTime)

		timePerBillion := (duration.Seconds() / float64(length)) * 1000000000

		fmt.Fprintf(os.Stderr, "Length: %d (time: %fs / billion)\n", length, timePerBillion)
	}

	canUseCache := cacheDepth > 0 && stepCount-cacheDepth > 0

	if canUseCache {

		if solution, found := (*cache)[*p]; found {
			// we have a cached solution for <p> up to a certain number of levels
			for i := range solution {
				runRules(&solution[i], g, stepCount-cacheDepth, characterCounts, lengths, cache, cacheDepth)
			}
			return
		}
	}

	if stepCount <= 0 {
		left := rune((*p)[0])
		(*characterCounts)[left]++
		return
	}

	(*lengths)[stepCount]++

	for i := range g.pairInsertionRules {
		rule := &g.pairInsertionRules[i]

		if rule.pair != *p {
			continue
		}

		// this rule now produces two new pairs
		left := pair{(*p)[0], rule.insert}
		right := pair{rule.insert, (*p)[1]}

		runRules(&left, g, stepCount-1, characterCounts, lengths, cache, cacheDepth)
		runRules(&right, g, stepCount-1, characterCounts, lengths, cache, cacheDepth)

		// We assume that only 1 rule will match per pair
		// that is, there are no duplicate rules
		return
	}

	runRules(p, g, stepCount-1, characterCounts, lengths, cache, cacheDepth)
}

func (p *pair) String() string {
	return string(
		[]rune{
			rune((*p)[0]),
			rune((*p)[1]),
		},
	)
}

func findMostAndLeastCommon(characterCounts map[rune]int) (rune, rune) {
	var mostCommon, leastCommon rune

	for r, count := range characterCounts {
		if mostCommon == rune(0) || count > characterCounts[mostCommon] {
			mostCommon = r
		}
		if leastCommon == rune(0) || count < characterCounts[leastCommon] {
			leastCommon = r
		}
	}

	return mostCommon, leastCommon
}

func parseInput(r io.Reader) game {
	b := bufio.NewScanner(r)
	result := game{}

	for b.Scan() {
		l := strings.TrimSpace(b.Text())

		if len(l) == 0 {
			continue
		}

		if len(result.initialPairs) == 0 {
			result.initialPairs = parsePairs(l)
			continue
		}

		parts := strings.Split(l, " -> ")
		if len(parts) != 2 {
			continue
		}
		result.pairInsertionRules = append(
			result.pairInsertionRules,
			pairInsertionRule{parsePairs(parts[0])[0], parts[1][0]},
		)
	}
	return result
}

func parsePairs(template string) []pair {
	templateLen := len(template)

	if templateLen < 2 {
		return []pair{}
	}

	result := make([]pair, 0, templateLen-1)

	for i := range template {
		if i < templateLen-1 {
			result = append(result, pair{template[i], template[i+1]})
		}
	}

	return result
}
