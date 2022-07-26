package types

type ShortAssignment struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	DueDate string `json:"dueDate"`
}

type AssignmentView struct {
	Assignments []ShortAssignment `json:"assignments"`
	Metadata    PaginateMetadata  `json:"meta"`
}
