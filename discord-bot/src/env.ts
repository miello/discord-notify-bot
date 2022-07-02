import { config } from 'dotenv'

config()

const DISCORD_TOKEN = process.env.DISCORD_TOKEN || ''
const APPLICATION_ID = process.env.APPLICATION_ID || ''

export { DISCORD_TOKEN, APPLICATION_ID }
