package shared_module

type SortDirection string

const (
	// Ascending sort order
	Asc SortDirection = "ASC"
	// Descending sort order
	Desc SortDirection = "DESC"
)

type DeleteOperation struct {
	Success bool
	Message string
}

type PaginationInput struct {
	Limit  *int
	Offset *int
}

type SortingInput struct {
	Field     string
	Direction SortDirection
}
