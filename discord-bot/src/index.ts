import { DISCORD_TOKEN } from './config/env'
import { updateSlashCommand } from './utils/slashCommand'
import { commandsList } from './config/routes'
import { client } from './config/clientBot'
import { initDB } from './config/mongo'
import mongoose from 'mongoose'

updateSlashCommand(Array.from(commandsList.values())).catch(console.error)
initDB().catch(console.error)

process.on('SIGTERM', () => {
  client.destroy()
  console.log('Bot is shutting down')
  mongoose.connection.close(false, () => {
    console.log('MongoDb connection closed.')
    process.exit(0)
  })
})

client.login(DISCORD_TOKEN)
