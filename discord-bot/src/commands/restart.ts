import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { client } from '../config/clientBot'
import { BOT_OWNER_ID, DISCORD_TOKEN } from '../config/env'
import { commandsList } from '../config/routes'
import { ICommand } from '../types/command'
import { updateSlashCommand } from '../utils/slashCommand'

const execute = async (interaction: CommandInteraction<CacheType>) => {
  if (interaction.user.id !== BOT_OWNER_ID) {
    await interaction.reply({
      content: 'Only bot owner can use this command',
    })
    return
  }

  await interaction.reply({
    content: 'Restart server in process',
  })

  client.destroy()
  client.login(DISCORD_TOKEN)
  console.log('Restart successfully')

  await updateSlashCommand(Array.from(commandsList.values()))
}

export default {
  name: 'restart',
  data: new SlashCommandBuilder()
    .setName('restart')
    .setDescription('Restart bot server'),
  execute,
} as ICommand
