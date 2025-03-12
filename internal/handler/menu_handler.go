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

func (menuHandler *MenuHandler) ReadMenuOfHandl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, msg, code := menuHandler.MenuService.ReadMenuOfSvc()
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, item, code)
		return
	}
}

func (menuHandler *MenuHandler) ReadMenuOfHandlByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		item, msg, code := menuHandler.MenuService.ReadMenuOfSvcByID(id)
		if code != http.StatusOK {
			utilsHandl.SendJSONMessages(w, msg, code)
			return
		}
		utilsHandl.SendJSONResponse(w, item, code)
		return
	}
}

func (menuHandler *MenuHandler) AddMenuOfHandl() http.HandlerFunc {
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

		msg, code = menuHandler.MenuService.AddMenuOfSvc(newMenuItem)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (menuHandler *MenuHandler) UpdateMenuOfHandl() http.HandlerFunc {
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

		msg, code = menuHandler.MenuService.UpdateMenuOfSvc(newMenuItem)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}

func (menuHandler *MenuHandler) DeleteMenuOfHandl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/menu/")
		if id == "" {
			utilsHandl.SendJSONMessages(w, "ID is empty", http.StatusBadRequest)
			return
		}
		msg, code := menuHandler.MenuService.DeleteMenuOfSvc(id)
		utilsHandl.SendJSONMessages(w, msg, code)
		return
	}
}
