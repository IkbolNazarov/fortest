package services

import (
	"fortest/models"
	"fortest/repositories"
)

type NewAppartmentServices interface {
	AddApartment(apartments []models.Apartment) error
}

type appartmentServiceImpl struct {
	repo repositories.NewAppartmentRepository
}

func NewAppartmentService(repo repositories.NewAppartmentRepository) NewAppartmentServices {
	return &appartmentServiceImpl{repo: repo}
}

func (s *appartmentServiceImpl) AddApartment(apartments []models.Apartment) error {
	return s.repo.CreateApartment(apartments)
}
