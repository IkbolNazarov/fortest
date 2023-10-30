package core

import (
	"fortest/db"
	"fortest/repositories"
	"fortest/services"
	"fortest/validator"
)

type Core struct {
	Validator            *validator.Validator
	NewObjectService     services.NewObjectServices
	NewAppartmentService services.NewAppartmentServices
	NewBuilderService    services.NewBuilderServices
}

func New() *Core {
	core := &Core{}
	core.Validator = validator.New()
	DB := db.GetDB()
	newObjectsRepository := repositories.NewObjectsRepository(DB)
	newBuilderRepository := repositories.NewBuilderRepository(DB)
	newApparmentRepository := repositories.NewAppartmentsRepository(DB)

	core.NewObjectService = services.NewObjectService(newObjectsRepository)
	core.NewBuilderService = services.NewBuilderService(newBuilderRepository)
	core.NewAppartmentService = services.NewAppartmentService(newApparmentRepository)

	return core
}
