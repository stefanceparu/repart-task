package bestfit

import (
	"sort"
)

type Calc struct {
	bestFit map[int]int
	bestSum int
}

func NewCalc() *Calc {
	return &Calc{
		bestFit: make(map[int]int),
		bestSum: 0,
	}
}

// CalculatePacks function is used to find the best fit for the target value
func (c *Calc) CalculatePacks(packs []int, target int) map[int]int {
	// clear any values stored previously
	c.reset()

	// sort packs in ascending order.
	sort.Ints(packs)

	// start processing combinations
	c.findCombinations(packs, map[int]int{}, 0, 0, target)

	// if there is a reminder, then we'll append it to the smaller pack
	rem := target - c.bestSum
	if rem > 0 {
		c.bestFit[packs[0]]++
	}

	// if possible, combine small packs into larger ones
	result := combinePacks(c.bestFit, packs)

	return result
}

func (c *Calc) findCombinations(packs []int, current map[int]int, currentSum int, start int, target int) {
	// stop condition
	if currentSum > target {
		return
	}

	// check to see if we found a better solution.
	if currentSum >= c.bestSum {
		c.bestSum = currentSum        // save the new best sum
		c.bestFit = make(map[int]int) // reset so that we store only best solution.

		// store new solution.
		for k, v := range current {
			c.bestFit[k] = v
		}
	}

	// loop through each pack to ensure that each pack will be considered at least once.
	for i := start; i < len(packs); i++ {
		// create new copy, so that we don't modify current map
		next := make(map[int]int)
		for k, v := range current {
			next[k] = v
		}

		// increment the current pack size count and explore new combination
		next[packs[i]]++

		c.findCombinations(packs, next, currentSum+packs[i], i, target)
	}
}

func (c *Calc) reset() {
	c.bestFit = make(map[int]int)
	c.bestSum = 0
}

// combinePacks helps accommodate smaller packs into larger ones.
func combinePacks(m map[int]int, packs []int) map[int]int {
	// sort available packs in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(packs)))

	for i := 0; i < len(packs); i++ {
		for j := i + 1; j < len(packs); j++ {
			// if the larger pack is a multiple of the smaller ones.
			// example: pack[i] = 500, pack[j] = 250
			if packs[i]%packs[j] == 0 {
				// calculate the number of larger packs that you can achieve, example:
				// if we have: m[250] = 2
				// pack 500 / pack 250 = 2
				// then pack 500 can be used 1 times instead of 2 x 250
				total := m[packs[j]] / (packs[i] / packs[j])

				// if the total is 0, we don't want to store it in map.
				if total == 0 {
					continue
				}

				// increment larger pack count
				m[packs[i]] += total

				// subtract from the count of smaller packs
				m[packs[j]] -= total * (packs[i] / packs[j])

				// if the remaining count is 0 then there is no need to keep the data in the map.
				if m[packs[j]] == 0 {
					delete(m, packs[j])
				}

				// check to see if this larger pack can also be combined into a larger one
				if m[packs[i]] > 1 {
					combinePacks(m, packs)
				}
			}
		}
	}

	return m
}
