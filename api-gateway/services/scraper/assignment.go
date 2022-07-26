package scraper

import (
	"api-gateway/models"
	"api-gateway/types"
	"api-gateway/utils"
	"time"

	"gorm.io/gorm"
)

type AssignmentService struct {
	db            *gorm.DB
	courseService *CourseService
}

func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{
		db:            db,
		courseService: NewCourseService(db),
	}
}

func convertToAssignmentView(assignment models.Assignment) types.ShortAssignment {
	return types.ShortAssignment{
		Title:   assignment.Title,
		Href:    assignment.Href,
		DueDate: assignment.DueDate.Format(time.RFC3339),
	}
}

func (c *AssignmentService) GetAssignments(id string, page int, limit int) (types.AssignmentView, error) {
	found, err := c.courseService.IsCourseIdExists(id)
	var res types.AssignmentView

	if err != nil {
		return res, utils.CreateError(500, err.Error())
	}

	if !found {
		return res, utils.CreateError(404, "Not found, maybe api owner does not attend this course")
	}

	var all_assignment []models.Assignment
	var total int64

	tx := c.db.Model(&models.Assignment{}).Where(&models.Assignment{
		CourseID: id,
	}).Count(&total).Order("due_date DESC").Offset(utils.GetOffset(page, limit)).Limit(limit).Find(&all_assignment)

	if tx.Error != nil {
		return res, utils.CreateError(500, tx.Error.Error())
	}

	var assignment_view []types.ShortAssignment

	for _, a := range all_assignment {
		assignment_view = append(assignment_view, convertToAssignmentView(a))
	}

	return types.AssignmentView{
		Assignments: assignment_view,
		Metadata: types.PaginateMetadata{
			CurrentPage: page,
			TotalPages:  utils.GetTotalPages(total, limit),
			TotalItems:  int(total),
		},
	}, nil
}
