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

	invRepo := dal.InventoryFilePath()
	invSvc := service.NewInventoryService(invRepo)
	invHandl := handler.NewInventoryHandler(invSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", invHandl.AddInventoryOfHandl())
	mux.HandleFunc("GET /orders", invHandl.ReadInventoryOfHandl())
	mux.HandleFunc("GET /orders/{id}", invHandl.ReadInventoryOfHandlByID())
	mux.HandleFunc("PUT /orders/{id}", invHandl.UpdateInventoryOfHandl())
	mux.HandleFunc("DELETE /orders/{id}", invHandl.DeleteInventoryOfHandl())

	menuRepo := dal.MenuFilePath()
	menuSvc := service.NewMenuService(menuRepo, invRepo)
	menuHandl := handler.NewMenuHandler(menuSvc)

	mux.HandleFunc("POST /menu", menuHandl.AddMenuOfHandl())
	mux.HandleFunc("GET /menu", menuHandl.ReadMenuOfHandl())
	mux.HandleFunc("GET /menu/{id}", menuHandl.ReadMenuOfHandlByID())
	mux.HandleFunc("PUT /menu/{id}", menuHandl.UpdateMenuOfHandl())
	mux.HandleFunc("DELETE /menu/{id}", menuHandl.DeleteMenuOfHandl())
	//mux.HandleFunc("")
}
