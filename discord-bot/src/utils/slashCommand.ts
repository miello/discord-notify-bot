import { SlashCommandBuilder } from '@discordjs/builders'
import { REST } from '@discordjs/rest'
import { Routes } from 'discord-api-types/v9'
import { DISCORD_TOKEN, APPLICATION_ID } from '../config/env'
import { ICommand } from '../types/command'
import { getCourseChoices } from './course'

const rest = new REST({ version: '9' }).setToken(DISCORD_TOKEN)

export const updateSlashCommand = async (executes: Array<ICommand>) => {
  try {
    const choices = await getCourseChoices()
    const commands = executes
      .map((val) => {
        const newCommand = new SlashCommandBuilder()
        newCommand.setName(val.name).setDescription(val.description)
        if (val.addCourseChoices) {
          newCommand.addStringOption((option) =>
            option
              .setName('course')
              .setDescription('Course to get material info for')
              .setRequired(true)
              .addChoices(...choices)
          )
        }
        return newCommand
      })
      .map((val) => val.toJSON())

    await rest.put(Routes.applicationCommands(APPLICATION_ID), {
      body: commands,
    })

    console.log('Successfully registered application commands.')
  } catch (e) {
    console.log(e)
  }
}
