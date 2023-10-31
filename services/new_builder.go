package services

import (
	"fortest/models"
	"fortest/repositories"
)

type NewBuilderServices interface {
	AddDeveloper(developer models.Developer)  (models.Developer,error)
	GetDeveloperWithObjectsByID(devID uint) (models.DeveloperWithObjectsResponse, error)
	GetDeveloperByID(id uint) (*models.Developer, error)
}

type builderServiceImpl struct {
	repo repositories.BuilderRepository
}

func NewBuilderService(repo repositories.BuilderRepository) NewBuilderServices {
	return &builderServiceImpl{repo: repo}
}

func (s *builderServiceImpl) AddDeveloper(developer models.Developer) (models.Developer,error) {
	return s.repo.CreateDeveloper(developer)
}

func (s *builderServiceImpl) GetDeveloperWithObjectsByID(devID uint) (models.DeveloperWithObjectsResponse, error) {
	return s.repo.GetDeveloperWithObjectsByID(devID)
}

func (s *builderServiceImpl) GetDeveloperByID(id uint) (*models.Developer, error) {
	return s.repo.GetDeveloperByID(id)
}
