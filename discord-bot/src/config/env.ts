import { config } from 'dotenv'

config()

const DISCORD_TOKEN = process.env.DISCORD_TOKEN || ''
const APPLICATION_ID = process.env.APPLICATION_ID || ''
const BASE_API_GATEWAY = process.env.BASE_API_GATEWAY || ''

export { DISCORD_TOKEN, APPLICATION_ID, BASE_API_GATEWAY }
