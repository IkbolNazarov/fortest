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
	CreateApartment(apartments []models.Apartment) ([]models.Apartment, error)
}

func (repo *newAppartmentRepositoryImpl) CreateApartment(apartments []models.Apartment) ([]models.Apartment, error) {
	for _, apartment := range apartments {
		if err := repo.Connection.Table(models.AppartmentsTable()).Create(&apartment).Error; err != nil {
			return nil, err
		}
	}
	return apartments, nil
}
