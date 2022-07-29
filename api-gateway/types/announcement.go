package types

type ShortAnnouncement struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Date  string `json:"publishDate"`
}

type OverviewAnnouncement struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	Date        string `json:"publishDate"`
	CourseTitle string `json:"courseTitle"`
}

type AnnouncementView struct {
	Announcements []ShortAnnouncement `json:"announcements"`
	Metadata      PaginateMetadata    `json:"meta"`
}

type OverviewAnnouncementView struct {
	Overview []OverviewAnnouncement `json:"overviews"`
	Metadata PaginateMetadata       `json:"meta"`
}
