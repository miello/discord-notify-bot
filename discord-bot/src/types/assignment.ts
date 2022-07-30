export interface IAssignment {
  title: string
  href: string
  dueDate: string
}

export interface IOverviewAssignment extends IAssignment {
  courseTitle: string
}
