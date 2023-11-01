package repositories

import (
	"fortest/models"

	"gorm.io/gorm"
)

type builderRepositoryImpl struct {
	Connection *gorm.DB
}

func NewBuilderRepository(conn *gorm.DB) BuilderRepository {
	return &builderRepositoryImpl{Connection: conn}
}

type BuilderRepository interface {
	CreateDeveloper(developer models.Developer)  (models.Developer,error)
	GetDeveloperWithObjectsByID(devID uint) (models.DeveloperWithObjectsResponse, error)
	GetDeveloperByID(id uint) (*models.Developer, error)
}

func (repo *builderRepositoryImpl) CreateDeveloper(developer models.Developer)  (models.Developer,error) {
	err :=  repo.Connection.Table(models.DevelopersTable()).Create(&developer).Error
	return developer, err 
}

func (repo *builderRepositoryImpl) GetDeveloperByID(id uint) (*models.Developer, error) {
	var developer models.Developer
	if err := repo.Connection.First(&developer, id).Error; err != nil {
		return nil, err
	}
	return &developer, nil
}

func (repo *builderRepositoryImpl) GetDeveloperWithObjectsByID(devID uint) (models.DeveloperWithObjectsResponse, error) {
	var developer models.Developer
	var objects []models.ObjectResp
	var response models.DeveloperWithObjectsResponse


	err := repo.Connection.First(&developer, devID).Error
	if err != nil {
		return response, err
	}

	err = repo.Connection.Table(models.ObjectsTable()).Where("developer_id = ?", devID).Find(&objects).Error
	if err != nil {
		return response, err
	}

	response.DeveloperInfo = developer
	response.Objects = objects

	return response, nil
}







