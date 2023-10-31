package repositories

import (
	"fortest/models"

	"gorm.io/gorm"
)

type objectRepositoryImpl struct {
	Connection *gorm.DB
}

func NewObjectsRepository(conn *gorm.DB) ObjectRepository {
	return &objectRepositoryImpl{Connection: conn}
}

type ObjectRepository interface {
	CreateObject(*models.Object) error
	DeleteObjectByID(int) error
	GetAllObjects(filter models.ObjectFilter, limit int, offset int) ([]models.Object, []models.Apartment, int64, error)
	GetPopularObjects(limit int, offset int) ([]models.Object, []models.Apartment, int64, error)
	GetObjectByID(id int) (models.Object, error)
	GetObjectsByDeveloperID(devID uint) ([]models.Object, error)
	UpdateObjectWithApartments(obj *models.Object) error
}

func (repo *objectRepositoryImpl) CreateObject(object *models.Object) error {
	return repo.Connection.Table(models.ObjectsTable()).Create(object).Error
}

func (repo *objectRepositoryImpl) DeleteObjectByID(id int) error {
	if err := repo.Connection.Delete(&models.Object{ID: uint(id)}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *objectRepositoryImpl) GetAllObjects(filter models.ObjectFilter, limit int, offset int) ([]models.Object, []models.Apartment, int64, error) {
	var objects []models.Object
	var apartments []models.Apartment
	var l int64
	query := repo.Connection.Table(models.ObjectsTable())

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}
	if filter.MinPrice != 0 {
		query = query.Where("min_price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice != 0 {
		query = query.Where("max_price <= ?", filter.MaxPrice)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Find(&objects).Limit(limit).Offset(offset).Count(&l).Error; err != nil {
		return nil, nil, 0, err
	}

	if err := repo.Connection.Table(models.AppartmentsTable()).Find(&apartments).Error; err != nil {
		return nil, nil, 0, err
	}

	return objects, apartments, l, nil
}

func (repo *objectRepositoryImpl) GetPopularObjects(limit int, offset int) ([]models.Object, []models.Apartment, int64, error) {
	var objects []models.Object
	var apartments []models.Apartment
	var count int64

	query := repo.Connection.Table(models.ObjectsTable()).Where("popular = ?", true)

	if err := query.Find(&objects).Error; err != nil {
		return nil, nil, 0, err
	}

	if err := repo.Connection.Table(models.AppartmentsTable()).Limit(limit).Offset(offset).Find(&apartments).Error; err != nil {
		return nil, nil, 0, err
	}

	amount := query.Count(&count)
	if amount.Error != nil {
		return nil, nil, 0, amount.Error
	}

	return objects, apartments, count, nil
}

func (repo *objectRepositoryImpl) GetObjectByID(id int) (models.Object, error) {
	var object models.Object
	err := repo.Connection.Table(models.ObjectsTable()).Preload("Apartments").Preload("Developer").Where("id = ?", id).First(&object).Error
	if err != nil {
		return models.Object{}, err
	}
	return object, nil
}

func (repo *objectRepositoryImpl) GetObjectsByDeveloperID(devID uint) ([]models.Object, error) {
	var objects []models.Object
	err := repo.Connection.Where("developer_id = ?", devID).Select("created_at, address, name, status").Find(&objects).Error
	return objects, err
}

func (repo *objectRepositoryImpl) UpdateObjectWithApartments(obj *models.Object) error {
	tx := repo.Connection.Begin()

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, apt := range obj.Apartments {
		if err := tx.Save(&apt).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
