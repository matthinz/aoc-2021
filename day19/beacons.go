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

	// vector describing how axes are rotated
	rotationVector point
}

const MinBeaconsInCommon = 12
const MaxScannerVisibility = 1000

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

// returns a new point with each coordinate inverted
func (p *point) inverse() point {
	return point{
		x: -1 * p.x,
		y: -1 * p.y,
		z: -1 * p.z,
	}
}

func (p *point) multiply(vector point) point {
	return point{
		x: p.x * vector.x,
		y: p.y * vector.y,
		z: p.z * vector.z,
	}
}

// returns a new point translated using the given vector
func (p *point) translate(vector point) point {
	return point{
		x: p.x + vector.x,
		y: p.y + vector.y,
		z: p.z + vector.z,
	}
}

func countBeaconsDetected(scanners []scanner) int {
	count := 0
	for _, scanner := range scanners {
		count += len(scanner.beacons)
	}
	return count
}

func countUniqueBeaconsDetected(scanners []scanner) int {
	uniqueBeacons := make(map[point]bool)

	for i, scanner := range scanners {
		for j := i + 1; j < len(scanners); j++ {
			otherScanner := scanners[j]
			rel := compareScanners(scanner, otherScanner)

			if len(rel.aBeacons) < MinBeaconsInCommon {
				continue
			}

			fmt.Printf("%s -> %s (%d)\n", scanner.name, otherScanner.name, len(rel.aBeacons))

		}

	}

	return len(uniqueBeacons)
}

// compareScanners takes two scanners and returns a structure describing the
// relation between each, including which beacons are equivalent
func compareScanners(a, b scanner) scannerRelation {

	rotationsToTry := make([]point, 0)
	for _, xRotation := range []float64{-1, 1} {
		for _, yRotation := range []float64{-1, 1} {
			for _, zRotation := range []float64{-1, 1} {
				rotationsToTry = append(rotationsToTry, point{xRotation, yRotation, zRotation})
			}
		}
	}

	facingsToTry := []func(point) point{
		func(p point) point { return point{p.x, p.y, p.z} },
		func(p point) point { return point{p.x, p.z, p.y} },
		func(p point) point { return point{p.y, p.x, p.z} },
		func(p point) point { return point{p.y, p.z, p.x} },
		func(p point) point { return point{p.z, p.x, p.y} },
		func(p point) point { return point{p.z, p.y, p.x} },
	}

	nullRotation := point{1, 1, 1}

	var bestResult *scannerRelation

	for _, rotation := range rotationsToTry {

		for _, facing := range facingsToTry {

			relation := scannerRelation{
				a:              a,
				b:              b,
				rotationVector: rotation,
			}

			for _, aBeacon := range a.beacons {
				for _, bBeacon := range b.beacons {

					// Here we assume aBeacon == bBeacon and come up with a new coordinate
					// system with this point as its origin. We can then translate all
					// beacons on <a> and <b> into this new coordinate system.

					// deltaA can be added to points from <a> to convert them into our
					// "universal" coordinate space
					deltaA := aBeacon.inverse()

					// deltaB can be added to points from <b> to convert them into our
					// "universal" coordinate space
					deltaB := bBeacon.inverse()

					aBeaconsInUniversalSpace := rotateAndTranslateBeacons(a.beacons, nullRotation, deltaA)
					aBeaconsInBSpace := rotateAndTranslateBeacons(aBeaconsInUniversalSpace, rotation, deltaB.inverse())
					aBeaconsInBSpace = filterBeacons(aBeaconsInBSpace, beaconIsVisible)
					aBeaconsInBSpace = mapBeacons(aBeaconsInBSpace, facing)

					bBeaconsInUniversalSpace := rotateAndTranslateBeacons(b.beacons, nullRotation, deltaB)
					bBeaconsInASpace := rotateAndTranslateBeacons(bBeaconsInUniversalSpace, rotation, deltaA.inverse())
					bBeaconsInASpace = filterBeacons(bBeaconsInASpace, beaconIsVisible)
					bBeaconsInASpace = mapBeacons(bBeaconsInASpace, facing)

					// If all the a beacons translated into b space are present in b.beacons
					// AND
					// all the b beacons translated into a space are present in a.beacons
					// then we have a rotation that worked!

					if allBeaconsFound(aBeaconsInBSpace, b.beacons) &&
						allBeaconsFound(bBeaconsInASpace, a.beacons) {
						// we have found it
						relation.aBeacons = append(relation.aBeacons, aBeacon)
						relation.bBeacons = append(relation.bBeacons, bBeacon)
					}
				}
			}

			if bestResult == nil || len(relation.aBeacons) > len(bestResult.aBeacons) {
				bestResult = &relation
			}
		}
	}

	return *bestResult
}

func mapBeacons(slice []point, f func(point) point) []point {
	result := make([]point, len(slice))
	for i := range slice {
		result[i] = f(slice[i])
	}
	return result
}

func beaconIsVisible(beacon point) bool {
	return math.Abs(beacon.x) <= MaxScannerVisibility && math.Abs(beacon.y) <= MaxScannerVisibility && math.Abs(beacon.z) <= MaxScannerVisibility
}

// returns true if every point in <slice> is found in <in>
func allBeaconsFound(slice []point, in []point) bool {

	for _, p := range slice {
		found := false
		for _, o := range in {
			if p == o {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func filterBeacons(slice []point, f func(point) bool) []point {
	result := make([]point, 0, len(slice))
	for _, beacon := range slice {
		if f(beacon) {
			result = append(result, beacon)
		}
	}
	return result
}

// rotates each point in <slice> by multiplying it by <rotationVector>, then
// adds <translationVector> to it.
func rotateAndTranslateBeacons(slice []point, rotationVector point, translationVector point) []point {
	result := make([]point, 0, len(slice))
	for i := range slice {
		rotated := slice[i].multiply(rotationVector)
		translated := rotated.translate(translationVector)
		result = append(result, translated)
	}
	return result
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

// countDistancesInCommon returns the number of values <a> and <b> have in common
func countDistancesInCommon(a, b []float64) int {
	counts := make(map[float64]int)

	for _, value := range a {
		counts[value]++
	}

	for _, value := range b {
		counts[value]++
	}

	result := 0
	for _, count := range counts {
		result += count / 2
	}

	return result
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
