package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x, y, z float64
}

type scanner struct {
	name    string
	beacons []point
}

type scannerRelation struct {
	a, b scanner
	// beacons in <a> that correspond to beacons in <b>
	aBeacons []point

	// becons in <b> that correspond to beacons in <a>
	bBeacons []point
}

const MinBeaconsInCommon = 12

func main() {
	l := log.New(os.Stderr, "", log.Default().Flags())
	result := Puzzle01(os.Stdin, l)
	fmt.Println(result)
}

func Puzzle01(r io.Reader, l *log.Logger) string {
	const minBeaconMatchesRequired = 12

	scanners := parseInput(os.Stdin)

	detectedCount := countBeaconsDetected(scanners)

	uniqueCount := countUniqueBeaconsDetected(scanners)

	l.Printf("%d beacons detected\n", detectedCount)
	l.Printf("%d *unique* beacons detected\n", uniqueCount)

	return strconv.Itoa(uniqueCount)
}

func countBeaconsDetected(scanners []scanner) int {
	count := 0
	for _, scanner := range scanners {
		count += len(scanner.beacons)
	}
	return count
}

func countUniqueBeaconsDetected(scanners []scanner) int {
	// To count the total number of beacons:
	// - go through each scanner
	// - for each beacon on the scanner:
	//   - if not visited, add 1
	//   - find all
	// - mark each

	visited := make(map[string]bool)
	count := 0

	for i, scanner := range scanners {
		for _, beacon := range scanner.beacons {

			key := fmt.Sprintf("%s.%v", scanner.name, beacon)
			if visited[key] {
				fmt.Printf("skip: %s\n", key)
				continue
			}
			visited[key] = true

			// find on all other scanners
			for j := range scanners {
				if i == j {
					continue
				}
				rel := compareScanners(scanner, scanners[j])

				if len(rel.aBeacons) < MinBeaconsInCommon {
					continue
				}

				for _, bBeacon := range rel.bBeacons {
					key := fmt.Sprintf("%s.%v", scanners[j].name, bBeacon)
					fmt.Printf("visited: %s\n", key)
					visited[key] = true
				}
			}

			count++
		}
	}

	return count
}

// compareScanners takes two scanners and returns a structure describing the
// relation between each
func compareScanners(a, b scanner) scannerRelation {

	aDistances := buildBeaconDistanceMap(a.beacons)
	bDistances := buildBeaconDistanceMap(b.beacons)

	relation := scannerRelation{
		a: a,
		b: b,
	}

	for aBeacon, distancesInA := range aDistances {

		var bestMatch *point
		var bestScore int

		// if a + b overlap on <aBeacon>, then there must be a single <bBeacon>
		// that corresponds to it. This will be the beacon in <b> that shares the
		// *most* sibling distances in common with *aBeacon*

		for bBeacon, distancesInB := range bDistances {

			score := similarity(distancesInA, distancesInB)

			if score > bestScore {
				bestMatch = &bBeacon
				bestScore = score
			}
		}

		if bestMatch != nil {
			// we believe we have found a match for this beacon
			relation.aBeacons = append(relation.aBeacons, aBeacon)
			relation.bBeacons = append(relation.bBeacons, *bestMatch)
		}
	}

	return relation
}

func buildBeaconDistanceMap(beacons []point) map[point][]float64 {
	result := make(map[point][]float64)

	for _, beacon := range beacons {
		result[beacon] = calcDistances(beacon, beacons)
	}

	return result
}

func calcDistances(beacon point, otherBeacons []point) []float64 {
	var result []float64
	for _, b := range otherBeacons {
		if b == beacon {
			continue
		}
		result = append(result, distance(beacon, b))
	}
	sort.Float64s(result)
	return result
}

func distance(a, b point) float64 {
	return math.Sqrt(
		math.Pow(a.x-b.x, 2) +
			math.Pow(a.y-b.y, 2) +
			math.Pow(a.z-b.z, 2),
	)
}

// similarity returns the number of values <a> and <b> have in common
func similarity(a, b []float64) int {
	var found int
	for _, aValue := range a {
		bIndex := sort.SearchFloat64s(b, aValue)
		if bIndex < len(b) {
			// value is found
			found++
		}
	}
	return found
}

func parseInput(r io.Reader) []scanner {
	s := bufio.NewScanner(r)

	var scanners []scanner

	for s.Scan() {

		l := strings.TrimSpace(s.Text())
		if len(l) == 0 {
			continue
		}

		if strings.Index(l, "---") == 0 {
			scanner := scanner{
				name: strings.TrimSpace(strings.ReplaceAll(l, "---", "")),
			}
			scanners = append(scanners, scanner)
			continue
		}

		tokens := strings.Split(l, ",")
		if len(tokens) != 3 {
			continue
		}

		var nums []float64
		for _, t := range tokens {
			value, err := strconv.ParseFloat(t, 64)
			if err != nil {
				panic(err)
			}
			nums = append(nums, float64(value))
		}

		scanner := &scanners[len(scanners)-1]
		scanner.beacons = append(scanner.beacons, point{nums[0], nums[1], nums[2]})
	}

	return scanners
}
