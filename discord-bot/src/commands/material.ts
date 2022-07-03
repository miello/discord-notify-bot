import { hyperlink, SlashCommandBuilder } from '@discordjs/builders'
import {
  CacheType,
  CommandInteraction,
  EmbedFieldData,
  MessageEmbed,
} from 'discord.js'
import { apiClient } from '../config/axios'
import { ICommand } from '../types/command'
import { IMaterial } from '../types/material'
import { extractInteractiveInfo } from '../utils/misc'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  const [courseId, title] = extractInteractiveInfo(interaction)
  const resp = await apiClient.get<Array<IMaterial>>(`/${courseId}/materials`)

  const materials = resp.data

  const embed = new MessageEmbed()
  embed.setTitle(`${title} Materials`)
  embed.setThumbnail(
    'https://images-ext-2.discordapp.net/external/4Q85mjDG7508BRnWbBibIMLsL1QYffvT7aq5b4HDaxM/https/www.mycourseville.com/sites/all/modules/courseville/files/logo/cv-logo.png'
  )
  embed.setColor('AQUA')

  const embedList: EmbedFieldData[] = []

  if (materials) {
    materials.forEach((val) => {
      const { folderName, file } = val

      embedList.push({
        name: folderName,
        value: file
          .splice(0, 4)
          .map((val) => {
            const { title, href } = val
            return hyperlink(title, href)
          })
          .join('\n'),
      })
    })
  }

  embed.addFields(...embedList)

  await interaction.reply({ embeds: [embed] })
}

export default {
  name: 'material',
  data: new SlashCommandBuilder()
    .setName('material')
    .setDescription('Get material from course'),
  execute,
  addCourseChoices: true,
} as ICommand
