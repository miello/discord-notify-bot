import { hyperlink } from '@discordjs/builders'
import { MessageActionRow, MessageButton, MessageEmbed } from 'discord.js'
import { MessageButtonStyles } from 'discord.js/typings/enums'
import { apiClient } from '../config/axios'
import { IAnnouncement } from '../types/announcement'
import { IAssignment } from '../types/assignment'
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

export const generateNewAnnouncement = async (
  courseId: string,
  title: string,
  id: string,
  page?: number
): Promise<[MessageEmbed, MessageActionRow]> => {
  const _page = page || 1
  const resp = await apiClient.get(
    `/${courseId}/announcements?page=${_page}&limit=5`
  )

  const announcements: Array<IAnnouncement> = resp.data.announcements
  const metadata: IPaginationMetadata = resp.data.meta

  const message = new MessageEmbed()
  message.setTitle(`${title} Announcement (${_page}/${metadata.totalPages})`)
  message.setThumbnail(
    'https://images-ext-2.discordapp.net/external/4Q85mjDG7508BRnWbBibIMLsL1QYffvT7aq5b4HDaxM/https/www.mycourseville.com/sites/all/modules/courseville/files/logo/cv-logo.png'
  )
  message.setColor('YELLOW')
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
  const resp = await apiClient.get(
    `/${courseId}/assignments?page=${_page}&limit=5`
  )
  const assignments: Array<IAssignment> = resp.data.assignments
  const metadata: IPaginationMetadata = resp.data.meta

  const message = new MessageEmbed()
  message.setTitle(`${title} Assignment (${_page}/${metadata.totalPages})`)
  message.setThumbnail(
    'https://images-ext-2.discordapp.net/external/4Q85mjDG7508BRnWbBibIMLsL1QYffvT7aq5b4HDaxM/https/www.mycourseville.com/sites/all/modules/courseville/files/logo/cv-logo.png'
  )
  message.setColor('YELLOW')
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

export const getOverviewNotification = async () => {
  const assignmentEmbed = new MessageEmbed()
  const announcementEmbed = new MessageEmbed()

  const rowsAction = new MessageActionRow().addComponents(new MessageButton())

  return [assignmentEmbed, announcementEmbed]
}
