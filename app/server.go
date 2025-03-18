package app

import (
	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"
	"frappuccino/tools"
	"net/http"
)

func StartServer() {
	tools.InitLogger()
	tools.CheckJsonFils()
	mux := http.NewServeMux()

	invRepo := dal.InventoryFilePath()
	menuRepo := dal.MenuFilePath()

	invSvc := service.NewInventoryService(invRepo, menuRepo)
	invHandl := handler.NewInventoryHandler(invSvc)

	menuSvc := service.NewMenuService(menuRepo, invRepo)
	menuHandl := handler.NewMenuHandler(menuSvc)

	mux.HandleFunc("POST /menu", menuHandl.AddMenuOfHandl())
	mux.HandleFunc("GET /menu", menuHandl.ReadMenuOfHandl())
	mux.HandleFunc("GET /menu/{id}", menuHandl.ReadMenuOfHandlByID())
	mux.HandleFunc("PUT /menu/{id}", menuHandl.UpdateMenuOfHandl())
	mux.HandleFunc("DELETE /menu/{id}", menuHandl.DeleteMenuOfHandl())

	mux.HandleFunc("POST /orders", invHandl.AddInventoryOfHandl())
	mux.HandleFunc("GET /orders", invHandl.ReadInventoryOfHandl())
	mux.HandleFunc("GET /orders/{id}", invHandl.ReadInventoryOfHandlByID())
	mux.HandleFunc("PUT /orders/{id}", invHandl.UpdateInventoryOfHandl())
	mux.HandleFunc("DELETE /orders/{id}", invHandl.DeleteInventoryOfHandl())

	//mux.HandleFunc("")
}
