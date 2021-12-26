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

type scanner struct {
	name    string
	beacons []point
}

type solution struct {
	// the full set of unique beacons, relative to the global origin
	beacons []point
	// each scanner, located in space, each with its beacons relative to that origin
	scanners []solvedScanner
}

type solvedScanner struct {
	scanner
	location point
}

// minimum number of beacons two scanners must have in common for them to be considered "the same"
const MinBeaconsInCommon = 12

// maximum visible distance for each scanner (a cube this many units on each side in either direction)
const MaxScannerVisibility = 1000

// degrees in which we increment our rotations about the x,y, and z axes
const RotationIncrementDegrees = 90

func main() {
	l := log.New(os.Stderr, "", log.Default().Flags())
	result := Puzzle01(os.Stdin, l)
	fmt.Println(result)
}

func Puzzle01(r io.Reader, l *log.Logger) string {
	scanners := parseInput(os.Stdin)

	solution := solve(scanners)

	return strconv.Itoa(len(solution.beacons))
}

////////////////////////////////////////////////////////////////////////////////

// given a set of scanners, returns a solution
func solve(scanners []scanner) solution {
	result := solution{
		scanners: make([]solvedScanner, 0, len(scanners)),
	}

	uniqueBeacons := make(map[point]bool)

	// Track original indices so we can sort result in the same way
	scannerIndices := make(map[string]int)
	for i, s := range scanners {
		scannerIndices[s.name] = i
	}

	scannerIndex := 0

	for scannerIndex < len(scanners) {
		scanner := scanners[scannerIndex]

		if len(result.scanners) == 0 {
			// first scanner becomes the reference for all others -- we locate it at
			// 0,0,0 in our space
			result.scanners = append(result.scanners, solvedScanner{
				scanner:  scanner,
				location: point{0, 0, 0},
			})

			fmt.Printf("Add beacons from %s\n", scanner.name)
			for _, b := range result.beacons {
				uniqueBeacons[b] = true
			}

			scannerIndex++
			continue
		}

		// try to find an overlap between <scanner> and other scanners we've already
		// overlapped
		foundOverlap := false
		for i := range result.scanners {

			success, solved := solveScanner(scanner, result.scanners[i])
			if !success {
				continue
			}

			result.scanners = append(result.scanners, *solved)

			fmt.Printf("Add beacons from %s\n", solved.name)
			for _, b := range solved.beacons {
				// b is relative to <solved>
				// translate it into our global space
				b = b.translate(solved.location)
				uniqueBeacons[b] = true
			}

			foundOverlap = true
			break
		}

		if foundOverlap {
			scannerIndex++
			continue
		}

		if scannerIndex == len(scanners)-1 {
			panic(fmt.Sprintf("Could not find overlap for scanner %s", scanner.name))
		}

		// move <scanner> to the end of the slice -- hopefully we can solve it later
		fmt.Printf("Moving %s to the end of slice to solve later\n", scanner.name)
		for i := scannerIndex; i < len(scanners)-1; i++ {
			scanners[i] = scanners[i+1]
		}
		scanners[len(scanners)-1] = scanner

	}

	for b := range uniqueBeacons {
		result.beacons = append(result.beacons, b)
	}

	sort.Slice(result.beacons, func(i, j int) bool {
		a := result.beacons[i]
		b := result.beacons[j]

		if a.x != b.x {
			return a.x < b.x
		}

		if a.y != b.y {
			return a.y < b.y
		}

		if a.z != b.z {
			return a.z < b.z
		}

		return false

	})

	sort.SliceStable(result.scanners, func(i, j int) bool {
		return scannerIndices[result.scanners[i].name] < scannerIndices[result.scanners[j].name]
	})

	return result
}

func solveScanner(a scanner, b solvedScanner) (bool, *solvedScanner) {

	var solution *solvedScanner

	tryRotationsAndOrientations(func(rotation point, orientation orientation) {

		aBeacons := orientBeacons(a.beacons, orientation)
		aBeacons = rotateBeacons(aBeacons, rotation)

		for _, aBeacon := range aBeacons {
			for _, bBeacon := range b.beacons {

				// Here we assume that aBeacon == bBeacon, then try to disprove that
				// assumption. We translate all of <a>'s beacons into <b>'s space
				// and search for overlap.

				aBeaconsInBSpace := translateBeacons(aBeacons, aBeacon.inverse())
				aBeaconsInBSpace = translateBeacons(aBeaconsInBSpace, bBeacon)

				anyIllegalOnesFound := false
				beaconsInCommon := 0
				for i := range aBeaconsInBSpace {
					// Don't consider beacons that should not be visible to b
					if !beaconIsVisible(aBeaconsInBSpace[i]) {
						continue
					}

					found := false
					for j := range b.beacons {
						if b.beacons[j] == aBeaconsInBSpace[i] {
							found = true
							break
						}
					}

					if !found {
						anyIllegalOnesFound = true
						break
					}

					beaconsInCommon++
				}

				if anyIllegalOnesFound {
					continue
				}

				if beaconsInCommon < MinBeaconsInCommon {
					continue
				}

				aLocation := bBeacon.translate(aBeacon.inverse())
				aLocation = aLocation.translate(b.location)

				// We have a potential solution
				solution = &solvedScanner{
					scanner: scanner{
						name:    a.name,
						beacons: aBeacons,
					},
					location: aLocation,
				}
			}
		}

	})

	if solution == nil {
		return false, nil
	}

	return true, solution
}

////////////////////////////////////////////////////////////////////////////////
// point methods

func (p *point) applyMatrix(matrix [3][3]float64) point {
	return point{
		math.Round(matrix[0][0]*p.x + matrix[0][1]*p.y + matrix[0][2]*p.z),
		math.Round(matrix[1][0]*p.x + matrix[1][1]*p.y + matrix[1][2]*p.z),
		math.Round(matrix[2][0]*p.x + matrix[2][1]*p.y + matrix[2][2]*p.z),
	}
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

// rotates <p> using the given vector (in degrees)
func (p *point) rotate(vector point) point {
	result := p.rotateX(vector.x)
	result = result.rotateY(vector.y)
	result = result.rotateZ(vector.z)
	return result
}

// rotates <p> <degrees> around the x axis
func (p *point) rotateX(degrees float64) point {
	radians := degrees * (math.Pi / 180)
	matrix := [3][3]float64{
		{1, 0, 0},
		{0, math.Cos(radians), -1 * math.Sin(radians)},
		{0, math.Sin(radians), math.Cos(radians)},
	}
	return p.applyMatrix(matrix)
}

// rotates <p> <degrees> around the y axis
func (p *point) rotateY(degrees float64) point {
	radians := degrees * (math.Pi / 180)
	matrix := [3][3]float64{
		{math.Cos(radians), 0, math.Sin(radians)},
		{0, 1, 0},
		{-1 * math.Sin(radians), 0, math.Cos(radians)},
	}
	return p.applyMatrix(matrix)
}

// rotates <p> <degrees> around the z axis
func (p *point) rotateZ(degrees float64) point {
	radians := degrees * (math.Pi / 180)
	matrix := [3][3]float64{
		{math.Cos(radians), -1 * math.Sin(radians), 0},
		{math.Sin(radians), math.Cos(radians), 0},
		{0, 0, 1},
	}
	return p.applyMatrix(matrix)

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

// calls <f> for a series of rotation / orientation options.
func tryRotationsAndOrientations(f func(point, orientation)) {

	rotationsToTry := make([]point, 0, int(math.Pow(360/RotationIncrementDegrees, 3)))

	for x := 0.0; x < 360; x += RotationIncrementDegrees {
		for y := 0.0; y < 360; y += RotationIncrementDegrees {
			for z := 0.0; z < 360; z += RotationIncrementDegrees {
				rotationsToTry = append(rotationsToTry, point{x, y, z})
			}
		}
	}

	orientationsToTry := []orientation{
		XyzOrientation,
		XzyOrientation,
		YxzOrientation,
		YzxOrientation,
		ZxyOrientation,
		ZyxOrientation,
	}

	for _, rotation := range rotationsToTry {
		for _, orientation := range orientationsToTry {
			f(rotation, orientation)
		}
	}
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
		p = p.rotateX(rotation.x)
		p = p.rotateY(rotation.y)
		p = p.rotateZ(rotation.z)
		return p
	})
}

// returns a new slice in which each point in <slice> is translated by <translationVector>
func translateBeacons(slice []point, translationVector point) []point {
	return mapPoints(slice, func(p point) point {
		return p.translate(translationVector)
	})
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
