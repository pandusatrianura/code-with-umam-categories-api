package utils

import "github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"

// GetMaxID returns the maximum ID value from a slice of Category records. If the slice is empty, it returns 0.
func GetMaxID(records []entity.Category) int64 {
	if len(records) == 0 {
		return 0
	}

	maxID := records[0].ID

	for _, record := range records {
		if record.ID > maxID {
			maxID = record.ID
		}
	}

	return maxID
}
