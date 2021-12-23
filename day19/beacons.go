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

type orientation string

const (
	UnknownOrientation = ""
	XyzOrientation     = "xyz"
	XzyOrientation     = "xzy"
	YxzOrientation     = "yxz"
	YzxOrientation     = "yzx"
	ZxyOrientation     = "zxy"
	ZyxOrientation     = "zyx"
)

var NullRotation = point{1, 1, 1}

type scanner struct {
	name    string
	beacons []point
	// orientation used to pick which direction is "up"
	orientation orientation
	// rotation vector used to align this scanner to the world
	rotation point
}

type scannerRelation struct {
	a, b *scanner
	// beacons in <a> that correspond to beacons in <b>
	aBeacons []point

	// becons in <b> that correspond to beacons in <a>
	bBeacons []point
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

func (p *point) orient(o orientation) point {
	switch o {
	case XyzOrientation:
		return *p
	case XzyOrientation:
		return point{p.x, p.z, p.y}
	case YxzOrientation:
		return point{p.y, p.x, p.z}
	case YzxOrientation:
		return point{p.y, p.z, p.x}
	case ZxyOrientation:
		return point{p.z, p.x, p.y}
	case ZyxOrientation:
		return point{p.z, p.y, p.x}
	default:
		panic(fmt.Sprintf("Invalid orientation: %s", string(o)))
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

	for i := range scanners {
		scanner := &scanners[i]
		for j := i + 1; j < len(scanners); j++ {
			otherScanner := &scanners[j]
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
func compareScanners(a, b *scanner) scannerRelation {

	bestResult := [][]point{
		{},
		{},
	}
	var bestARotation, bestBRotation point
	var bestAOrientation, bestBOrientation orientation

	tryRotationsAndOrientations(a, func(aRotation point, aOrientation orientation) {
		tryRotationsAndOrientations(b, func(bRotation point, bOrientation orientation) {

			correspondingBeacons := findOverlap(a, b, aOrientation, bOrientation, aRotation, bRotation)

			if len(correspondingBeacons[0]) > len(bestResult[0]) {
				bestResult = correspondingBeacons
				bestAOrientation = aOrientation
				bestBOrientation = bOrientation
				bestARotation = aRotation
				bestBRotation = bRotation
			}

		})

	})

	if len(bestResult[0]) < MinBeaconsInCommon {
		// not enough beacons in common == not the same
		return scannerRelation{
			a: a,
			b: b,
		}
	}

	a.orientation = bestAOrientation
	a.rotation = bestARotation

	b.orientation = bestBOrientation
	b.rotation = bestBRotation

	return scannerRelation{
		a:        a,
		b:        b,
		aBeacons: bestResult[0],
		bBeacons: bestResult[1],
	}
}

// given two scanners, two orientations and two rotations, try to figure out
// what of each scanner's beacons overlap with each other
func findOverlap(a, b *scanner, aOrientation, bOrientation orientation, aRotation, bRotation point) [][]point {

	bestResult := [][]point{
		{},
		{},
	}

	for _, aBeacon := range a.beacons {

		aIsSpecial := aOrientation == XyzOrientation && bOrientation == ZxyOrientation &&
			aRotation == point{-1, 1, -1} && bRotation == point{-1, -1, 1}

		aIsSpecial = aIsSpecial && aBeacon == point{-391, 539, -444}

		// isSpecial := aOrientation == XyzOrientation && bOrientation == XyzOrientation &&
		// 	aRotation == point{-1, 1, -1} && bRotation == point{-1, 1, -1}

		aBeacon = aBeacon.orient(aOrientation)

		if aIsSpecial {
			fmt.Printf("oriented aBeacon: %v\n", aBeacon)
		}

		for _, bBeacon := range b.beacons {

			isSpecial := aIsSpecial && (math.Abs(bBeacon.x) == 660 || math.Abs(bBeacon.y) == 660 || math.Abs(bBeacon.z) == 660)

			if isSpecial {
				fmt.Println(bBeacon)
			}

			bBeacon = bBeacon.orient(bOrientation)

			isSpecial = isSpecial && bBeacon == point{-426, -660, -479}

			if isSpecial {
				fmt.Printf("oriented bBeacon: %v\n", bBeacon)
			}

			// Here we assume aBeacon == bBeacon. Then we work to prove that this is false.

			// First, project all beacons in <a> and <b> into the coordinate space where
			// the origin is at <aBeacon> or <bBeacon>.

			// deltaA can be added to points from <a> to convert them into our
			// "universal" coordinate space
			deltaA := aBeacon.inverse()

			// deltaB can be added to points from <b> to convert them into our
			// "universal" coordinate space
			deltaB := bBeacon.inverse()

			if isSpecial {
				fmt.Printf("deltaA: %v\ndeltaB: %v\n", deltaA, deltaB)
			}

			aBeaconsInUniversalSpace := orientBeacons(a.beacons, aOrientation)

			if isSpecial {
				fmt.Printf("aBeaconsInUniversalSpace (oriented): %v\n", aBeaconsInUniversalSpace)
			}

			aBeaconsInUniversalSpace = translateBeacons(aBeaconsInUniversalSpace, deltaA)
			if isSpecial {
				fmt.Printf("aBeaconsInUniversalSpace (translated): %v\n", aBeaconsInUniversalSpace)
			}

			aBeaconsInUniversalSpace = rotateBeacons(aBeaconsInUniversalSpace, aRotation)

			if isSpecial {
				fmt.Printf("aBeaconsInUniversalSpace (rotated): %v\n", aBeaconsInUniversalSpace)
			}

			aBeaconsInBSpace := translateBeacons(aBeaconsInUniversalSpace, deltaB.inverse())
			if isSpecial {
				fmt.Printf("aBeaconsInBSpace (translated): %v\n", aBeaconsInBSpace)
			}

			aBeaconsInBSpace = filterBeacons(aBeaconsInBSpace, beaconIsVisible)

			if isSpecial {
				fmt.Printf("aBeaconsInBSpace (filtered): %v\n", aBeaconsInBSpace)
			}

			bBeaconsInUniversalSpace := orientBeacons(b.beacons, bOrientation)
			if isSpecial {
				fmt.Printf("bBeaconsInUniversalSpace (oriented): %v\n", bBeaconsInUniversalSpace)
			}

			bBeaconsInUniversalSpace = translateBeacons(bBeaconsInUniversalSpace, deltaB)
			if isSpecial {
				fmt.Printf("bBeaconsInUniversalSpace (translated): %v\n", bBeaconsInUniversalSpace)
			}

			bBeaconsInUniversalSpace = rotateBeacons(bBeaconsInUniversalSpace, bRotation)

			if isSpecial {
				fmt.Printf("bBeaconsInUniversalSpace (rotated): %v\n", bBeaconsInUniversalSpace)
			}

			bBeaconsInASpace := translateBeacons(bBeaconsInUniversalSpace, deltaA.inverse())

			if isSpecial {
				fmt.Printf("bBeaconsInASpace (translated): %v\n", bBeaconsInASpace)
			}

			bBeaconsInASpace = filterBeacons(bBeaconsInASpace, beaconIsVisible)

			if isSpecial {
				fmt.Printf("bBeaconsInASpace (filtered): %v\n", bBeaconsInASpace)
			}

			if len(aBeaconsInBSpace) != len(bBeaconsInASpace) {
				if isSpecial {
					fmt.Printf("different # of beacons found in a/b\n")
				}
				continue
			}

			if !allBeaconsFound(aBeaconsInBSpace, b.beacons) {
				if isSpecial {
					fmt.Printf("Not all a beacons found in b space for %v == %v\n", aBeacon, bBeacon)
				}
				continue
			}

			if !allBeaconsFound(bBeaconsInASpace, a.beacons) {
				if isSpecial {
					fmt.Printf("Not all b beacons found in a space for %v == %v\n", aBeacon, bBeacon)
				}
				continue
			}

			// In these rotations / orientation, len(aBeaconsInBSpace) beacons
			// are shared in common between a + b
			if len(bestResult[0]) < len(aBeaconsInBSpace) {
				if isSpecial {
					fmt.Printf("new best result: %v == %v yields %d in common\n", aBeacon, bBeacon, len(aBeaconsInBSpace))
				}
				bestResult = [][]point{
					bBeaconsInASpace,
					aBeaconsInBSpace,
				}
			} else {
				if isSpecial {
					fmt.Printf("NOT a new best result: %v == %v yields %d in common\n", aBeacon, bBeacon, len(aBeaconsInBSpace))
				}

			}
		}
	}

	return bestResult
}

// calls <f> for a series of rotation / orientation options. If <f> returns
// true, this means the try was successful, and the rotation/orientation combo
// is recorded on <scanner> for future use
func tryRotationsAndOrientations(scanner *scanner, f func(point, orientation)) {

	var rotationsToTry []point
	var orientationsToTry []orientation

	if scanner.rotation != NullRotation {
		rotationsToTry = []point{scanner.rotation}
	} else {
		for _, x := range []float64{-1, 1} {
			for _, y := range []float64{-1, 1} {
				for _, z := range []float64{-1, 1} {
					rotationsToTry = append(rotationsToTry, point{x, y, z})
				}
			}
		}

	}

	if scanner.orientation != UnknownOrientation {
		orientationsToTry = []orientation{scanner.orientation}
	} else {
		orientationsToTry = []orientation{
			XyzOrientation,
			XzyOrientation,
			YxzOrientation,
			YzxOrientation,
			ZxyOrientation,
			ZyxOrientation,
		}
	}

	for _, rotation := range rotationsToTry {
		for _, orientation := range orientationsToTry {
			// fmt.Printf("try rotation %v, orientation %s\n", rotation, orientation)
			f(rotation, orientation)
		}
	}
}

func mapBeacons(slice []point, f func(point) point) []point {
	result := make([]point, len(slice))
	for i := range slice {
		result[i] = f(slice[i])
	}
	return result
}

func intersection(a, b []point) []point {
	var result []point

	for _, aPoint := range a {
		for _, bPoint := range b {
			if aPoint == bPoint {
				result = append(result, aPoint)
			}
		}
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

func orientBeacons(slice []point, o orientation) []point {
	return mapBeacons(slice, func(p point) point {
		return p.orient(o)
	})
}

func rotateBeacons(slice []point, rotation point) []point {
	return mapBeacons(slice, func(p point) point {
		return p.multiply(rotation)
	})
}

// returns a new slice in which each point in <slice> is translated by <translationVector>
func translateBeacons(slice []point, translationVector point) []point {
	result := make([]point, 0, len(slice))
	for i := range slice {
		translated := slice[i].translate(translationVector)
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
				name:        strings.TrimSpace(strings.ReplaceAll(l, "---", "")),
				orientation: UnknownOrientation,
				rotation:    NullRotation,
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
