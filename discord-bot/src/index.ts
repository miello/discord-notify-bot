import { DISCORD_TOKEN } from './config/env'
import { updateSlashCommand } from './utils/slashCommand'
import { commandsList } from './config/routes'
import { client } from './config/clientBot'

updateSlashCommand(Array.from(commandsList.values()))

client.login(DISCORD_TOKEN)
