package d19

import (
	"strings"
	"testing"
)

var TEST_INPUT = strings.TrimSpace(`
	--- scanner 0 ---
	404,-588,-901
	528,-643,409
	-838,591,734
	390,-675,-793
	-537,-823,-458
	-485,-357,347
	-345,-311,381
	-661,-816,-575
	-876,649,763
	-618,-824,-621
	553,345,-567
	474,580,667
	-447,-329,318
	-584,868,-557
	544,-627,-890
	564,392,-477
	455,729,728
	-892,524,684
	-689,845,-530
	423,-701,434
	7,-33,-71
	630,319,-379
	443,580,662
	-789,900,-551
	459,-707,401

	--- scanner 1 ---
	686,422,578
	605,423,415
	515,917,-361
	-336,658,858
	95,138,22
	-476,619,847
	-340,-569,-846
	567,-361,727
	-460,603,-452
	669,-402,600
	729,430,532
	-500,-761,534
	-322,571,750
	-466,-666,-811
	-429,-592,574
	-355,545,-477
	703,-491,-529
	-328,-685,520
	413,935,-424
	-391,539,-444
	586,-435,557
	-364,-763,-893
	807,-499,-711
	755,-354,-619
	553,889,-390

	--- scanner 2 ---
	649,640,665
	682,-795,504
	-784,533,-524
	-644,584,-595
	-588,-843,648
	-30,6,44
	-674,560,763
	500,723,-460
	609,671,-379
	-555,-800,653
	-675,-892,-343
	697,-426,-610
	578,704,681
	493,664,-388
	-671,-858,530
	-667,343,800
	571,-461,-707
	-138,-166,112
	-889,563,-600
	646,-828,498
	640,759,510
	-630,509,768
	-681,-892,-333
	673,-379,-804
	-742,-814,-386
	577,-820,562

	--- scanner 3 ---
	-589,542,597
	605,-692,669
	-500,565,-823
	-660,373,557
	-458,-679,-417
	-488,449,543
	-626,468,-788
	338,-750,-386
	528,-832,-391
	562,-778,733
	-938,-730,414
	543,643,-506
	-524,371,-870
	407,773,750
	-104,29,83
	378,-903,-323
	-778,-728,485
	426,699,580
	-438,-605,-362
	-469,-447,-387
	509,732,623
	647,635,-688
	-868,-804,481
	614,-800,639
	595,780,-596

	--- scanner 4 ---
	727,592,562
	-293,-554,779
	441,611,-461
	-714,465,-776
	-743,427,-804
	-660,-479,-426
	832,-632,460
	927,-485,-438
	408,393,-506
	466,436,-512
	110,16,151
	-258,-428,682
	-393,719,612
	-211,-452,876
	808,-476,-593
	-575,615,604
	-485,667,467
	-680,325,-822
	-627,-443,-432
	872,-547,-609
	833,512,582
	807,604,487
	839,-516,451
	891,-625,532
	-652,-548,-490
	30,-46,-14
	`)

func TestOrient(t *testing.T) {
	p := point{1, 2, 3}
	tests := map[orientation]point{
		XyzOrientation: point{1, 2, 3},
		XzyOrientation: point{1, 3, 2},
		YxzOrientation: point{2, 1, 3},
		YzxOrientation: point{2, 3, 1},
		ZxyOrientation: point{3, 1, 2},
		ZyxOrientation: point{3, 2, 1},
	}

	for o, expected := range tests {
		actual := p.orient(o)
		if expected != actual {
			t.Errorf("orient(%s) failed: expected %v, got %v", o, expected, actual)
		}
	}
}

func TestRotateX(t *testing.T) {

	p := point{4, 3, 2}

	expected := []point{
		point{4, 3, 2},
		point{4, -2, 3},
		point{4, -3, -2},
		point{4, 2, -3},
	}

	for i := 0; i < len(expected); i++ {
		degrees := 90 * i
		actual := p.rotateX(float64(degrees))
		if actual != expected[i] {
			t.Errorf("%d: expected %v, got %v", degrees, expected[i], actual)
		}
	}

}

func TestRotateY(t *testing.T) {

	p := point{4, 3, 2}

	expected := []point{
		point{4, 3, 2},
		point{2, 3, -4},
		point{-4, 3, -2},
		point{-2, 3, 4},
	}

	for i := 0; i < len(expected); i++ {
		degrees := 90 * i
		actual := p.rotateY(float64(degrees))
		if actual != expected[i] {
			t.Errorf("%d: expected %v, got %v", degrees, expected[i], actual)
		}
	}

}

func TestRotateZ(t *testing.T) {

	p := point{4, 3, 2}

	expected := []point{
		point{4, 3, 2},
		point{-3, 4, 2},
		point{-4, -3, 2},
		point{3, -4, 2},
	}

	for i := 0; i < len(expected); i++ {
		degrees := 90 * i
		actual := p.rotateZ(float64(degrees))
		if actual != expected[i] {
			t.Errorf("%d: expected %v, got %v", degrees, expected[i], actual)
		}
	}

}

func TestUndoOrient(t *testing.T) {
	p := point{1, 2, 3}

	orientations := []orientation{
		XyzOrientation,
		XzyOrientation,
		YxzOrientation,
		YzxOrientation,
		ZxyOrientation,
		ZyxOrientation,
	}

	for _, o := range orientations {
		oriented := p.orient(o)
		undone := oriented.undoOrient(o)
		if undone != p {
			t.Errorf("undoOrient(%v) failed. Expected %v, got %v", o, p, undone)
		}
	}

}

func TestSolveScanners0And1(t *testing.T) {
	scanners := parseInput(strings.NewReader(TEST_INPUT))

	solution := solve(scanners[0:2])

	expectedScanners := 2
	if len(solution.scanners) != expectedScanners {
		for _, s := range solution.scanners {
			t.Log(s.name)
		}
		t.Fatalf("Solution should include %d scanners, but has %d", expectedScanners, len(solution.scanners))
	}

	scanner0 := solution.scanners[0]
	if scanner0.name != "scanner 0" {
		t.Fatalf("Scanners in wrong order")
	}

	scanner1 := solution.scanners[1]

	if scanner1.name != "scanner 1" {
		t.Fatalf("Scanners in wrong order")
	}

	expectedLocation := point{68, -1246, -43}
	if scanner1.location != expectedLocation {
		t.Errorf("Scanner 1 should be at %v, but was at %v", expectedLocation, scanner1.location)
	}

}

func TestSolveScanners0And1And4(t *testing.T) {
	allScanners := parseInput(strings.NewReader(TEST_INPUT))

	scanners := []scanner{
		allScanners[0],
		allScanners[1],
		allScanners[4],
	}

	solution := solve(scanners)

	expectedLocation := point{-20, -1133, 1061}
	if solution.scanners[2].location != expectedLocation {
		t.Fatalf("Expected %s to be at %v, but was at %v", solution.scanners[2].name, expectedLocation, solution.scanners[2].location)
	}

	expectedScanner1And4Overlaps := []point{
		point{459, -707, 401},
		point{-739, -1745, 668},
		point{-485, -357, 347},
		point{432, -2009, 850},
		point{528, -643, 409},
		point{423, -701, 434},
		point{-345, -311, 381},
		point{408, -1815, 803},
		point{534, -1912, 768},
		point{-687, -1600, 576},
		point{-447, -329, 318},
		point{-635, -1737, 486},
	}

	scanner1BeaconsInGlobalSpace := translateBeacons(solution.scanners[1].beacons, solution.scanners[1].location)
	scanner4BeaconsInGlobalSpace := translateBeacons(solution.scanners[2].beacons, solution.scanners[2].location)

	actualOverlaps := intersection(scanner1BeaconsInGlobalSpace, scanner4BeaconsInGlobalSpace)

	if len(actualOverlaps) != len(expectedScanner1And4Overlaps) {
		t.Errorf("Wrong # of overlaps between scanners 1 + 4 (expected %d, got %d)", len(expectedScanner1And4Overlaps), len(actualOverlaps))
	}

	for _, b := range actualOverlaps {

		found := false
		for i := range expectedScanner1And4Overlaps {
			if expectedScanner1And4Overlaps[i] == b {
				found = true
			}
		}
		if !found {
			t.Errorf("Overlap %v was not expected", b)
		}

	}

	for _, b := range expectedScanner1And4Overlaps {

		found := false
		for i := range actualOverlaps {
			if actualOverlaps[i] == b {
				found = true
			}
		}
		if !found {
			t.Errorf("Overlap %v was expected but not found", b)
		}

	}

}

func TestSolveAllScanners(t *testing.T) {
	scanners := parseInput(strings.NewReader(TEST_INPUT))
	solution := solve(scanners)

	expectedNames := []string{
		"scanner 0",
		"scanner 1",
		"scanner 2",
		"scanner 3",
		"scanner 4",
	}

	for i, s := range solution.scanners {
		if s.name != expectedNames[i] {
			t.Errorf("Expected #%d to have name %v, but got %v", i, expectedNames[i], s.name)
		}
	}

	expectedLocations := []point{
		point{0, 0, 0},
		point{68, -1246, -43},
		point{1105, -1205, 1229},
		point{-92, -2380, -20},
		point{-20, -1133, 1061},
	}

	for i, s := range solution.scanners {
		if s.location != expectedLocations[i] {
			t.Fatalf("Expected #%d to have location %v, but got %v", i, expectedLocations[i], s.location)
		}
	}

	expectedBeacons := []point{
		point{-892, 524, 684},
		point{-876, 649, 763},
		point{-838, 591, 734},
		point{-789, 900, -551},
		point{-739, -1745, 668},
		point{-706, -3180, -659},
		point{-697, -3072, -689},
		point{-689, 845, -530},
		point{-687, -1600, 576},
		point{-661, -816, -575},
		point{-654, -3158, -753},
		point{-635, -1737, 486},
		point{-631, -672, 1502},
		point{-624, -1620, 1868},
		point{-620, -3212, 371},
		point{-618, -824, -621},
		point{-612, -1695, 1788},
		point{-601, -1648, -643},
		point{-584, 868, -557},
		point{-537, -823, -458},
		point{-532, -1715, 1894},
		point{-518, -1681, -600},
		point{-499, -1607, -770},
		point{-485, -357, 347},
		point{-470, -3283, 303},
		point{-456, -621, 1527},
		point{-447, -329, 318},
		point{-430, -3130, 366},
		point{-413, -627, 1469},
		point{-345, -311, 381},
		point{-36, -1284, 1171},
		point{-27, -1108, -65},
		point{7, -33, -71},
		point{12, -2351, -103},
		point{26, -1119, 1091},
		point{346, -2985, 342},
		point{366, -3059, 397},
		point{377, -2827, 367},
		point{390, -675, -793},
		point{396, -1931, -563},
		point{404, -588, -901},
		point{408, -1815, 803},
		point{423, -701, 434},
		point{432, -2009, 850},
		point{443, 580, 662},
		point{455, 729, 728},
		point{456, -540, 1869},
		point{459, -707, 401},
		point{465, -695, 1988},
		point{474, 580, 667},
		point{496, -1584, 1900},
		point{497, -1838, -617},
		point{527, -524, 1933},
		point{528, -643, 409},
		point{534, -1912, 768},
		point{544, -627, -890},
		point{553, 345, -567},
		point{564, 392, -477},
		point{568, -2007, -577},
		point{605, -1665, 1952},
		point{612, -1593, 1893},
		point{630, 319, -379},
		point{686, -3108, -505},
		point{776, -3184, -501},
		point{846, -3110, -434},
		point{1135, -1161, 1235},
		point{1243, -1093, 1063},
		point{1660, -552, 429},
		point{1693, -557, 386},
		point{1735, -437, 1738},
		point{1749, -1800, 1813},
		point{1772, -405, 1572},
		point{1776, -675, 371},
		point{1779, -442, 1789},
		point{1780, -1548, 337},
		point{1786, -1538, 337},
		point{1847, -1591, 415},
		point{1889, -1729, 1762},
		point{1994, -1805, 1792},
	}

	for _, b := range solution.beacons {
		isExpected := false
		for _, eb := range expectedBeacons {
			if eb == b {
				isExpected = true
				break
			}
		}
		if !isExpected {
			t.Errorf("UNEXPECTED BEACON: %d,%d,%d", int(b.x), int(b.y), int(b.z))
		}
	}

	for _, eb := range expectedBeacons {
		wasFound := false
		for _, b := range solution.beacons {
			if b == eb {
				wasFound = true
				break
			}
		}
		if !wasFound {
			t.Errorf("MISSING BEACON: %d,%d,%d", int(eb.x), int(eb.y), int(eb.z))
		}
	}

	if len(solution.beacons) != len(expectedBeacons) {
		t.Errorf("Expected solution to have %d beacons, but got %d", len(expectedBeacons), len(solution.beacons))
	}

}
