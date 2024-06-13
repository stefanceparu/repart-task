package firstfit

import "sort"

type Calc struct{}

func (c Calc) CalculatePacks(input []int, orderQuantity int) map[int]int {
	packsNeeded := make(map[int]int)
	// sort packs sizes in descending order
	packSizes := input
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})

	packSizesMap := convertToMap(packSizes)
	sm := packSizes[len(packSizes)-1]

	// using divide and modulo operations to avoid recursion
	for _, size := range packSizes {
		if orderQuantity >= size {
			packsNeeded[size] = orderQuantity / size
			orderQuantity = orderQuantity % size
		}
	}

	if orderQuantity > 0 {
		packsNeeded[sm]++

		// case when we end up with multiple smallest packs, we need to try to combine them into one
		if packsNeeded[sm] > 1 {
			check := packsNeeded[sm] * sm
			if _, ok := packSizesMap[check]; ok {
				delete(packsNeeded, sm) // remove the smallest pack
				packsNeeded[check] += 1 // add the next pack that is double in size
			}
		}
	}

	return packsNeeded
}

func convertToMap(packSizes []int) map[int]int {
	resultMap := make(map[int]int, len(packSizes))
	for _, size := range packSizes {
		resultMap[size] = size
	}

	return resultMap
}
