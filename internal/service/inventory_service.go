package service

import (
	"frappuccino/models"
	"net/http"
)

type InventoryRepository interface {
	ReadInventoryOfDal() ([]models.InventoryItem, string, int)
	AddInventoryOfDal(item models.InventoryItem) (string, int)
	UpdateInventoryOfDal(itemUpdate models.InventoryItem) (string, int)
	DeleteInventoryOfDal(id string) (string, int)
}

type InventoryService struct {
	InventoryRepository InventoryRepository
}

func NewInventoryService(inventoryRepository InventoryRepository) *InventoryService {
	return &InventoryService{InventoryRepository: inventoryRepository}
}

func (svc *InventoryService) AddInventoryOfSvc(item models.InventoryItem) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()

	if code != http.StatusOK {
		return "Repository: " + msg, code
	}

	for _, item := range items {
		if item.IngredientID == item.IngredientID {
			return "According to this ID, there is an inventory ingredient", 400
		}
	}
	msg, code = svc.InventoryRepository.AddInventoryOfDal(item)
	if code != http.StatusCreated {
		return "Repository: " + msg, code
	}
	return "Success", http.StatusCreated
}

func (svc *InventoryService) UpdateInventoryOfSvc(itemUpdate models.InventoryItem) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()

	if code != http.StatusOK {
		return "Repository: " + msg, code
	}

	for _, item := range items {
		if item.IngredientID == item.IngredientID {
			msg, code = svc.InventoryRepository.UpdateInventoryOfDal(item)
			if code != http.StatusOK {
				return "Repository: " + msg, code
			}
			return "Success", http.StatusOK
		}
	}
	return "Service: There is no inventory ingredient for such an ID", http.StatusNotFound
}

func (svc *InventoryService) DeleteInventoryOfDal(id string) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return "Repository: " + msg, code
	}

	for _, item := range items {
		if item.IngredientID == id {
			msg, code = svc.InventoryRepository.DeleteInventoryOfDal(id)
			if code != http.StatusOK {
				return "Repository: " + msg, code
			}
			return "No Content", http.StatusNoContent
		}
	}
	return "Service: There is no inventory ingredient for this ID", http.StatusNotFound
}

func (svc *InventoryService) ReadInventoryOfSvc() ([]models.InventoryItem, string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return items, "Repository: " + msg, code
	}
	return items, "Success", http.StatusOK
}

func (svc *InventoryService) ReadInventoryOfSvcById(id string) (models.InventoryItem, string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return models.InventoryItem{}, "Repository: " + msg, code
	}
	for _, item := range items {
		if item.IngredientID == id {
			return item, "Success", http.StatusOK
		}
	}
	return models.InventoryItem{}, "Service: There is no inventory for this ID", http.StatusNotFound
}
