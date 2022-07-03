import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'

export interface ICommand {
  name: string
  data: SlashCommandBuilder
  execute: (x: CommandInteraction<CacheType>) => Promise<void>
  addCourseChoices: boolean
}
