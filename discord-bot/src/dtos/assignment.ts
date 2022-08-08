import { IAssignment, IOverviewAssignment } from '../types/assignment'
import { IPaginationMetadata } from '../types/common'

export interface AssignmentDTO {
  assignments: Array<IAssignment>
  meta: IPaginationMetadata
}

export interface OverviewAssignmentDTO {
  assignments: Array<IOverviewAssignment>
  meta: IPaginationMetadata
}
