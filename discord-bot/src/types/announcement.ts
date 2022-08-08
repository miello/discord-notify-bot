export interface IAnnouncement {
  title: string
  href: string
  publishDate: string
}

export interface IOverviewAnnouncement extends IAnnouncement {
  courseTitle: string
}
