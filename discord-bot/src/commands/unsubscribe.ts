import { SlashCommandBuilder } from '@discordjs/builders'
import { CacheType, CommandInteraction } from 'discord.js'
import { ICommand } from '../types/command'
import { unsubscribe } from '../utils/subscribe'

async function execute(interaction: CommandInteraction<CacheType>) {
  unsubscribe(interaction.guildId || '', interaction.channelId || '')

  await interaction.reply({
    content: 'Unsubscribe to daily notification',
  })
}

export default {
  name: 'unsubscribe',
  data: new SlashCommandBuilder()
    .setName('unsubscribe')
    .setDescription('Unsubscribe for daily notification'),
  execute,
} as ICommand
