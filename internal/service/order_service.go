package service

import (
	"frappuccino/models"
)

type OrderRepository interface {
	ReadOrderOfDal() (map[string]models.Order, string, int)
	UpdateOrderOfDal(item models.Order) (string, int)
	DeleteOrderOfDal(item models.Order) (string, int)
	AddOrderOfDal(item models.Order) (string, int)
}
