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

func convertToAssignmentView(assignment models.Assignment) types.AssignmentView {
	return types.AssignmentView{
		Title:   assignment.Title,
		Href:    assignment.Href,
		DueDate: assignment.DueDate.Format(time.RFC3339),
	}
}

func (c *AssignmentService) GetAssignments(id string, page int, limit int) ([]types.AssignmentView, error) {
	found, err := c.courseService.IsCourseIdExists(id)

	if err != nil {
		return nil, utils.CreateError(500, err.Error())
	}

	if !found {
		return nil, utils.CreateError(404, "Not found, maybe api owner does not attend this course")
	}

	var all_assignment []models.Assignment

	tx := c.db.Where(&models.Assignment{
		CourseID: id,
	}).Offset(utils.GetOffset(page, limit)).Take(limit).Order("due_date DESC").Find(&all_assignment)

	if tx.Error != nil {
		return nil, utils.CreateError(500, tx.Error.Error())
	}

	var assignment_view []types.AssignmentView

	for _, a := range all_assignment {
		assignment_view = append(assignment_view, convertToAssignmentView(a))
	}

	return assignment_view, nil
}
