package d03

func calculateGammaAndEpsilon(numbers []int) (int, int) {
	bitLength := 16
	setBits := make([]int, bitLength)

	// Create a map of bit index to number of elements that have the corresponding bit set
	for _, value := range numbers {
		for i := 0; i < bitLength; i++ {
			mask := 1 << i
			bitIsSet := value&mask != 0
			if bitIsSet {
				setBits[i]++
			}
		}
	}

	var gamma, epsilon int

	for i := 0; i < bitLength; i++ {
		countWithBitSet := setBits[i]

		if countWithBitSet == 0 {
			continue
		}

		if countWithBitSet >= len(numbers)/2 {
			// when >= 50% have bit set, set corresponding bit in gamma
			gamma = gamma | 1<<i
		} else {
			// when < 50% have bit set, set corresponding bit in epsilon
			epsilon = epsilon | 1<<i
		}
	}

	return gamma, epsilon
}
