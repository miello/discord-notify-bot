import { CacheType, CommandInteraction } from 'discord.js'

export interface ICommand {
  name: string
  commandName: string
  description: string
  execute: (x: CommandInteraction<CacheType>) => Promise<void>
  addCourseChoices: boolean
}
