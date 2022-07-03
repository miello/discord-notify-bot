import { apiClient } from '../config/axios'
import { ICourse } from '../types/course'

export async function getCourseChoices() {
  const resp = await apiClient.get<Array<ICourse>>('/courses')
  const courses = resp.data

  const choices = courses.map((val) => {
    return {
      name: `${val.courseTitle} (${val.semester}/${val.year})`,
      value: `${val.key} ${val.courseTitle}`,
    }
  })

  return choices
}
