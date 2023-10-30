package services

import (
	"fmt"
	"fortest/models"
	"fortest/repositories"
)

type NewObjectServices interface {
	CreateObject(*models.Object) error
	DeleteObjectByID(int) error
	GetAllObjects(filter models.ObjectFilter, pagination *models.Pagination) ([]models.Object, []models.Apartment, int64, error)
	GetPopularObjects(pagination *models.Pagination) ([]models.Object, []models.Apartment, int64, error)
	GetObjectByID(id int) (models.Object, error)
	GetDeveloperObjectsSummaryByID(devID uint) (models.DeveloperObjectsResponse, error)
	UpdateObjectWithApartments(obj *models.Object) error
}

type objectServiceImpl struct {
	repo repositories.ObjectRepository
}

func NewObjectService(repo repositories.ObjectRepository) NewObjectServices {
	return &objectServiceImpl{repo: repo}
}

func (s *objectServiceImpl) CreateObject(object *models.Object) error {
	return s.repo.CreateObject(object)
}

func (s *objectServiceImpl) DeleteObjectByID(id int) error {
	return s.repo.DeleteObjectByID(id)
}

func (s *objectServiceImpl) GetAllObjects(filter models.ObjectFilter, pagination *models.Pagination) ([]models.Object, []models.Apartment, int64, error) {
	offset := (pagination.Page - 1) * pagination.Limit
	return s.repo.GetAllObjects(filter, pagination.Limit, offset)
}

func (s *objectServiceImpl) GetPopularObjects(pagination *models.Pagination) ([]models.Object, []models.Apartment, int64, error) {
	offset := (pagination.Page - 1) * pagination.Limit
	return s.repo.GetPopularObjects(pagination.Limit, offset)
}

func (s *objectServiceImpl) GetObjectByID(id int) (models.Object, error) {
	return s.repo.GetObjectByID(id)
}

func (s *objectServiceImpl) GetDeveloperObjectsSummaryByID(devID uint) (models.DeveloperObjectsResponse, error) {
	var response models.DeveloperObjectsResponse
	objects, err := s.repo.GetObjectsByDeveloperID(devID)
	if err != nil {
		return response, err
	}

	for _, obj := range objects {
		summary := models.DeveloperObjectsSummary{
			CreationDate: obj.CreationDate,
			Address:      obj.Address,
			Name:         obj.Name,
			Status:       obj.Status,
		}
		response.Objects = append(response.Objects, summary)
	}
	response.DeveloperID = devID
	return response, nil
}

func (s *objectServiceImpl) UpdateObjectWithApartments(obj *models.Object) error {
	if err := s.repo.UpdateObjectWithApartments(obj); err != nil {
		return fmt.Errorf("unable to update object and its apartments: %w", err)
	}
	return nil
}
