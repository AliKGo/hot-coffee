package service

import (
	"net/http"

	"frappuccino/models"
)

type InventoryRepository interface {
	ReadInventoryOfDal() ([]models.InventoryItem, string, int)
	AddInventoryOfDal(item models.InventoryItem) (string, int)
	UpdateInventoryOfDal(itemUpdate models.InventoryItem) (string, int)
	DeleteInventoryOfDal(id string) (string, int)
}
type InventoryServiceImpl struct {
	InventoryRepository InventoryRepository
}

func NewInventoryService(inventoryRepository InventoryRepository) *InventoryServiceImpl {
	return &InventoryServiceImpl{InventoryRepository: inventoryRepository}
}

func (svc *InventoryServiceImpl) AddInventoryOfSvc(item models.InventoryItem) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()

	if code != http.StatusOK {
		return msg, code
	}

	for _, item := range items {
		if item.IngredientID == item.IngredientID {
			return "Service: According to this ID, there is an inventory ingredient", http.StatusBadRequest
		}
	}
	msg, code = svc.InventoryRepository.AddInventoryOfDal(item)
	if code != http.StatusCreated {
		return msg, code
	}
	return "Success", http.StatusCreated
}

func (svc *InventoryServiceImpl) UpdateInventoryOfSvc(itemUpdate models.InventoryItem) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()

	if code != http.StatusOK {
		return msg, code
	}

	for _, item := range items {
		if item.IngredientID == item.IngredientID {
			msg, code = svc.InventoryRepository.UpdateInventoryOfDal(item)
			if code != http.StatusOK {
				return msg, code
			}
			return "Success", http.StatusOK
		}
	}
	return "Service: There is no inventory ingredient for such an ID", http.StatusNotFound
}

func (svc *InventoryServiceImpl) DeleteInventoryOfSvc(id string) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	for _, item := range items {
		if item.IngredientID == id {
			msg, code = svc.InventoryRepository.DeleteInventoryOfDal(id)
			if code != http.StatusOK {
				return msg, code
			}
			return "Service: No Content", http.StatusNoContent
		}
	}
	return "Service: There is no inventory ingredient for this ID", http.StatusNotFound
}

func (svc *InventoryServiceImpl) ReadInventoryOfSvc() ([]models.InventoryItem, string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return items, msg, code
	}
	return items, "Success", http.StatusOK
}

func (svc *InventoryServiceImpl) ReadInventoryOfSvcById(id string) (models.InventoryItem, string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return models.InventoryItem{}, msg, code
	}
	for _, item := range items {
		if item.IngredientID == id {
			return item, "Success", http.StatusOK
		}
	}
	return models.InventoryItem{}, "Service: There is no inventory for this ID", http.StatusNotFound
}
