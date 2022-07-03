import { Collection } from 'discord.js'

const subscriber = new Collection<string, Set<string>>()

export const subscribe = (guildId: string, channelId: string) => {
  const prev = subscriber.get(guildId) || new Set()
  prev.add(channelId)
  subscriber.set(guildId, prev)
}

export const unsubscribe = (guildId: string, channelId: string) => {
  if (!subscriber.get(guildId)) return

  const prev = subscriber.get(guildId) || new Set()
  prev.delete(channelId)
  subscriber.set(guildId, prev)
}

export const getSubscriberList = () => {
  return Array.from(subscriber.entries())
}
