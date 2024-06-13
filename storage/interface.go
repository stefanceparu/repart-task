package storage

type Storage interface {
	AddPacks(sizes []int) error
	RemovePack(size int) error
	RemovePacks()
	GetPacks() []int
}
