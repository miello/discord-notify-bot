import { hyperlink } from '@discordjs/builders'
import { MessageActionRow, MessageButton, MessageEmbed } from 'discord.js'
import { MessageButtonStyles } from 'discord.js/typings/enums'
import { apiClient } from '../config/axios'
import { AnnouncementDTO, OverviewAnnouncementDTO } from '../dtos/announcement'
import { AssignmentDTO, OverviewAssignmentDTO } from '../dtos/assignment'
import { IAnnouncement } from '../types/announcement'
import { IAssignment, IOverviewAssignment } from '../types/assignment'
import { IPaginationMetadata } from '../types/common'
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

export const generateNewEmbed = () => {
  const message = new MessageEmbed()
  message.setThumbnail(
    'https://images-ext-2.discordapp.net/external/4Q85mjDG7508BRnWbBibIMLsL1QYffvT7aq5b4HDaxM/https/www.mycourseville.com/sites/all/modules/courseville/files/logo/cv-logo.png'
  )
  message.setColor('YELLOW')
  return message
}

export const generateNewAnnouncement = async (
  courseId: string,
  title: string,
  id: string,
  page?: number
): Promise<[MessageEmbed, MessageActionRow]> => {
  const _page = page || 1
  const resp = await apiClient.get<AnnouncementDTO>(
    `/${courseId}/announcements?page=${_page}&limit=5`
  )

  const announcements = resp.data.announcements
  const metadata = resp.data.meta

  const message = generateNewEmbed()
  message.setTitle(`${title} Announcement (${_page}/${metadata.totalPages})`)
  message.setURL(
    `https://www.mycourseville.com/?q=courseville/course/${courseId}`
  )

  if (announcements) {
    announcements.forEach((val) => {
      const { publishDate, title, href } = val

      const publishDateTime = new Date(publishDate)

      let publishDateString = publishDateTime
        .toString()
        .split(' ')
        .slice(1, 4)
        .join(' ')
      publishDateString = `${publishDateString}`
      message.addField(publishDateString, hyperlink(title, href))
    })
  }

  const row = new MessageActionRow().addComponents(
    new MessageButton()
      .setCustomId(`Prev-${_page}-${id}`)
      .setLabel('Prev')
      .setStyle(MessageButtonStyles.PRIMARY)
      .setDisabled(_page === 1),
    new MessageButton()
      .setCustomId(`Next-${_page}-${id}`)
      .setLabel('Next')
      .setStyle(MessageButtonStyles.SECONDARY)
      .setDisabled(metadata.totalPages === _page)
  )

  return [message, row]
}

export const generateNewAssignment = async (
  courseId: string,
  title: string,
  id: string,
  page?: number
): Promise<[MessageEmbed, MessageActionRow]> => {
  const _page = page || 1
  const resp = await apiClient.get<AssignmentDTO>(
    `/${courseId}/assignments?page=${_page}&limit=5`
  )
  const assignments: Array<IAssignment> = resp.data.assignments
  const metadata: IPaginationMetadata = resp.data.meta

  const message = generateNewEmbed()
  message.setTitle(`${title} Assignment (${_page}/${metadata.totalPages})`)
  message.setURL(
    `https://www.mycourseville.com/?q=courseville/course/${courseId}/assignment`
  )

  assignments.forEach((val) => {
    const { dueDate, title, href } = val
    const dueDateTime = new Date(dueDate)

    // if (isBefore(dueDateTime, new Date())) return

    let dueDateString = dueDateTime.toString().split(' ').slice(1, 5).join(' ')
    dueDateString = `Due on ${dueDateString}`

    message.addField(title, hyperlink(dueDateString, href))
  })

  const row = new MessageActionRow().addComponents(
    new MessageButton()
      .setCustomId(`Prev-${_page}-${id}`)
      .setLabel('Prev')
      .setStyle(MessageButtonStyles.PRIMARY)
      .setDisabled(_page === 1),
    new MessageButton()
      .setCustomId(`Next-${_page}-${id}`)
      .setLabel('Next')
      .setStyle(MessageButtonStyles.SECONDARY)
      .setDisabled(metadata.totalPages === _page)
  )

  return [message, row]
}

export const generateNewOverview = async (
  id: string,
  page: number,
  type: 'assignments' | 'announcements',
  courseId: string[]
): Promise<[MessageEmbed[], MessageActionRow[]]> => {
  const message = generateNewEmbed()

  const resp = await apiClient.get(
    `/${type}/overview?page=${page}&limit=5&id=${courseId}`,
    {
      params: {
        id: courseId,
        page,
        limit: 5,
      },
    }
  )

  let maxPage = 0

  if (type === 'assignments') {
    const data = resp.data as OverviewAssignmentDTO
    const assignments = data.assignments
    const metadata = data.meta

    message.setTitle(
      `Assignment Daily Notification (${metadata.currentPage}/${metadata.totalPages})`
    )

    assignments.forEach((val) => {
      const { dueDate, title, href, courseTitle } = val
      const dueDateTime = new Date(dueDate)

      // if (isBefore(dueDateTime, new Date())) return

      let dueDateString = dueDateTime
        .toString()
        .split(' ')
        .slice(1, 5)
        .join(' ')
      dueDateString = `Due on ${dueDateString}`

      message.addField(
        dueDateString,
        `${courseTitle}: ${hyperlink(title, href)}`
      )
    })

    maxPage = metadata.totalPages
  }

  if (type === 'announcements') {
    const data = resp.data as OverviewAnnouncementDTO
    const announcements = data.announcements
    const metadata = data.meta

    message.setTitle(
      `Announcement Daily Notification (${metadata.currentPage}/${metadata.totalPages})`
    )

    announcements.forEach((val) => {
      const { publishDate, title, href, courseTitle } = val

      const publishDateTime = new Date(publishDate)

      let publishDateString = publishDateTime
        .toString()
        .split(' ')
        .slice(1, 4)
        .join(' ')
      publishDateString = `${publishDateString}`

      message.addField(
        publishDateString,
        `${courseTitle}: ${hyperlink(title, href)}`
      )
    })

    maxPage = metadata.totalPages
  }

  const typeAction = new MessageActionRow().addComponents(
    new MessageButton()
      .setStyle('SUCCESS')
      .setLabel('Assignments')
      .setCustomId(`overview_${id}-Change-assignments`)
      .setDisabled(type === 'assignments'),
    new MessageButton()
      .setStyle('SECONDARY')
      .setLabel('Announcements')
      .setCustomId(`overview_${id}-Change-announcements`)
      .setDisabled(type === 'announcements')
  )

  const pageAction = new MessageActionRow().addComponents(
    new MessageButton()
      .setStyle('SECONDARY')
      .setLabel('Prev')
      .setCustomId(`overview_${id}-Prev-${type}-${page}`)
      .setDisabled(page === 1),
    new MessageButton()
      .setStyle('PRIMARY')
      .setLabel('Next')
      .setCustomId(`overview_${id}-Next-${type}-${page}`)
      .setDisabled(page >= maxPage)
  )

  return [[message], [typeAction, pageAction]]
}

// export const getOverviewNotification = async (
//   id: string,
//   page: number,
//   type: string,
//   courseId: string[]
// ): Promise<[MessageEmbed[], MessageActionRow[]]> => {
//   const title = type[0].toLocaleUpperCase() + type.slice(1)
//   message.setTitle(`${title} daily notification (page ${page})`)

// }
