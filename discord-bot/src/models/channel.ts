import mongoose from 'mongoose'

const guildSchema = new mongoose.Schema({
  guildId: String,
  channelId: String,
  courseId: [String],
})

export const Guild = mongoose.model('Guild', guildSchema)
