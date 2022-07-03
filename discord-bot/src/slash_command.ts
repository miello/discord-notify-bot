import { SlashCommandBuilder } from '@discordjs/builders'
import { REST } from '@discordjs/rest'
import { Routes } from 'discord-api-types/v9'
import { DISCORD_TOKEN, APPLICATION_ID } from './env'
import { apiClient } from './axios'

const rest = new REST({ version: '9' }).setToken(DISCORD_TOKEN)

interface ICourse {
  key: string
  courseId: string
  courseTitle: string
  courseHref: string
  semester: number
  year: number
}

export const updateSlashCommand = async () => {
  try {
    const resp = await apiClient.get<Array<ICourse>>('/courses')
    const courses = resp.data
    const choices = courses.map((val) => {
      return {
        name: `${val.courseTitle} (${val.semester}/${val.year})`,
        value: val.key,
      }
    })

    const commands = [
      new SlashCommandBuilder()
        .setName('material')
        .setDescription('Get material from course')
        .addStringOption((option) =>
          option
            .setName('course')
            .setDescription('Course to get material info for')
            .setRequired(true)
            .addChoices(...choices)
        ),
      new SlashCommandBuilder()
        .setName('assignment')
        .setDescription('Get assignment from course')
        .addStringOption((option) =>
          option
            .setName('course')
            .setDescription('Course to get assignment info for')
            .setRequired(true)
            .addChoices(...choices)
        ),
      new SlashCommandBuilder()
        .setName('announcement')
        .setDescription('Get announcement from course')
        .addStringOption((option) =>
          option
            .setName('course')
            .setDescription('Course to get announcement info for')
            .setRequired(true)
            .addChoices(...choices)
        ),
    ].map((command) => command.toJSON())

    await rest.put(Routes.applicationCommands(APPLICATION_ID), {
      body: commands,
    })

    console.log('Successfully registered application commands.')
  } catch (e) {
    console.log(e)
  }
}
