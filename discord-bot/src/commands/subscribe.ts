import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { ICommand } from '../types/command'
import { schedule } from 'node-cron'
import { client } from '../config/clientBot'
import { getSubscriberList, subscribe } from '../utils/subscribe'
import { getOverviewNotification } from '../utils/course'

schedule(
  '30 12 * * *',
  () => {
    console.log('Running daily cron job')
    getSubscriberList().forEach((val) => {
      const guild = client.guilds.cache.get(val[0])
      if (!guild) return

      val[1].forEach(async (channelId) => {
        const channel = guild.channels.cache.get(channelId)

        if (!channel) return

        const textChannel = channel.isText()
        if (!textChannel) return

        const embeds = await getOverviewNotification()
        const message = await channel.send({ embeds })

        const collector = message.channel?.createMessageComponentCollector({
          time: 24 * 60 * 60 * 1000,
        })
      })
    })
  },
  {
    timezone: 'Asia/Bangkok',
  }
)

async function execute(interaction: CommandInteraction<CacheType>) {
  await interaction.reply({
    content: 'Not Implemented Yet',
  })
  return
  subscribe(interaction.guildId || '', interaction.channelId || '')

  await interaction.reply({
    content: 'This channel have subscribed to daily notification (every 12:00)',
  })
}

export default {
  name: 'subscribe',
  data: new SlashCommandBuilder()
    .setName('subscribe')
    .setDescription('Subscribe for daily notification'),
  execute,
} as ICommand
