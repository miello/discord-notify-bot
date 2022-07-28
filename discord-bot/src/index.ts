import { DISCORD_TOKEN } from './config/env'
import { updateSlashCommand } from './utils/slashCommand'
import { commandsList } from './config/routes'
import { client } from './config/clientBot'
import { initDB } from './config/mongo'

updateSlashCommand(Array.from(commandsList.values())).catch(console.error)
initDB().catch(console.error)

client.login(DISCORD_TOKEN)
