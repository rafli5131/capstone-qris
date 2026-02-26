package repository

import (
	"capstone-qris/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ApiClientRepository struct {
	Repository[entity.ApiClient]
	Log *logrus.Logger
}

func NewApiClientRepository(db *gorm.DB, log *logrus.Logger) *ApiClientRepository {
	return &ApiClientRepository{
		Repository: Repository[entity.ApiClient]{DB: db},
		Log:        log,
	}
}

func (r *ApiClientRepository) FindByClientID(db *gorm.DB, clientID string) (*entity.ApiClient, error) {
	var client entity.ApiClient
	if err := db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
