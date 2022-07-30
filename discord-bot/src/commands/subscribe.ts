import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { ICommand } from '../types/command'
import { schedule } from 'node-cron'
import { client } from '../config/clientBot'
import { getSubscriberList } from '../utils/subscribe'
// import { getOverviewNotification } from '../utils/course'
import { Guild } from '../models/channel'
import { extractInteractiveInfo } from '../utils/misc'

// schedule(
//   '30 12 * * *',
//   () => {
//     console.log('Running daily cron job')
//     getSubscriberList().forEach((val) => {
//       const guild = client.guilds.cache.get(val[0])
//       if (!guild) return

//       val[1].forEach(async (channelId) => {
//         const channel = guild.channels.cache.get(channelId)

//         if (!channel) return

//         const textChannel = channel.isText()
//         if (!textChannel) return

//         const [embeds, row] = await getOverviewNotification()
//         const message = await channel.send({ embeds, components: row })

//         const collector = message.channel?.createMessageComponentCollector({
//           time: 24 * 60 * 60 * 1000,
//         })
//       })
//     })
//   },
//   {
//     timezone: 'Asia/Bangkok',
//   }
// )

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
