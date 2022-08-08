package types

type CourseView struct {
	ID       string `json:"key"`
	Key      string `json:"courseId"`
	Title    string `json:"courseTitle"`
	Href     string `json:"courseHref"`
	Semester int    `json:"semester"`
	Year     int    `json:"year"`
}
