import mongoose from 'mongoose'
import { MONGO_PASSWORD, MONGO_URI, MONGO_USERNAME } from './env'

export async function initDB() {
  console.log('Try to connect to database')
  await mongoose.connect(MONGO_URI, {
    auth: {
      username: MONGO_USERNAME,
      password: MONGO_PASSWORD,
    },
    authSource: 'admin',
  })
  console.log('Connected to database')
}
