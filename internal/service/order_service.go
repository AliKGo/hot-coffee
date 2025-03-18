package service

import (
	"frappuccino/internal/service/utilsService"
	"frappuccino/models"
	"net/http"
)

type OrderRepository interface {
	ReadOrderOfDal() (map[string]models.Order, string, int)
	UpdateOrderOfDal(item models.Order) (string, int)
	DeleteOrderOfDal(item models.Order) (string, int)
	AddOrderOfDal(item models.Order) (string, int)
}

type OrderService struct {
	repoOrder OrderRepository
	repoMenu  MenuRepository
	repoInv   InventoryRepository
}

type OrderServiceImpl interface {
	ReadOrderOfService() (map[string]models.Order, string, int)
	AddOrderOfService(order models.Order) (string, int)
	UpdateOrderOfService(orderUpdate models.Order) (string, int)
	DeleteOrderOfService(id string) (string, int)
	ReadOrderOfServiceByID(id string) (models.Order, string, int)
}

func NewOrderService(repoOrder OrderRepository, repoMenu MenuRepository, repoInv InventoryRepository) *OrderService {
	return &OrderService{repoOrder, repoMenu, repoInv}
}

func (svc *OrderService) ReadOrderOfService() (map[string]models.Order, string, int) {
	return svc.repoOrder.ReadOrderOfDal()
}

func (svc *OrderService) AddOrderOfService(order models.Order) (string, int) {
	listMenu, msg, code := svc.repoMenu.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	listInv, msg, code := svc.repoInv.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	var listUpdate []string

	for _, orderItem := range order.Items {
		if _, exists := listMenu[orderItem.ProductID]; !exists {
			return "Service: There is no menu item for this ID", http.StatusBadRequest
		}
		for _, menuItemIngredient := range listMenu[orderItem.ProductID].Ingredients {
			if _, exists := listInv[menuItemIngredient.IngredientID]; !exists {
				return "Service: There is no inventory item for this ID", http.StatusBadRequest
			}

			if listInv[menuItemIngredient.IngredientID].Quantity < menuItemIngredient.Quantity*float64(orderItem.Quantity) {
				return "Service: There is not enough inventory", http.StatusBadRequest
			}
			newItemInv := listInv[menuItemIngredient.IngredientID]
			listUpdate = append(listUpdate, menuItemIngredient.IngredientID)
			newItemInv.Quantity = listInv[menuItemIngredient.IngredientID].Quantity - menuItemIngredient.Quantity*float64(orderItem.Quantity)
			listInv[menuItemIngredient.IngredientID] = newItemInv
		}
	}

	for _, item := range listUpdate {
		if msg, code = svc.repoInv.UpdateInventoryOfDal(listInv[item]); code != http.StatusOK {
			return msg, code
		}
	}

	listOrder, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := listOrder[order.ID]; exists {
		for exists {
			order.ID = utilsService.GenerateRandomString()
			if _, exists = listOrder[order.ID]; !exists {
				break
			}
		}
	}

	return svc.repoOrder.AddOrderOfDal(order)
}

func (svc *OrderService) UpdateOrderOfService(orderUpdate models.Order) (string, int) {
	listOrder, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	order := listOrder[orderUpdate.ID]

	if _, exists := listOrder[orderUpdate.ID]; !exists {
		return "Service: There is no order for this ID", http.StatusBadRequest
	}

	listMenu, msg, code := svc.repoMenu.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	listInv, msg, code := svc.repoInv.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	listUpdate := []string{}

	for _, orderItem := range order.Items {
		if _, exists := listMenu[orderItem.ProductID]; !exists {
			return "Service: There is no menu item for this ID", http.StatusBadRequest
		}
		for _, menuItemIngredient := range listMenu[orderItem.ProductID].Ingredients {
			if _, exists := listInv[menuItemIngredient.IngredientID]; !exists {
				return "Service: There is no inventory item for this ID", http.StatusBadRequest
			}
			newItemInv := listInv[menuItemIngredient.IngredientID]
			listUpdate = append(listUpdate, menuItemIngredient.IngredientID)
			newItemInv.Quantity = listInv[menuItemIngredient.IngredientID].Quantity + menuItemIngredient.Quantity*float64(orderItem.Quantity)
			listInv[menuItemIngredient.IngredientID] = newItemInv
		}
	}

	listAddorder := []string{}

	for _, orderItem := range orderUpdate.Items {
		if _, exists := listMenu[orderItem.ProductID]; !exists {
			return "Service: There is no menu item for this ID", http.StatusBadRequest
		}

		for _, menuItemIngredient := range listMenu[orderItem.ProductID].Ingredients {
			if _, exists := listInv[menuItemIngredient.IngredientID]; !exists {
				return "Service: There is no inventory item for this ID", http.StatusBadRequest
			}

			if listInv[menuItemIngredient.IngredientID].Quantity < menuItemIngredient.Quantity*float64(orderItem.Quantity) {
				return "Service: There is not enough inventory", http.StatusBadRequest
			}
			newItemInv := listInv[menuItemIngredient.IngredientID]
			listAddorder = append(listAddorder, menuItemIngredient.IngredientID)
			newItemInv.Quantity = listInv[menuItemIngredient.IngredientID].Quantity - menuItemIngredient.Quantity*float64(orderItem.Quantity)
			listInv[menuItemIngredient.IngredientID] = newItemInv
		}
	}

	for _, item := range listUpdate {
		if msg, code = svc.repoInv.UpdateInventoryOfDal(listInv[item]); code != http.StatusOK {
			return msg, code
		}
	}

	for _, item := range listAddorder {
		if msg, code = svc.repoInv.UpdateInventoryOfDal(listInv[item]); code != http.StatusOK {
			return msg, code
		}
	}

	return svc.repoOrder.UpdateOrderOfDal(orderUpdate)
}

func (svc *OrderService) DeleteOrderOfService(id string) (string, int) {
	items, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	order := items[id]

	if _, exists := items[id]; !exists {
		return "Service: There is no order for this ID", http.StatusBadRequest
	}

	listMenu, msg, code := svc.repoMenu.ReadMenuOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	listInv, msg, code := svc.repoInv.ReadInventoryOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	listUpdate := []string{}
	for _, orderItem := range order.Items {
		if _, exists := listMenu[orderItem.ProductID]; !exists {
			return "Service: There is no menu item for this ID", http.StatusBadRequest
		}
		for _, menuItemIngredient := range listMenu[orderItem.ProductID].Ingredients {
			if _, exists := listInv[menuItemIngredient.IngredientID]; !exists {
				return "Service: There is no inventory item for this ID", http.StatusBadRequest
			}
			newItemInv := listInv[menuItemIngredient.IngredientID]
			listUpdate = append(listUpdate, menuItemIngredient.IngredientID)
			newItemInv.Quantity = listInv[menuItemIngredient.IngredientID].Quantity + menuItemIngredient.Quantity*float64(orderItem.Quantity)
			listInv[menuItemIngredient.IngredientID] = newItemInv
		}
	}

	for _, item := range listUpdate {
		if msg, code = svc.repoInv.UpdateInventoryOfDal(listInv[item]); code != http.StatusOK {
			return msg, code
		}
	}
	return svc.repoOrder.DeleteOrderOfDal(order)
}

func (svc *OrderService) ReadOrderOfServiceByID(id string) (models.Order, string, int) {
	items, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return models.Order{}, msg, code
	}
	if order, ok := items[id]; ok {
		return order, "Success", http.StatusOK
	}
	return models.Order{}, "Service: Not Found", http.StatusNotFound
}
