package service

import (
	"net/http"

	"frappuccino/models"
)

type MenuRepository interface {
	ReadMenuOfDal() (map[string]models.MenuItem, string, int)
	AddMenuOfDal(item models.MenuItem) (string, int)
	UpdateMenuOfDal(item models.MenuItem) (string, int)
	DeleteMenuOfDal(id string) (string, int)
}

type MenuServiceImpl struct {
	MenuRepository      MenuRepository
	InventoryRepository InventoryRepository
}

func NewMenuService(menuRepository MenuRepository, inventoryRepository InventoryRepository) *MenuServiceImpl {
	return &MenuServiceImpl{
		MenuRepository:      menuRepository,
		InventoryRepository: inventoryRepository,
	}
}

func (svc *MenuServiceImpl) ReadMenuOfSvc() (map[string]models.MenuItem, string, int) {
	return svc.MenuRepository.ReadMenuOfDal()
}

func (svc *MenuServiceImpl) ReadMenuOfSvcByID(id string) (models.MenuItem, string, int) {
	items, msg, code := svc.MenuRepository.ReadMenuOfDal()
	if code != http.StatusOK {
		return models.MenuItem{}, msg, code
	}
	if item, ok := items[id]; ok {
		return item, "Success", http.StatusOK
	}

	return models.MenuItem{}, "Service: Not Found", http.StatusNotFound
}

func (svc *MenuServiceImpl) AddMenuOfSvc(itemMenu models.MenuItem) (string, int) {
	items, msg, code := svc.MenuRepository.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[itemMenu.ID]; exists {
		return "Service: according to this ID, there is already a menu item", http.StatusBadRequest
	}

	invItems, msg, code := svc.InventoryRepository.ReadInventoryOfDal()

	for _, ingredient := range itemMenu.Ingredients {
		if _, exists := invItems[ingredient.IngredientID]; !exists {
			return "Service: there is no ingredient for such an ID: " + ingredient.IngredientID, http.StatusBadRequest
		}
	}
	return svc.MenuRepository.AddMenuOfDal(itemMenu)
}

func (svc *MenuServiceImpl) UpdateMenuOfSvc(itemMenu models.MenuItem) (string, int) {
	items, msg, code := svc.MenuRepository.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}
	if _, exists := items[itemMenu.ID]; !exists {
		return "Service: there is no item menu", http.StatusBadRequest
	}
	invItems, msg, code := svc.InventoryRepository.ReadInventoryOfDal()

	for _, ingredient := range itemMenu.Ingredients {
		if _, exists := invItems[ingredient.IngredientID]; !exists {
			return "Service: there is no ingredient for such an ID: " + ingredient.IngredientID, http.StatusBadRequest
		}
	}
	return svc.MenuRepository.UpdateMenuOfDal(itemMenu)
}

func (svc *MenuServiceImpl) DeleteMenuOfSvc(id string) (string, int) {
	items, msg, code := svc.MenuRepository.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := items[id]; !exists {
		return "Service: there is no item menu", http.StatusBadRequest
	}
	return svc.MenuRepository.DeleteMenuOfDal(id)
}
