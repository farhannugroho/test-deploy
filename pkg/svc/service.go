package svc

import (
	"food_delivery_api/pkg/model"
	"food_delivery_api/pkg/storage/mysql"
)

type Service interface {
	// User
	AddUser(model.User) (model.User, error)
	GetUsers() ([]model.User, error)
	GetUser(model.User) (model.User, error)
	EditUser(model.User) (model.User, error)
	RemoveUser(model.User) (model.User, error)

	// Merchant
	AddMerchant(model.Merchant) (model.Merchant, error)
	GetMerchants() ([]model.Merchant, error)
	GetMerchant(model.Merchant) (model.Merchant, error)
	EditMerchant(model.Merchant) (model.Merchant, error)
	RemoveMerchant(model.Merchant) (model.Merchant, error)
}

type service struct {
	rmy mysql.RepositoryMySQL
}

func NewService(rmy mysql.RepositoryMySQL) Service {
	return &service{rmy}
}
