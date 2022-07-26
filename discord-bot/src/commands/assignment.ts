import { SlashCommandBuilder, hyperlink } from '@discordjs/builders'
import {
  CacheType,
  CommandInteraction,
  MessageActionRow,
  MessageButton,
  MessageEmbed,
} from 'discord.js'
import { apiClient } from '../config/axios'
import { IAssignment } from '../types/assignment'
import { ICommand } from '../types/command'
import { extractInteractiveInfo } from '../utils/misc'
// import { isBefore } from 'date-fns'
import { nanoid } from 'nanoid'
import { IPaginationMetadata } from '../types/common'
import { MessageButtonStyles } from 'discord.js/typings/enums'

const generateNewAssignment = async (
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

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const [courseId, title] = extractInteractiveInfo(interaction)
  const newId = nanoid()
  const [message, row] = await generateNewAssignment(courseId, title, newId)

  const collector = interaction.channel?.createMessageComponentCollector({
    time: 60000,
  })

  collector?.on('collect', async (msg) => {
    const splitMsg = msg.customId.split('-')
    const [command, page, id] = splitMsg

    if (id !== newId) return

    let newPage = +page

    if (command === 'Prev') --newPage
    if (command === 'Next') ++newPage

    const [newMessage, newRow] = await generateNewAssignment(
      courseId,
      title,
      newId,
      newPage
    )

    await msg.update({ embeds: [newMessage], components: [newRow] })
  })

  await interaction.reply({
    embeds: [message],
    components: [row],
  })
}

export default {
  name: 'assignment',
  data: new SlashCommandBuilder()
    .setName('assignment')
    .setDescription('Get assignment from course'),
  execute,
  addCourseChoices: true,
} as ICommand
