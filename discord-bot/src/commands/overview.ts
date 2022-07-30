import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { nanoid } from 'nanoid'
import { client } from '../config/clientBot'
import { Guild } from '../models/channel'
import { ICommand } from '../types/command'
import { generateNewOverview } from '../utils/course'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const newId = nanoid()
  const subscribe_guild = await Guild.find()
  subscribe_guild.forEach(async ({ guildId, channelId, courseId }) => {
    if (!guildId || !channelId || courseId.length === 0) return

    const guild = client.guilds.cache.get(guildId)
    if (!guild) return

    const channel = guild.channels.cache.get(channelId)
    if (!channel) return

    const textChannel = channel.isText()
    if (!textChannel) return

    const [embeds, row] = await generateNewOverview(
      newId,
      1,
      'announcements',
      courseId
    )
    await interaction.reply({ embeds, components: row })

    const collector = interaction.channel?.createMessageComponentCollector({
      time: 60000,
    })

    collector?.on('collect', async (msg) => {
      const splitMsg = msg.customId.split('-')

      const [id, command, type, page] = splitMsg
      const overviewType = type as 'announcements' | 'assignments'

      if (id.split('_')[1] !== newId) return

      let newPage = +page || 1

      if (command === 'Prev') --newPage
      if (command === 'Next') ++newPage
      if (command === 'Change') {
        newPage = 1
      }

      const [newMessage, newRow] = await generateNewOverview(
        newId,
        newPage,
        overviewType,
        courseId
      )

      await msg.update({ embeds: newMessage, components: newRow })
    })
  })
}

export default {
  name: 'overview',
  data: new SlashCommandBuilder()
    .setName('overview')
    .setDescription('Testing overview feature'),
  execute,
} as ICommand
