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
}

type OrderHandler struct {
	menuService OrderService
}

func NewOrderHandler(menuService OrderService) *OrderHandler {
	return &OrderHandler{menuService: menuService}
}

func (h *OrderHandler) ReadOrderOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, msg, code := h.menuService.ReadOrderOfService()
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
		order, msg, code := h.menuService.ReadOrderOfServiceByID(id)
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

		msg, code = h.menuService.AddOrderOfService(newOrder)
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

		msg, code = h.menuService.UpdateOrderOfService(newOrder)
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
		msg, code := h.menuService.DeleteOrderOfService(id)
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
		msg, code := h.menuService.CloseOrderOfService(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}
