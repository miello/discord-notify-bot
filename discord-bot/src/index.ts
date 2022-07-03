import { Client, Intents } from 'discord.js'
import { DISCORD_TOKEN } from './config/env'
import { updateSlashCommand } from './utils/slashCommand'
import { commandsList } from './routes'

const client = new Client({ intents: [Intents.FLAGS.GUILDS] })

client.once('ready', () => {
  console.log('Ready!')
})

client.on('interactionCreate', async (interaction) => {
  if (!interaction.isCommand()) return

  const { commandName } = interaction
  const command = commandsList.get(commandName)

  if (!command) return

  try {
    await command.execute(interaction)
  } catch (err) {
    console.error(err)
    await interaction.reply({
      content: 'There was an error while executing this command!',
      ephemeral: true,
    })
  }
})

updateSlashCommand(Array.from(commandsList.values()))

client.login(DISCORD_TOKEN)
