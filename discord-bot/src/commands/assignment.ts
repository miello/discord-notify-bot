import { CacheType, CommandInteraction } from 'discord.js'
import { ICommand } from '../types/command'
import { extractInteractiveInfo } from '../utils/misc'
// import { isBefore } from 'date-fns'
import { nanoid } from 'nanoid'
import { generateNewAssignment } from '../utils/course'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const [courseId, title] = extractInteractiveInfo(interaction)
  const newId = nanoid()
  const [message, row] = await generateNewAssignment(courseId, title, newId)

  const collector = interaction.channel?.createMessageComponentCollector({
    time: 60000,
  })

  collector?.on('collect', async (msg) => {
    const splitMsg = msg.customId.split('#')
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
  commandName: 'assignment',
  description: 'Get assignment from course',
  execute,
  addCourseChoices: true,
} as ICommand
