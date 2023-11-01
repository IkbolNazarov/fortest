package repositories

import (
	"fortest/models"
	"gorm.io/gorm"
)

type newAppartmentRepositoryImpl struct {
	Connection *gorm.DB
}

func NewAppartmentsRepository(conn *gorm.DB) NewAppartmentRepository {
	return &newAppartmentRepositoryImpl{Connection: conn}
}

type NewAppartmentRepository interface {
	CreateApartment(apartments models.ApartmentsResponse) (models.ApartmentsResponse, error)
}

func (repo *newAppartmentRepositoryImpl) CreateApartment(apartments models.ApartmentsResponse) (models.ApartmentsResponse, error) {
	for _, apartment := range apartments.Apartments {
		if err := repo.Connection.Table(models.AppartmentsTable()).Create(&apartment).Error; err != nil {
			return models.ApartmentsResponse{}, err
		}
	}
	return apartments, nil
}
