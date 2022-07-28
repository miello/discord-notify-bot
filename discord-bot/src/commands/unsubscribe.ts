import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { Guild } from '../models/channel'
import { ICommand } from '../types/command'
import { extractInteractiveInfo } from '../utils/misc'

async function execute(interaction: CommandInteraction<CacheType>) {
  const [courseId, title] = extractInteractiveInfo(interaction)

  const currentGuild = await Guild.where({
    guildId: interaction.guildId,
    channelId: interaction.channelId,
  }).findOne()

  const courseIdIdx = currentGuild?.courseId.indexOf(courseId[0])

  if (!currentGuild || courseIdIdx === -1) {
    await interaction.reply({
      content: `This channel have not subscribed to ${title} daily notification`,
    })
    return
  }

  currentGuild.courseId.splice(courseIdIdx || 0, 1)
  await interaction.reply({
    content: `This channel have unsubscribed to ${title} daily notification`,
  })
  return
}

export default {
  name: 'unsubscribe',
  data: new SlashCommandBuilder()
    .setName('unsubscribe')
    .setDescription('Unsubscribe for daily notification'),
  execute,
} as ICommand
