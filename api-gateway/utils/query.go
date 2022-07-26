package utils

import "math"

func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}

func GetTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
