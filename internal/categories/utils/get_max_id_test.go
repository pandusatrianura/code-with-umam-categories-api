package utils

import (
	"testing"

	"github.com/pandusatrianura/code-with-umam-categories-api/internal/categories/entity"
)

func TestGetMaxID(t *testing.T) {
	tests := []struct {
		name     string
		records  []entity.Category
		expected int64
	}{
		{
			name:     "empty slice",
			records:  []entity.Category{},
			expected: 0,
		},
		{
			name: "single element",
			records: []entity.Category{
				{ID: 5, Name: "Category A", Description: "Desc A"},
			},
			expected: 5,
		},
		{
			name: "multiple elements with increasing IDs",
			records: []entity.Category{
				{ID: 1, Name: "Category A", Description: "Desc A"},
				{ID: 2, Name: "Category B", Description: "Desc B"},
				{ID: 3, Name: "Category C", Description: "Desc C"},
			},
			expected: 3,
		},
		{
			name: "multiple elements with random IDs",
			records: []entity.Category{
				{ID: 10, Name: "Category A", Description: "Desc A"},
				{ID: 3, Name: "Category B", Description: "Desc B"},
				{ID: 15, Name: "Category C", Description: "Desc C"},
				{ID: 7, Name: "Category D", Description: "Desc D"},
			},
			expected: 15,
		},
		{
			name: "multiple elements with duplicate max IDs",
			records: []entity.Category{
				{ID: 12, Name: "Category A", Description: "Desc A"},
				{ID: 7, Name: "Category B", Description: "Desc B"},
				{ID: 12, Name: "Category C", Description: "Desc C"},
			},
			expected: 12,
		},
		{
			name: "negative and positive IDs",
			records: []entity.Category{
				{ID: -10, Name: "Category A", Description: "Desc A"},
				{ID: 0, Name: "Category B", Description: "Desc B"},
				{ID: 5, Name: "Category C", Description: "Desc C"},
				{ID: -3, Name: "Category D", Description: "Desc D"},
			},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMaxID(tt.records)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
