import { Collection } from 'discord.js'
import { ICommand } from './types/command'
import { readdirSync } from 'fs'
import { join } from 'path'

const commandsPath = join(__dirname, 'commands')
const commandFiles = readdirSync('./src/commands').filter((file) =>
  file.endsWith('.ts')
)

const commandsList = new Collection<string, ICommand>()

for (const file of commandFiles) {
  const filePath = join(
    commandsPath,
    file.substring(0, file.length - 3) + '.js'
  )

  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const command: ICommand = require(filePath).default
  commandsList.set(command.data.name, command)
}

export { commandsList }
