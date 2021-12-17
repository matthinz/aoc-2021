package main

import "fmt"

const MaxSteps = 1000
const MinInitialYVelocity = -10000
const MaxInitialYVelocity = 100000

func main() {

	// target area: x=138..184, y=-125..-71
	const minX = 138
	const maxX = 184
	const minY = -125
	const maxY = -71

	recordInitialXVelocity := 0
	recordInitialYVelocity := 0
	allTimeRecordY := 0

	hits := 0
	misses := 0

	minInitialXVelocity := -1 * maxX
	maxInitialXVelocity := maxX + 1

	for initialXVelocity := minInitialXVelocity; initialXVelocity < maxInitialXVelocity; initialXVelocity++ {
		for initialYVelocity := MinInitialYVelocity; initialYVelocity < MaxInitialYVelocity; initialYVelocity++ {

			x := 0
			y := 0
			xVelocity := initialXVelocity
			yVelocity := initialYVelocity
			localRecordY := 0

			for step := 0; step < MaxSteps; step++ {
				x, y, xVelocity, yVelocity = doStep(x, y, xVelocity, yVelocity)

				if y > localRecordY {
					localRecordY = y
				}

				inTarget := x >= minX && x <= maxX && y >= minY && y <= maxY

				if inTarget {
					hits++
					if localRecordY > allTimeRecordY {
						allTimeRecordY = localRecordY
						recordInitialXVelocity = initialXVelocity
						recordInitialYVelocity = initialYVelocity
						fmt.Printf("New all-time record height of %d for velocity %d,%d!!!\n", allTimeRecordY, recordInitialXVelocity, recordInitialYVelocity)
					}
					break
				}

				missedTarget := y < minY
				if missedTarget {
					misses++
					break
				}
			}
		}
	}

	fmt.Printf("%d hits\n", hits)
	fmt.Printf("%d misses\n", misses)
	fmt.Printf("highest: %d\n", allTimeRecordY)

}

func doStep(x, y, xVelocity, yVelocity int) (int, int, int, int) {
	// The probe's x position increases by its x velocity.
	newX := x + xVelocity

	// Due to drag, the probe's x velocity changes by 1 toward the value 0; that
	// is, it decreases by 1 if it is greater than 0, increases by 1 if it is less
	// than 0, or does not change if it is already 0.
	if xVelocity > 0 {
		xVelocity--
	} else if xVelocity < 0 {
		xVelocity++
	}

	// The probe's y position increases by its y velocity.
	newY := y + yVelocity

	// Due to gravity, the probe's y velocity decreases by 1.
	yVelocity--

	return newX, newY, xVelocity, yVelocity
}
