import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { ICommand } from '../types/command'
import { schedule } from 'node-cron'
import { client } from '../config/clientBot'
import { Guild } from '../models/channel'
import { extractInteractiveInfo } from '../utils/misc'
import { generateNewOverview } from '../utils/course'
import { nanoid } from 'nanoid'

schedule(
  '30 12 * * *',
  async () => {
    console.log('Running daily cron job')
    const subscribe_guild = await Guild.find()
    subscribe_guild.forEach(async ({ guildId, channelId, courseId }) => {
      const newId = nanoid()
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
      const currentMessage = await channel.send({ embeds, components: row })

      const collector = currentMessage.channel?.createMessageComponentCollector(
        {
          time: 60000,
        }
      )

      collector?.on('collect', async (msg) => {
        const splitMsg = msg.customId.split('#')

        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const [_, id, command, type, page] = splitMsg
        const overviewType = type as 'announcements' | 'assignments'

        if (id !== newId) return

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
  },
  {
    timezone: 'Asia/Bangkok',
  }
)

async function execute(interaction: CommandInteraction<CacheType>) {
  const [courseId, title] = extractInteractiveInfo(interaction)

  const currentGuild = await Guild.where({
    guildId: interaction.guildId,
    channelId: interaction.channelId,
  }).findOne()

  if (!currentGuild) {
    const newGuild = new Guild({
      guildId: interaction.guildId,
      channelId: interaction.channelId,
      courseId: [courseId],
    })

    await newGuild.save()
    await interaction.reply({
      content: `This channel have subscribed to ${title} daily notification (every 12:00)`,
    })
    return
  }

  if (currentGuild.courseId.includes(courseId)) {
    await interaction.reply({
      content: `This channel have already subscribed to ${title}`,
    })
    return
  }

  currentGuild.courseId.push(courseId)
  await currentGuild.save()
  await interaction.reply({
    content: `This channel have subscribed to ${title} daily notification (every 12:00)`,
  })
}

export default {
  name: 'subscribe',
  data: new SlashCommandBuilder()
    .setName('subscribe')
    .setDescription('Subscribe for daily notification'),
  execute,
  addCourseChoices: true,
} as ICommand
