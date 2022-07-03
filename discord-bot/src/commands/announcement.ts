import { SlashCommandBuilder, hyperlink } from '@discordjs/builders'
import { CacheType, CommandInteraction, MessageEmbed } from 'discord.js'
import { apiClient } from '../config/axios'
import { ICommand } from '../types/command'
import { extractInteractiveInfo } from '../utils/misc'
import { IAnnouncement } from '../types/announcement'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const [courseId, title] = extractInteractiveInfo(interaction)

  const resp = await apiClient.get(`/${courseId}/announcements`)
  const assignments: Array<IAnnouncement> = resp.data

  const message = new MessageEmbed()
  message.setTitle(`${title} Announcement`)
  message.setThumbnail(
    'https://images-ext-2.discordapp.net/external/4Q85mjDG7508BRnWbBibIMLsL1QYffvT7aq5b4HDaxM/https/www.mycourseville.com/sites/all/modules/courseville/files/logo/cv-logo.png'
  )
  message.setColor('YELLOW')
  message.setURL(
    `https://www.mycourseville.com/?q=courseville/course/${courseId}`
  )

  if (assignments) {
    assignments.forEach((val) => {
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

  await interaction.reply({ embeds: [message] })
}

export default {
  name: 'announcement',
  data: new SlashCommandBuilder()
    .setName('announcement')
    .setDescription('Get announcement from course'),
  execute,
  addCourseChoices: true,
} as ICommand
