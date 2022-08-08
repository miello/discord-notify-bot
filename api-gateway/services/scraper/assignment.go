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
	})

	tx = tx.Count(&total).Order("due_date DESC").Offset(utils.GetOffset(page, limit)).Limit(limit).Find(&all_assignment)

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

func (c *AssignmentService) GetOverviewAssignments(id []string, page int, limit int) (types.OverviewAssignmentView, error) {
	offset := utils.GetOffset(page, limit)
	now_date := time.Now()
	date := time.Now().Add(14 * 24 * time.Hour)

	var res types.OverviewAssignmentView
	var assignment []models.Assignment
	var count int64

	tx := c.db.Model(&models.Assignment{}).Where(gorm.Expr("due_date < ? AND due_date > ?", date, now_date)).Where(gorm.Expr("course_id IN (?)", id))
	if tx.Error != nil {
		return res, utils.CreateError(500, tx.Error.Error())
	}

	tx = tx.Count(&count).Order("due_date asc").Offset(offset).Limit(limit).Preload("Course").Find(&assignment)
	if tx.Error != nil {
		return res, utils.CreateError(500, tx.Error.Error())
	}

	var assignment_view []types.OverviewAssignment = make([]types.OverviewAssignment, 0)

	for _, a := range assignment {
		assignment_view = append(assignment_view, convertToOverviewAssignment(a))
	}

	res = types.OverviewAssignmentView{
		Assignments: assignment_view,
		Metadata: types.PaginateMetadata{
			CurrentPage: page,
			TotalPages:  utils.GetTotalPages(count, limit),
			TotalItems:  int(count),
		},
	}

	return res, nil
}

func convertToOverviewAssignment(a models.Assignment) types.OverviewAssignment {
	return types.OverviewAssignment{
		CourseTitle: a.Course.Title,
		Title:       a.Title,
		Href:        a.Href,
		DueDate:     a.DueDate.Format(time.RFC3339),
	}
}
