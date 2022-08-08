import { config } from 'dotenv'

config()

const DISCORD_TOKEN = process.env.DISCORD_TOKEN || ''
const APPLICATION_ID = process.env.APPLICATION_ID || ''
const BASE_API_GATEWAY = process.env.BASE_API_GATEWAY || ''

const MONGO_HOST = process.env.MONGO_HOST || ''
const MONGO_PORT = process.env.MONGO_PORT || ''
const MONGO_USERNAME = process.env.MONGO_USERNAME || ''
const MONGO_PASSWORD = process.env.MONGO_PASSWORD || ''
const MONGO_DATABASE = process.env.MONGO_DATABASE || ''
const BOT_OWNER_ID = process.env.BOT_OWNER_ID || ''

const MONGO_URI = `mongodb://${MONGO_HOST}:${MONGO_PORT}/${MONGO_DATABASE}`

export {
  DISCORD_TOKEN,
  APPLICATION_ID,
  BASE_API_GATEWAY,
  MONGO_URI,
  MONGO_USERNAME,
  MONGO_PASSWORD,
  BOT_OWNER_ID,
}
