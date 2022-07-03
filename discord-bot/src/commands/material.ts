import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { ICommand } from '../types/command'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  await interaction.reply({ content: 'Material' })
}

export default {
  name: 'material',
  data: new SlashCommandBuilder()
    .setName('material')
    .setDescription('Get material from course'),
  execute,
  addCourseChoices: true,
} as ICommand
