package types

type ShortAnnouncement struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Date  string `json:"publishDate"`
}

type AnnouncementView struct {
	Announcements []ShortAnnouncement `json:"announcements"`
	Metadata      PaginateMetadata    `json:"metadata"`
}
