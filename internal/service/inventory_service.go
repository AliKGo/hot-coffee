package service

import (
	"net/http"

	"frappuccino/models"
)

type InventoryRepository interface {
	ReadInventoryOfDal() (map[string]models.InventoryItem, string, int)
	AddInventoryOfDal(item models.InventoryItem) (string, int)
	UpdateInventoryOfDal(itemUpdate models.InventoryItem) (string, int)
	DeleteInventoryOfDal(id string) (string, int)
}
type InventoryServiceImpl struct {
	MenuRepository      MenuRepository
	InventoryRepository InventoryRepository
}

func NewInventoryService(inventoryRepository InventoryRepository, menuRepository MenuRepository) *InventoryServiceImpl {
	return &InventoryServiceImpl{InventoryRepository: inventoryRepository, MenuRepository: menuRepository}
}

func (svc *InventoryServiceImpl) AddInventoryOfSvc(item models.InventoryItem) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[item.IngredientID]; exists {
		return "Service: According to this ID, there is an inventory ingredient", http.StatusBadRequest
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

	if _, exists := items[itemUpdate.IngredientID]; !exists {
		return "Service: There is no inventory ingredient for such an ID", http.StatusNotFound
	}

	msg, code = svc.InventoryRepository.UpdateInventoryOfDal(itemUpdate)
	if code != http.StatusOK {
		return msg, code
	}

	return "Success", http.StatusOK
}

func (svc *InventoryServiceImpl) DeleteInventoryOfSvc(id string) (string, int) {
	items, msg, code := svc.InventoryRepository.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[id]; !exists {
		return "Service: There is no inventory ingredient for this ID", http.StatusNotFound
	}

	itemsMenu, msg, code := svc.MenuRepository.ReadMenuOfDal()

	for _, itemMenu := range itemsMenu {
		for _, itemIngredients := range itemMenu.Ingredients {
			if _, data := items[itemIngredients.IngredientID]; data {
				return "Service: This inventory item is used in the menu!", http.StatusBadRequest
			}
		}
	}

	msg, code = svc.InventoryRepository.DeleteInventoryOfDal(id)
	if code != http.StatusOK {
		return msg, code
	}

	return "Service: No Content", http.StatusNoContent
}

func (svc *InventoryServiceImpl) ReadInventoryOfSvc() (map[string]models.InventoryItem, string, int) {
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

	if item, ok := items[id]; ok {
		return item, "Success", http.StatusOK
	}

	return models.InventoryItem{}, "Service: There is no inventory for this ID", http.StatusNotFound
}
