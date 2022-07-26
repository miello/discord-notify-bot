package utils

func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}
