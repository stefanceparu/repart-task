package memory

import (
	"fmt"
)

type MemDB struct {
	data []int
}

func NewMemDB() *MemDB {
	return &MemDB{}
}

func (db *MemDB) AddPacks(sizes []int) error {
	// validate & remove duplicates, if exist.
	existing := db.convertToMap()
	var valid []int
	for _, size := range sizes {
		if size <= 0 {
			return fmt.Errorf("pack size must be positive %d", size)
		}

		if _, ok := existing[size]; !ok {
			valid = append(valid, size)
		}
	}

	db.data = append(db.data, valid...)
	return nil
}

func (db *MemDB) RemovePack(size int) error {
	for i := 0; i < len(db.data); i++ {
		if db.data[i] == size {
			db.data = append(db.data[:i], db.data[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("size %d not found", size)
}

func (db *MemDB) GetPacks() []int {
	return db.data
}

func (db *MemDB) RemovePacks() {
	db.data = []int{}
}

func (db *MemDB) convertToMap() map[int]int {
	result := map[int]int{}
	for _, size := range db.data {
		result[size] = size
	}

	return result
}
