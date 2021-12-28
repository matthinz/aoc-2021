package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type pair struct{ left, right byte }

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

	pairCounts := make(map[pair]int)

	// initialize pair counts
	for i := range g.initialPairs {
		pairCounts[g.initialPairs[i]]++
	}

	for stepIndex := 0; stepIndex < stepCount; stepIndex++ {
		pairCounts = tick(&g, pairCounts)
	}

	result := make(map[rune]int)
	for pair, count := range pairCounts {
		result[rune(pair.left)] += count
	}

	// Make sure we count the rightmost element of the last pair
	lastInitialPair := g.initialPairs[len(g.initialPairs)-1]
	result[rune(lastInitialPair.right)]++

	return result
}

func tick(g *game, pairCounts map[pair]int) map[pair]int {
	result := make(map[pair]int, len(pairCounts))

	for p, count := range pairCounts {

		ruleApplied := false

		for _, rule := range g.pairInsertionRules {

			if rule.pair != p {
				continue
			}

			left := pair{p.left, rule.insert}
			result[left] += count

			right := pair{rule.insert, p.right}
			result[right] += count

			ruleApplied = true
			break
		}

		if !ruleApplied {
			// when no rule was applied, the pair survives to the next step
			result[p] += count
		}
	}

	return result

}

func (p *pair) String() string {
	return string(
		[]rune{
			rune((*p).left),
			rune((*p).right),
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
