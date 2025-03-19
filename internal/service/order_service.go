package service

import (
	"frappuccino/internal/service/utilsService"
	"frappuccino/models"
	"net/http"
	"sort"
	"time"
)

type OrderRepository interface {
	ReadOrderOfDal() (map[string]models.Order, string, int)
	UpdateOrderOfDal(item models.Order) (string, int)
	DeleteOrderOfDal(item models.Order) (string, int)
	AddOrderOfDal(item models.Order) (string, int)
}

type OrderServiceImpl struct {
	repoOrder OrderRepository
	repoMenu  MenuRepository
	repoInv   InventoryRepository
}

func NewOrderService(repoOrder OrderRepository, repoMenu MenuRepository, repoInv InventoryRepository) *OrderServiceImpl {
	return &OrderServiceImpl{repoOrder, repoMenu, repoInv}
}

func (svc *OrderServiceImpl) ReadOrderOfService() (map[string]models.Order, string, int) {
	return svc.repoOrder.ReadOrderOfDal()
}

func (svc *OrderServiceImpl) AddOrderOfService(order models.Order) (string, int) {
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
	order.CreatedAt = time.Now().String()
	order.Status = "open"

	return svc.repoOrder.AddOrderOfDal(order)
}

func (svc *OrderServiceImpl) UpdateOrderOfService(orderUpdate models.Order) (string, int) {
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

func (svc *OrderServiceImpl) DeleteOrderOfService(id string) (string, int) {
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

func (svc *OrderServiceImpl) ReadOrderOfServiceByID(id string) (models.Order, string, int) {
	items, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return models.Order{}, msg, code
	}
	if order, ok := items[id]; ok {
		return order, "Success", http.StatusOK
	}
	return models.Order{}, "Service: Not Found", http.StatusNotFound
}

func (svc *OrderServiceImpl) CloseOrderOfService(id string) (string, int) {
	listOrder, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return msg, code
	}

	if _, exists := listOrder[id]; !exists {
		return "Service: There is no order for this ID", http.StatusBadRequest
	}
	if listOrder[id].Status != "open" {
		return "Service: This order is not open", http.StatusBadRequest
	}
	order := listOrder[id]
	order.Status = "closed"
	return svc.repoOrder.UpdateOrderOfDal(order)
}

func (svc *OrderServiceImpl) TotalSalesOfSvc() (float64, string, int) {
	items, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return 0, msg, code
	}

	listMenu, msg, code := svc.repoMenu.ReadMenuOfDal()
	if code != http.StatusOK {
		return 0, msg, code
	}

	var totalSales float64 = 0

	for _, order := range items {
		if order.Status != "closed" {
			continue
		}
		for _, orderItem := range order.Items {
			totalSales += listMenu[orderItem.ProductID].Price * float64(orderItem.Quantity)
		}
	}
	return totalSales, "Success", http.StatusOK
}

func (svc *OrderServiceImpl) PopularItemsOfSvc() ([]models.OrderItem, string, int) {
	items, msg, code := svc.repoOrder.ReadOrderOfDal()
	if code != http.StatusOK {
		return nil, msg, code
	}

	var popularItems map[string]int

	for _, order := range items {
		if order.Status != "closed" {
			continue
		}
		for _, orderItem := range order.Items {
			if _, exists := popularItems[orderItem.ProductID]; exists {
				quantity := popularItems[orderItem.ProductID] + orderItem.Quantity
				popularItems[orderItem.ProductID] = quantity
			}
			popularItems[orderItem.ProductID] = orderItem.Quantity
		}
	}

	var orderItems []models.OrderItem
	for id, quantity := range popularItems {
		orderItems = append(orderItems, models.OrderItem{id, quantity})
	}

	sort.Slice(orderItems, func(i, j int) bool {
		return orderItems[i].Quantity > orderItems[j].Quantity // Меняй `<` для сортировки по возрастанию
	})

	return orderItems, "Success", http.StatusOK
}
