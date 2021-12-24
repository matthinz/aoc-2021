package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
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

	// orientation used to pick which direction is "up" and help align it with the world
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

// minimum number of beacons two scanners must have in common for them to be considered "the same"
const MinBeaconsInCommon = 12

// maximum visible distance for each scanner (a cube this many units on each side in either direction)
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

////////////////////////////////////////////////////////////////////////////////
// point methods

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

// undoes the application of the given orientation to a point
func (p *point) undoOrient(o orientation) point {
	switch o {
	case XyzOrientation:
		return *p
	case XzyOrientation:
		return p.orient(XzyOrientation)
	case YxzOrientation:
		return p.orient(YxzOrientation)
	case YzxOrientation:
		return p.orient(ZxyOrientation)
	case ZxyOrientation:
		return p.orient(YzxOrientation)
	case ZyxOrientation:
		return p.orient(ZyxOrientation)
	default:
		panic(fmt.Sprintf("Unknown orientation: %s", o))
	}
}

////////////////////////////////////////////////////////////////////////////////

// compareScanners takes two scanners and returns a structure describing the
// relation between each, including which beacons are equivalent
func compareScanners(a, b *scanner) scannerRelation {

	bestOverlap := [2][]point{
		{},
		{},
	}

	tryRotationsAndOrientations(a, func(aRotation point, aOrientation orientation) {
		tryRotationsAndOrientations(b, func(bRotation point, bOrientation orientation) {

			overlap := findOverlap(a, b, aOrientation, bOrientation, aRotation, bRotation)

			if len(overlap[0]) > len(bestOverlap[0]) {
				bestOverlap = overlap
			}

		})
	})

	return scannerRelation{
		a:        a,
		b:        b,
		aBeacons: bestOverlap[0],
		bBeacons: bestOverlap[1],
	}

}

// given two scanners, two orientations and two rotations, try to figure out
// what of each scanner's beacons overlap with each other. return slices
// of beacons from a and b (of the same length), containing the beacons that
// overlap.
func findOverlap(a, b *scanner, aOrientation, bOrientation orientation, aRotation, bRotation point) [2][]point {

	var bestA, bestB *point
	var bestOverlap int

	// Here we move through all combinations of beacons in a and b. For each one,
	// we pretend they are equivalent and project all other beacons on their
	// scanner into a shared coordinate space with the unified beacon at its origin.

	var logf func(string, ...interface{})

	silentLogf := func(format string, a ...interface{}) {}

	verboseLogf := func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}

	rotationIsSpecial := aRotation == NullRotation && bRotation == point{-1, 1, -1}
	orientationIsSpecial := aOrientation == XyzOrientation && bOrientation == XyzOrientation
	SPECIAL_A_POINT := point{-618, -824, -621}
	SPECIAL_B_POINT := point{686, 422, 578}

	for i := range a.beacons {
		aBeacon := a.beacons[i]

		if aBeacon == SPECIAL_A_POINT && rotationIsSpecial && orientationIsSpecial {
			logf = verboseLogf
		} else {
			logf = silentLogf
		}
		// deltaA is a vector that can be *added* to (oriented + rotated) points in <a>
		// to convert them into the shared coordinate space
		deltaA := aBeacon.orient(aOrientation)
		deltaA = deltaA.multiply(aRotation)
		deltaA = deltaA.inverse()

		for j := range b.beacons {
			bBeacon := b.beacons[j]

			if aBeacon == SPECIAL_A_POINT && bBeacon == SPECIAL_B_POINT && rotationIsSpecial && orientationIsSpecial {
				logf = verboseLogf
			} else {
				logf = silentLogf
			}

			bBeaconsInASpace := findBBeaconsInASpace(a, b, aBeacon, bBeacon, aOrientation, bOrientation, aRotation, bRotation)

			if len(bBeaconsInASpace) == 0 {
				continue
			}

			logf("bBeaconsInASpace (%d): %v\n", len(bBeaconsInASpace), bBeaconsInASpace)

			// Now we have a list of beacons from <b> that *should* be in <a> assuming
			// they have lined up correctly. If all of our set of beacons are not
			// present in a, this test failed and we can move on
			if !allBeaconsFound(bBeaconsInASpace, a.beacons) {
				logf("Not all beacons translated from b found in a!\n")
				continue
			}

			logf("All found! bestOverlap is currently %d (%v == %v)\n", bestOverlap, bestA, bestB)

			// At this point, we've found two beacons that, if we assume they are
			// equivalent, give us a potential solution. But we'll keep going
			// in case there exists another combination that provide a better
			// solution. I feel this is unlikely though?
			if len(bBeaconsInASpace) > bestOverlap {
				bestA = &aBeacon
				bestB = &bBeacon
				bestOverlap = len(bBeaconsInASpace)
			}
		}
	}

	if rotationIsSpecial && orientationIsSpecial {
		logf = verboseLogf
	} else {
		logf = silentLogf
	}

	if bestA == nil || bestB == nil {
		// no overlap found
		return [2][]point{
			{},
			{},
		}
	}

	result := [2][]point{
		findBBeaconsInASpace(a, b, *bestA, *bestB, aOrientation, bOrientation, aRotation, bRotation),
		findBBeaconsInASpace(b, a, *bestB, *bestA, bOrientation, aOrientation, bRotation, aRotation),
	}

	if len(result[0]) != len(result[1]) {
		// panic(fmt.Sprintf("result slices for %s (%s,%v) and %s (%s,%v) have different lengths: %d vs %d", a.name, aOrientation, aRotation, b.name, bOrientation, bRotation, len(result[0]), len(result[1])))
		return [2][]point{
			{},
			{},
		}
	}

	return result
}

func findBBeaconsInASpace(
	a, b *scanner,
	aBeacon, bBeacon point,
	aOrientation, bOrientation orientation,
	aRotation, bRotation point,
) []point {

	// deltaA/deltaB are vectors that can be *added* to (oriented + rotated)
	// points in their resepective spaces to convert them into the shared coordinate space
	deltaA := aBeacon.orient(aOrientation)
	deltaA = deltaA.multiply(aRotation)
	deltaA = deltaA.inverse()
	deltaB := bBeacon.orient(bOrientation)
	deltaB = bBeacon.multiply(bRotation)
	deltaB = deltaB.inverse()

	// The process of converting points from <a> into <b>'s spaces is as follows:
	//
	// 1. *Orient* the points, by putting their x,y, and z coordinates in the correct order
	// 2. *Rotate* the points, by applying a rotation (multiplying by e.g. <-1,1,1> to flip across the y axis)
	// 3. *Project* the points into the shared coordinate space by adding the inverse of the origin point to each
	// 4. *Translate* the points out of shared space by adding the (oriented, rotated) origin point from the other space
	// 5. *Unrotate* the points by multiplying by the rotation vector of the other space
	// 6. *Unorient* the points by undoing the orientation applied to the other space
	// 7. *Filter* any points that would not be visible from <a>
	// 8. *Bail* if we did not find enough points to meet the minimum threshold (12)

	bBeaconsInASpace := orientBeacons(b.beacons, bOrientation)
	bBeaconsInASpace = rotateBeacons(bBeaconsInASpace, bRotation)
	bBeaconsInASpace = translateBeacons(bBeaconsInASpace, deltaB)
	bBeaconsInASpace = translateBeacons(bBeaconsInASpace, deltaA.inverse())
	bBeaconsInASpace = rotateBeacons(bBeaconsInASpace, aRotation)
	bBeaconsInASpace = unorientBeacons(bBeaconsInASpace, aOrientation)

	// Since we know that a would not be able to see anything > 1000 units
	// away in any direction, we can remove those
	bBeaconsInASpace = filterPoints(bBeaconsInASpace, beaconIsVisible)

	// One of the rules of the game is that this whole "different perspectives
	// on the same thing" thing only works if the two scanners have at least
	// 12 beacons in common.
	if len(bBeaconsInASpace) < MinBeaconsInCommon {
		return []point{}
	}

	return bBeaconsInASpace
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
			f(rotation, orientation)
		}
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

	count := 0

	for i := range scanners {
		scanner := &scanners[i]
		count += len(scanner.beacons)

		for j := i + 1; j < len(scanners); j++ {
			otherScanner := &scanners[j]
			rel := compareScanners(scanner, otherScanner)

			if len(rel.aBeacons) < MinBeaconsInCommon {
				continue
			}

			count -= len(rel.aBeacons)

		}

	}

	return count
}

////////////////////////////////////////////////////////////////////////////////
// point helpers

func filterPoints(slice []point, f func(point) bool) []point {
	result := make([]point, 0, len(slice))
	for _, beacon := range slice {
		if f(beacon) {
			result = append(result, beacon)
		}
	}
	return result
}

func mapPoints(slice []point, f func(point) point) []point {
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

func orientBeacons(slice []point, o orientation) []point {
	return mapPoints(slice, func(p point) point {
		return p.orient(o)
	})
}

func unorientBeacons(slice []point, o orientation) []point {
	return mapPoints(slice, func(p point) point {
		return p.undoOrient(o)
	})
}

func rotateBeacons(slice []point, rotation point) []point {
	return mapPoints(slice, func(p point) point {
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

////////////////////////////////////////////////////////////////////////////////
// parseInput()

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
