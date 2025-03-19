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

func (h *OrderHandler) ReadOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, msg, code := h.orderService.ReadOrderOfService()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, orders, code)
		return
	}
}

func (h *OrderHandler) ReadOrderOfHandleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/order/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		order, msg, code := h.orderService.ReadOrderOfServiceByID(id)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		utilsHandl.SendJSONResponse(w, order, code)
		return
	}
}

func (h *OrderHandler) AddOrderOfHandle() http.HandlerFunc {
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

		msg, code = h.orderService.AddOrderOfService(newOrder)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (h *OrderHandler) UpdateOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/order/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
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

		newOrder.ID = id

		msg, code = h.orderService.UpdateOrderOfService(newOrder)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (h *OrderHandler) DeleteOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/order/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		msg, code := h.orderService.DeleteOrderOfService(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (h *OrderHandler) CloseOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/order/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		msg, code := h.orderService.CloseOrderOfService(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (h *OrderHandler) TotalSalesOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if totalSales, msg, code := h.orderService.TotalSalesOfSvc(); code != http.StatusOK {
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

func (h *OrderHandler) PopularItemsOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, msg, code := h.orderService.PopularItemsOfSvc()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		utilsHandl.SendJSONResponse(w, items, code)
		return
	}
}
