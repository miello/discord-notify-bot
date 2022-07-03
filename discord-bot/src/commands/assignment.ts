import { SlashCommandBuilder, hyperlink } from '@discordjs/builders'
import { CacheType, CommandInteraction, MessageEmbed } from 'discord.js'
import { apiClient } from '../config/axios'
import { IAssignment } from '../types/assignment'
import { ICommand } from '../types/command'
import { extractInteractiveInfo } from '../utils/misc'
import { isBefore } from 'date-fns'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const [courseId, title] = extractInteractiveInfo(interaction)

  const resp = await apiClient.get(`/${courseId}/assignments`)
  const assignments: Array<IAssignment> = resp.data

  const message = new MessageEmbed()
  message.setTitle(`${title} Assignment`)
  message.setThumbnail(
    'https://images-ext-2.discordapp.net/external/4Q85mjDG7508BRnWbBibIMLsL1QYffvT7aq5b4HDaxM/https/www.mycourseville.com/sites/all/modules/courseville/files/logo/cv-logo.png'
  )
  message.setColor('YELLOW')

  assignments.forEach((val) => {
    const { dueDate, title, href } = val
    const dueDateTime = new Date(dueDate)

    if (isBefore(dueDateTime, new Date())) return

    let dueDateString = dueDateTime.toString().split(' ').slice(1, 5).join(' ')
    dueDateString = `Due on ${dueDateString}`

    message.addField(title, hyperlink(dueDateString, href))
  })

  await interaction.reply({ embeds: [message] })
}

export default {
  name: 'assignment',
  data: new SlashCommandBuilder()
    .setName('assignment')
    .setDescription('Get assignment from course'),
  execute,
  addCourseChoices: true,
} as ICommand
