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
	CreateApartment(apartments []models.Apartment) error
}

func (repo *newAppartmentRepositoryImpl) CreateApartment(apartments []models.Apartment) error {
	return repo.Connection.Table(models.AppartmentsTable()).Create(&apartments).Error
}
