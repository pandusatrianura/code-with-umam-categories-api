package entity

// Category represents a grouping of certain entities, containing an ID, name, and a description.
type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
