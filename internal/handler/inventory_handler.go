package handler

import (
	"net/http"
	"strings"

	"frappuccino/internal/handler/utilsHandl"
	"frappuccino/models"
)

type InventoryService interface {
	AddInventoryOfSvc(item models.InventoryItem) (string, int)
	UpdateInventoryOfSvc(itemUpdate models.InventoryItem) (string, int)
	DeleteInventoryOfSvc(id string) (string, int)
	ReadInventoryOfSvc() (map[string]models.InventoryItem, string, int)
	ReadInventoryOfSvcById(id string) (models.InventoryItem, string, int)
}

type InventoryHandler struct {
	InventoryService InventoryService
}

func NewInventoryHandler(inventoryService InventoryService) *InventoryHandler {
	return &InventoryHandler{InventoryService: inventoryService}
}

func (invHandl *InventoryHandler) AddInventoryOfHandl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newInventory models.InventoryItem
		msg, code := utilsHandl.ParseJSONBody(r, &newInventory)

		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, code)
			return
		}

		if msg = utilsHandl.ValidateInventory(newInventory); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, http.StatusBadRequest)
			return
		}

		msg, code = invHandl.InventoryService.AddInventoryOfSvc(newInventory)
		utilsHandl.SendJSONMessages(w, msg, code)
	}
}

func (invHandl *InventoryHandler) UpdateInventoryOfHandl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newInventory models.InventoryItem
		msg, code := utilsHandl.ParseJSONBody(r, &newInventory)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, code)
			return
		}
		if msg = utilsHandl.ValidateInventory(newInventory); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, http.StatusBadRequest)
			return
		}
		msg, code = invHandl.InventoryService.UpdateInventoryOfSvc(newInventory)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (invHandl *InventoryHandler) DeleteInventoryOfHandl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		msg, code := invHandl.InventoryService.DeleteInventoryOfSvc(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (invHandl *InventoryHandler) ReadInventoryOfHandl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, msg, code := invHandl.InventoryService.ReadInventoryOfSvc()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, items, code)
		return
	}
}

func (invHandl *InventoryHandler) ReadInventoryOfHandlByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}

		item, msg, code := invHandl.InventoryService.ReadInventoryOfSvcById(id)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, item, code)
		return
	}
}
