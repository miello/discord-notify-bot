import { CommandInteraction, CacheType } from 'discord.js'

export const extractInteractiveInfo = (info: CommandInteraction<CacheType>) => {
  const raw_input = (info.options.data[0].value as string) || ''
  const split_input = raw_input.split(' ')
  const courseId = split_input[0]
  split_input.splice(0, 1)

  const title = split_input.join(' ')

  return [courseId, title]
}
