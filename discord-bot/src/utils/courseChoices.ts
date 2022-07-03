import { apiClient } from '../config/axios'

interface ICourse {
  key: string
  courseId: string
  courseTitle: string
  courseHref: string
  semester: number
  year: number
}

export async function getCourseChoices() {
  const resp = await apiClient.get<Array<ICourse>>('/courses')
  const courses = resp.data
  const choices = courses.map((val) => {
    return {
      name: `${val.courseTitle} (${val.semester}/${val.year})`,
      value: val.key,
    }
  })

  return choices
}
