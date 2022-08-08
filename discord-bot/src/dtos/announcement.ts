import { IAnnouncement, IOverviewAnnouncement } from '../types/announcement'
import { IPaginationMetadata } from '../types/common'

export interface AnnouncementDTO {
  announcements: Array<IAnnouncement>
  meta: IPaginationMetadata
}

export interface OverviewAnnouncementDTO {
  announcements: Array<IOverviewAnnouncement>
  meta: IPaginationMetadata
}
