package repository

import (
	"gorm.io/gorm"
)

// Repository is a generic base providing basic CRUD operations.
type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Save(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) FindByID(db *gorm.DB, entity *T, id any) error {
	return db.First(entity, "id = ?", id).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}
