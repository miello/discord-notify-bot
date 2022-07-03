import { hyperlink } from '@discordjs/builders'
import { isBefore, sub } from 'date-fns'
import { MessageEmbed } from 'discord.js'
import { apiClient } from '../config/axios'
import { IAnnouncement } from '../types/announcement'
import { IAssignment } from '../types/assignment'
import { ICourse } from '../types/course'

export async function getCourseChoices() {
  const resp = await apiClient.get<Array<ICourse>>('/courses')
  const courses = resp.data

  const choices = courses.map((val) => {
    return {
      name: `${val.courseTitle} (${val.semester}/${val.year})`,
      value: `${val.key} ${val.courseTitle}`,
    }
  })

  return choices
}

export const getOverviewNotification = async () => {
  const assignmentEmbed = new MessageEmbed()
  const announcementEmbed = new MessageEmbed()

  const resp = await apiClient.get<Array<ICourse>>('/courses')
  const courses = resp.data.map((val) => ({
    name: val.courseTitle,
    key: val.key,
  }))

  const promiseAssignment = courses.map(async (course) => {
    const resp = await apiClient.get<Array<IAssignment>>(
      `/${course.key}/assignments`
    )

    if (!resp.data) return

    const filteredAssignments = resp.data.filter((val) => {
      return isBefore(new Date(val.dueDate), sub(new Date(), { days: 14 }))
    })

    filteredAssignments.forEach((assignment) => {
      const dueDateTime = new Date(assignment.dueDate)

      let dueDateString = dueDateTime
        .toString()
        .split(' ')
        .slice(1, 5)
        .join(' ')

      dueDateString = `Due on ${dueDateString}`
      assignmentEmbed.addField(
        course.name,
        `${hyperlink(assignment.title, assignment.href)} (${dueDateString})`
      )
    })
  })

  const promiseAnnouncement = courses.map(async (course) => {
    const resp = await apiClient.get<Array<IAnnouncement>>(
      `/${course.key}/announcements`
    )

    if (!resp.data) return

    const filteredAnnouncement = resp.data.filter((val) => {
      return isBefore(new Date(val.publishDate), sub(new Date(), { days: 14 }))
    })

    filteredAnnouncement.map((announcement) => {
      const publishedDate = new Date(announcement.publishDate)
        .toString()
        .split(' ')
        .slice(1, 5)
        .join(' ')

      const publishedDateString = `Due on ${publishedDate}`

      announcementEmbed.addField(
        course.name,
        `${hyperlink(
          announcement.title,
          announcement.href
        )} (${publishedDateString})`
      )
    })
  })

  await Promise.all(promiseAssignment)
  await Promise.all(promiseAnnouncement)

  return [assignmentEmbed, announcementEmbed]
}
