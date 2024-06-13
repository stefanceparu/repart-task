package service

type Calculator interface {
	CalculatePacks(input []int, orderQuantity int) map[int]int
}
