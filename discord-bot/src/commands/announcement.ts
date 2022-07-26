import { SlashCommandBuilder, hyperlink } from '@discordjs/builders'
import {
  CacheType,
  CommandInteraction,
  MessageActionRow,
  MessageButton,
  MessageEmbed,
} from 'discord.js'
import { apiClient } from '../config/axios'
import { ICommand } from '../types/command'
import { extractInteractiveInfo } from '../utils/misc'
import { IAnnouncement } from '../types/announcement'
import { MessageButtonStyles } from 'discord.js/typings/enums'
import { IPaginationMetadata } from '../types/common'

const generateNewAnnouncement = async (
  courseId: string,
  title: string,
  page?: number
): Promise<[MessageEmbed, MessageActionRow]> => {
  const _page = page || 1
  const resp = await apiClient.get(
    `/${courseId}/announcements?page=${_page}&limit=5`
  )
  const announcements: Array<IAnnouncement> = resp.data.announcements
  const metadata: IPaginationMetadata = resp.data.meta

  const message = new MessageEmbed()
  message.setTitle(`${title} Announcement`)
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
      .setCustomId(`Prev ${_page}`)
      .setLabel('Prev')
      .setStyle(MessageButtonStyles.PRIMARY)
      .setDisabled(_page === 1),
    new MessageButton()
      .setCustomId(`Next ${_page}`)
      .setLabel('Next')
      .setStyle(MessageButtonStyles.SECONDARY)
      .setDisabled(metadata.totalPages === _page)
  )

  return [message, row]
}

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const [courseId, title] = extractInteractiveInfo(interaction)
  const [message, row] = await generateNewAnnouncement(courseId, title)

  const collector = interaction.channel?.createMessageComponentCollector({
    time: 60000,
  })

  collector?.on('collect', async (msg) => {
    const splitMsg = msg.customId.split(' ')
    const [command, page] = splitMsg
    let newPage = +page

    if (command === 'Prev') --newPage
    if (command === 'Next') ++newPage

    const [newMessage, newRow] = await generateNewAnnouncement(
      courseId,
      title,
      newPage
    )

    await msg.update({ embeds: [newMessage], components: [newRow] })
  })

  await interaction.reply({ embeds: [message], components: [row] })
}

export default {
  name: 'announcement',
  data: new SlashCommandBuilder()
    .setName('announcement')
    .setDescription('Get announcement from course'),
  execute,
  addCourseChoices: true,
} as ICommand
