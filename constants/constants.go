package constants

const (

	// SuccessCode represents the standard code for a successful operation, typically used in API responses.
	SuccessCode = "1000"

	// ErrorCode represents the constant code for general errors in the system.
	ErrorCode = "2000"

	// ErrCategoryNotFound indicates that the specified category could not be found in the data source.
	ErrCategoryNotFound = "kategori tidak ditemukan"

	// ErrInvalidCategoryID indicates that the provided category ID is invalid or cannot be processed.
	ErrInvalidCategoryID = "id kategori tidak valid"

	// ErrInvalidRequest represents an error message for an invalid category request during request parsing or validation.
	ErrInvalidRequest = "request kategori tidak valid"
)
