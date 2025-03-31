package handler

import (
	"frappuccino/internal/handler/utilsHandl"
	"frappuccino/models"
	"net/http"
	"strings"
)

type OrderService interface {
	ReadOrderOfService() (map[string]models.Order, string, int)
	AddOrderOfService(order models.Order) (string, int)
	UpdateOrderOfService(orderUpdate models.Order) (string, int)
	DeleteOrderOfService(id string) (string, int)
	ReadOrderOfServiceByID(id string) (models.Order, string, int)
	CloseOrderOfService(id string) (string, int)
	TotalSalesOfSvc() (float64, string, int)
	PopularItemsOfSvc() ([]models.OrderItem, string, int)
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(menuService OrderService) *OrderHandler {
	return &OrderHandler{orderService: menuService}
}

func (orderHandle *OrderHandler) ReadOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, msg, code := orderHandle.orderService.ReadOrderOfService()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, orders, code)
		return
	}
}

func (orderHandle *OrderHandler) ReadOrderOfHandleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/orders/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "Handler: No ID was entered", http.StatusBadRequest)
			return
		}
		order, msg, code := orderHandle.orderService.ReadOrderOfServiceByID(id)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		utilsHandl.SendJSONResponse(w, order, code)
		return
	}
}

func (orderHandle *OrderHandler) AddOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newOrder models.Order
		msg, code := utilsHandl.ParseJSONBody(r, &newOrder)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		if msg = utilsHandl.ValidateOrder(newOrder); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, http.StatusBadRequest)
			return
		}

		msg, code = orderHandle.orderService.AddOrderOfService(newOrder)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (orderHandle *OrderHandler) UpdateOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/orders/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "Handler: No ID was entered", http.StatusBadRequest)
			return
		}
		var newOrder models.Order
		msg, code := utilsHandl.ParseJSONBody(r, &newOrder)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		if msg = utilsHandl.ValidateOrder(newOrder); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, http.StatusBadRequest)
			return
		}

		if newOrder.ID != id {
			utilsHandl.SendJSONMessages(w, "Handler: The ID in the url and in the request body cannot be different", http.StatusBadRequest)
			return

		}

		msg, code = orderHandle.orderService.UpdateOrderOfService(newOrder)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (orderHandle *OrderHandler) DeleteOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/orders/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		msg, code := orderHandle.orderService.DeleteOrderOfService(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (orderHandle *OrderHandler) CloseOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/orders/")
		id = strings.TrimSuffix(id, "/close")

		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		
		msg, code := orderHandle.orderService.CloseOrderOfService(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (orderHandle *OrderHandler) TotalSalesOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if totalSales, msg, code := orderHandle.orderService.TotalSalesOfSvc(); code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		} else {
			response := models.Response{
				TotalSales: totalSales,
			}
			utilsHandl.SendJSONResponse(w, response, code)
			return
		}
	}
}

func (orderHandle *OrderHandler) PopularItemsOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, msg, code := orderHandle.orderService.PopularItemsOfSvc()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		utilsHandl.SendJSONResponse(w, items, code)
		return
	}
}
