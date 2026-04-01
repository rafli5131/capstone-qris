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

func (r *ApiClientRepository) FindAll(db *gorm.DB) ([]entity.ApiClient, error) {
	var clients []entity.ApiClient
	if err := db.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ApiClientRepository) Update(db *gorm.DB, client *entity.ApiClient) error {
	return db.Save(client).Error
}
