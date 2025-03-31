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

func (invHandle *InventoryHandler) AddInventoryOfHandle() http.HandlerFunc {
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

		msg, code = invHandle.InventoryService.AddInventoryOfSvc(newInventory)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (invHandle *InventoryHandler) UpdateInventoryOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newInventory models.InventoryItem
		msg, code := utilsHandl.ParseJSONBody(r, &newInventory)
		if id := strings.TrimPrefix(r.URL.Path, "/inventory/"); id != newInventory.IngredientID {
			utilsHandl.SendJSONMessages(w, "Handler: The ID in the url and in the request body cannot be different", http.StatusBadRequest)
			return
		}
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, code)
			return
		}

		if msg = utilsHandl.ValidateInventory(newInventory); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, http.StatusBadRequest)
			return
		}

		msg, code = invHandle.InventoryService.UpdateInventoryOfSvc(newInventory)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (invHandle *InventoryHandler) DeleteInventoryOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "Handler: ID is empty", http.StatusBadRequest)
			return
		}

		msg, code := invHandle.InventoryService.DeleteInventoryOfSvc(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (invHandle *InventoryHandler) ReadInventoryOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, msg, code := invHandle.InventoryService.ReadInventoryOfSvc()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		utilsHandl.SendJSONResponse(w, items, code)
		return
	}
}

func (invHandle *InventoryHandler) ReadInventoryOfHandleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/inventory/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}

		item, msg, code := invHandle.InventoryService.ReadInventoryOfSvcById(id)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		utilsHandl.SendJSONResponse(w, item, code)
		return
	}
}
