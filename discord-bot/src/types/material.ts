export interface IFile {
  title: string
  href: string
}

export interface IMaterial {
  folderName: string
  file: Array<IFile>
}
