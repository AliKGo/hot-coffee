package handler

import (
	"net/http"
	"strings"

	"frappuccino/internal/handler/utilsHandl"
	"frappuccino/models"
)

type MenuService interface {
	ReadMenuOfSvc() (map[string]models.MenuItem, string, int)
	ReadMenuOfSvcByID(id string) (models.MenuItem, string, int)
	AddMenuOfSvc(itemMenu models.MenuItem) (string, int)
	UpdateMenuOfSvc(itemMenu models.MenuItem) (string, int)
	DeleteMenuOfSvc(id string) (string, int)
}

type MenuHandler struct {
	MenuService MenuService
}

func NewMenuHandler(menuService MenuService) *MenuHandler {
	return &MenuHandler{MenuService: menuService}
}

func (menuHandle *MenuHandler) ReadMenuOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, msg, code := menuHandle.MenuService.ReadMenuOfSvc()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, item, code)
		return
	}
}

func (menuHandle *MenuHandler) ReadMenuOfHandleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "Handler: No ID was entered", http.StatusBadRequest)
			return
		}
		item, msg, code := menuHandle.MenuService.ReadMenuOfSvcByID(id)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, item, code)
		return
	}
}

func (menuHandle *MenuHandler) AddMenuOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newMenuItem models.MenuItem
		msg, code := utilsHandl.ParseJSONBody(r, &newMenuItem)

		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		if msg = utilsHandl.ValidateMenu(newMenuItem); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, code)
			return
		}

		msg, code = menuHandle.MenuService.AddMenuOfSvc(newMenuItem)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (menuHandle *MenuHandler) UpdateMenuOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newMenuItem models.MenuItem
		msg, code := utilsHandl.ParseJSONBody(r, &newMenuItem)

		if id := strings.TrimPrefix(r.URL.Path, "/menu/"); id != newMenuItem.ID {
			utilsHandl.SendJSONMessages(w, "Handler: The ID in the url and in the request body cannot be different", http.StatusBadRequest)
			return
		}
		
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}

		if msg = utilsHandl.ValidateMenu(newMenuItem); msg != "OK" {
			utilsHandl.SendJSONMessages(w, "Handler: "+msg, code)
			return
		}

		msg, code = menuHandle.MenuService.UpdateMenuOfSvc(newMenuItem)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (menuHandle *MenuHandler) DeleteMenuOfHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "Handler: ID is empty", http.StatusBadRequest)
			return
		}
		msg, code := menuHandle.MenuService.DeleteMenuOfSvc(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}
