package scraper

import (
	"api-gateway/models"
	"api-gateway/types"
	"api-gateway/utils"

	"gorm.io/gorm"
)

type MaterialService struct {
	db            *gorm.DB
	courseService *CourseService
}

func NewMaterialService(db *gorm.DB) *MaterialService {
	return &MaterialService{
		db:            db,
		courseService: NewCourseService(db),
	}
}

func convertToMaterialView(materials *[]models.Material) []types.MaterialView {
	result_key := map[string][]types.File{}
	for _, material := range *materials {
		result_key[material.FolderName] = append(result_key[material.FolderName], types.File{
			Title: material.Title,
			Href:  material.Href,
		})
	}

	var materialView []types.MaterialView

	for folderName, value := range result_key {
		materialView = append(materialView, types.MaterialView{
			FolderName: folderName,
			File:       value,
		})
	}

	return materialView
}

func (c *MaterialService) GetMaterials(id string, folderName string) ([]types.MaterialView, error) {
	found, err := c.courseService.IsCourseIdExists(id)

	if err != nil {
		return nil, utils.CreateError(500, err.Error())
	}

	if !found {
		return nil, utils.CreateError(404, "Not found, maybe api owner does not attend this course")
	}

	query := models.Material{
		CourseID: id,
	}

	if folderName != "" {
		query.FolderName = folderName
	}

	var raw_material []models.Material

	tx := c.db.Where(&query).Find(&raw_material)
	if tx.Error != nil {
		return nil, utils.CreateError(500, tx.Error.Error())
	}

	return convertToMaterialView(&raw_material), nil
}
