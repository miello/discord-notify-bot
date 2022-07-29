package types

type ShortAssignment struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	DueDate string `json:"dueDate"`
}

type OverviewAssignment struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	DueDate     string `json:"dueDate"`
	CourseTitle string `json:"courseTitle"`
}

type AssignmentView struct {
	Assignments []ShortAssignment `json:"assignments"`
	Metadata    PaginateMetadata  `json:"meta"`
}

type OverviewAssignmentView struct {
	Overview []OverviewAssignment `json:"overviews"`
	Metadata PaginateMetadata     `json:"meta"`
}
